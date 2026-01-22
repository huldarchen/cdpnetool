package browser

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// Options 浏览器启动选项
type Options struct {
	ExecPath            string   // 浏览器可执行文件路径
	UserDataDir         string   // 用户数据目录
	RemoteDebuggingPort int      // CDP端口，0表示自动选择
	Headless            bool     // 是否以无头模式启动
	Args                []string // 额外启动参数
	Env                 []string // 额外环境变量
}

// Browser 已启动的浏览器进程句柄
type Browser struct {
	cmd         *exec.Cmd
	DevToolsURL string
	port        int
}

// Start 启动浏览器并等待CDP服务就绪
func Start(ctx context.Context, opts Options) (*Browser, error) {
	exe := opts.ExecPath
	if exe == "" {
		exe = defaultChromePath()
	}
	if exe == "" {
		return nil, errors.New("chrome executable not found")
	}

	port := opts.RemoteDebuggingPort
	if port == 0 {
		port = 9222
	}

	finalPort, err := pickPort(port)
	if err != nil {
		return nil, fmt.Errorf("failed to pick port: %w", err)
	}

	port = finalPort
	args := buildLaunchArgs(port, opts)
	cmd := exec.CommandContext(ctx, exe, args...)
	if len(opts.Env) > 0 {
		cmd.Env = append(os.Environ(), opts.Env...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start browser: %w", err)
	}

	b := &Browser{cmd: cmd, DevToolsURL: fmt.Sprintf("http://127.0.0.1:%d", port), port: port}
	waitCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := waitDevToolsReady(waitCtx, b.DevToolsURL); err != nil {
		_ = b.Stop(2 * time.Second)
		return nil, fmt.Errorf("devtools not ready: %w", err)
	}

	return b, nil
}

// Stop 关闭浏览器进程
func (b *Browser) Stop(timeout time.Duration) error {
	if b == nil || b.cmd == nil || b.cmd.Process == nil {
		return nil
	}
	done := make(chan error, 1)
	go func() { done <- b.cmd.Wait() }()
	// Windows上直接Kill以避免悬挂
	_ = b.cmd.Process.Kill()
	select {
	case <-time.After(timeout):
		return errors.New("browser stop timeout")
	case err := <-done:
		return err
	}
}

// defaultChromePath 返回常见的 Chrome 可执行路径（跨平台）
func defaultChromePath() string {
	candidates := getChromePaths()
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	for _, name := range []string{"chrome", "google-chrome", "chromium", "chromium-browser"} {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}

	return ""
}

// getChromePaths 根据操作系统返回可能的 Chrome 路径
func getChromePaths() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			filepath.Join(os.Getenv("ProgramFiles"), "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(os.Getenv("ProgramFiles(x86)"), "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "Application", "chrome.exe"),
		}
	case "darwin":
		return []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			filepath.Join(os.Getenv("HOME"), "Applications", "Google Chrome.app", "Contents", "MacOS", "Google Chrome"),
		}
	case "linux":
		return []string{
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/snap/bin/chromium",
		}
	default:
		return nil
	}
}

// pickPort 尝试使用指定端口，如果被占用则选择随机空闲端口
func pickPort(preferred int) (int, error) {
	// 先尝试首选端口
	if preferred > 0 {
		l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", preferred))
		if err == nil {
			_ = l.Close()
			return preferred, nil
		}
	}

	// 首选端口不可用，选择随机空闲端口
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("failed to find free port: %w", err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

// buildLaunchArgs 构建浏览器启动参数
func buildLaunchArgs(port int, opts Options) []string {
	args := []string{
		fmt.Sprintf("--remote-debugging-port=%d", port),
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-background-networking",
		"--disable-background-timer-throttling",
		"--disable-backgrounding-occluded-windows",
		"--disable-breakpad",
		"--disable-client-side-phishing-detection",
		"--disable-default-apps",
		"--disable-extensions",
		"--disable-hang-monitor",
		"--disable-prompt-on-repost",
		"--disable-renderer-backgrounding",
		"--disable-sync",
		"--disable-translate",
		"--metrics-recording-only",
		"--safebrowsing-disable-auto-update",
	}

	// Linux 环境下添加额外参数
	if runtime.GOOS == "linux" {
		args = append(args, "--disable-dev-shm-usage")
	}

	// 用户数据目录
	if opts.UserDataDir != "" {
		_ = os.MkdirAll(opts.UserDataDir, 0o755)
		args = append(args, fmt.Sprintf("--user-data-dir=%s", opts.UserDataDir))
	} else {
		// 使用带时间戳的临时目录，避免冲突
		dir := filepath.Join(os.TempDir(), fmt.Sprintf("cdpnetool-chrome-%d", time.Now().Unix()))
		_ = os.MkdirAll(dir, 0o755)
		args = append(args, fmt.Sprintf("--user-data-dir=%s", dir))
	}

	// 无头模式
	if opts.Headless {
		args = append(args, "--headless=new", "--disable-gpu")
	}

	// 额外参数
	if len(opts.Args) > 0 {
		args = append(args, opts.Args...)
	}

	return args
}

// waitDevToolsReady 轮询 DevTools 服务是否就绪
func waitDevToolsReady(ctx context.Context, base string) error {
	url := fmt.Sprintf("%s/json/version", base)
	cli := &http.Client{Timeout: 500 * time.Millisecond}
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("devtools not ready after timeout: %w", ctx.Err())
		case <-ticker.C:
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				continue
			}
			resp, err := cli.Do(req)
			if err == nil && resp.StatusCode == 200 {
				resp.Body.Close()
				return nil
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
}
