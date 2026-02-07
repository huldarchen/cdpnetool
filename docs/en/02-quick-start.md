# Quick Start

This guide will help you complete the installation, configuration, and creation of your first network interception rule in cdpnetool within 5 minutes.

---

## System Requirements

- **Operating System**: Windows 10/11 (official installer provided), macOS, Linux (theoretically supported, feedback welcome)
- **Browser**: Chrome/Chromium/Edge browser (optional, tool can launch built-in browser)
- **Other**: No additional dependencies required

---

## Installation and Running

### Method 1: Using Precompiled Version (Recommended)

#### Windows Users

1. Visit [GitHub Releases](https://github.com/241x/cdpnetool/releases) page
2. Download the latest Windows installer
3. Run the `.exe` installer and follow the prompts
4. Or download the portable `.zip`, extract and run directly

#### macOS and Linux Users

> ⚠️ **Note**: Currently, official builds are only provided for Windows. macOS and Linux platforms are theoretically supported but not fully tested. If you encounter issues on these platforms, please provide feedback through [Issues](https://github.com/241x/cdpnetool/issues).

Since precompiled versions are not provided, please refer to the "Building from Source" section below.

---

### Method 2: Building from Source

If you want to run from source or participate in development, you need the following environment:

#### Prerequisites

- Go 1.22+
- Node.js 22+
- Wails CLI v2

#### Build Steps

**1. Clone the project**

```bash
git clone https://github.com/241x/cdpnetool.git
cd cdpnetool
```

**2. Install Wails CLI** (if not already installed)

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

**3. Run in development mode**

```bash
wails dev
```

**4. Build production version**

```bash
wails build
```

After building, the executable file will be located in the `build/bin/` directory.

> 💡 **Tip**: For detailed development configuration and build instructions, please refer to the Development Guide

---

## 5-Minute Tutorial

### Step 1: Launch Browser

1. Open the cdpnetool application
2. Click the "Launch Browser" button on the left side of the toolbar
3. The tool will automatically launch an independent Chrome instance
4. The DevTools URL will be automatically filled in the input box (usually `http://localhost:9222`)

> 💡 **Tip**: If you want to use an already running browser, you can manually enter its DevTools URL

---

### Step 2: Connect to Browser

1. Confirm the DevTools URL is correctly filled in
2. Click the "Connect" button
3. After successful connection, the status indicator turns green and shows "Connected"

> ⚠️ **Common Issues**: If connection fails, please check:
> - Whether the browser has started normally
> - Whether the DevTools URL format is correct (`http://localhost:port`)
> - Whether the port is blocked by the firewall

---

### Step 3: Attach Target Page

1. Switch to the "Targets" tab
2. Click the refresh button in the upper right corner to get the page list
3. Find the page you want to intercept in the target list
4. Click the "Attach" button on the right side of the page
5. The button changes to "Attached" status

> 💡 **Note**: cdpnetool uses single-target mode, only one page can be attached at a time. Switching targets will automatically detach the old target.

---

### Step 4: Create Your First Rule Configuration

1. Switch to the "Rules" tab
2. If there is no configuration yet, click the "+" button on the left to create a new configuration
3. The configuration will automatically generate a default name, which can be double-clicked to rename to "Test Configuration"

#### Add Rule

Click the "Add Rule" button on the right to start configuring the rule:

**Basic Information:**
- Modify rule name: `Mock User Info API`
- Keep rule enabled (default)
- Stage selection: `response` (response stage)

**Configure Match Conditions:**
1. Click "Add Condition" in the "AllOf" area
2. Condition type selection: `urlContains` (URL contains)
3. Enter match value: `/api/user/info`

**Configure Actions:**
1. Click "Add Action"
2. Action type selection: `setBody` (set response body)
3. Enter JSON data:

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "Test User",
    "email": "test@example.com"
  }
}
```

4. Press `Ctrl+S` (Windows) or `Cmd+S` (macOS) to save the configuration

---

### Step 5: Enable Interception

1. In the left configuration list, find the "Test Configuration" you just created
2. Click the switch button on the left side of the configuration
3. The tool will prompt "Configuration 'Test Configuration' enabled"
4. The configuration item will display a green "· Running" indicator

> ⚠️ **Note**: You must complete both "Connect Browser" and "Attach Target" before enabling interception

---

### Step 6: Test Effect

1. Visit a page containing the `/api/user/info` API in the browser
2. Switch to the "Events" tab
3. You can see the intercepted request in the "Matched Requests" list
4. Click the event card to view details and confirm the response body has been replaced with our configured JSON data

> 🎉 **Congratulations**! You have successfully completed your first network interception rule

---

## Interface Features

### Top Toolbar

- **Launch Browser**: One-click launch of independent Chrome instance
- **DevTools URL Input Box**: Manually enter or auto-fill browser debugging address
- **Connect/Disconnect Button**: Manage connection status with browser
- **Status Indicator**: Display connection status and number of attached targets
- **Language Switch**: Chinese/English interface switching
- **Settings Button**: Open system settings dialog

### Targets Panel (Target Management)

- View all page tabs in current browser
- Display page title and URL
- Attach/detach target pages

### Rules Panel (Rule Configuration)

- **Left Side**: Configuration list, supports create, switch, rename, delete
- **Right Side**: Rule editor
  - Visual editing mode (default)
  - JSON source code editing mode (advanced)
- **Configuration Management**: Import/export, save, enable/disable

### Events Panel (Event Monitoring)

- **Matched Requests**: Display list of requests that hit rules
- **Unmatched Requests**: Display list of requests that did not hit rules
- **Event Details**: DevTools-style request response viewing panel
  - Request details (URL, method, headers, body)
  - Response details (status code, headers, body)
  - Matched rule information

### Network Panel (Request Preview)

- **Full Traffic View**: Display all captured network requests (regardless of rule matching)
- **Enable/Disable Capture**: Control whether to listen to all network requests
- **Details Panel**: DevTools-style request response viewing
- **Clear History**: Clear captured request list

> 💡 **Tip**: Network panel requires attaching a target page before enabling capture

---

## System Settings Usage

### Open Settings Dialog

Click the gear icon button on the right side of the top toolbar to open the system settings dialog.

### General Settings

**Language Settings**
- Support: Simplified Chinese, English
- Takes effect immediately after switching, interface updates automatically
- Can also click the language button on the top toolbar for quick switching

**Theme Settings**
- Support: Light, Dark, Follow System
- Takes effect immediately after switching
- Follow System: Automatically adapts to the operating system's dark/light mode

### Browser Settings

**Extra Launch Arguments**
- Custom browser launch parameters, one parameter per line
- Common parameter examples:
  - `--disable-extensions`: Disable extensions
  - `--disable-gpu`: Disable GPU acceleration
  - `--incognito`: Incognito mode
  - `--window-size=1920,1080`: Set window size

**Browser Path**
- Can specify the full path to the browser executable file
- Leave blank to automatically detect system Chrome/Edge
- Click "Select File" button to browse and select

### About Page

Click the info icon in the bottom status bar to view:
- Application version number
- GitHub repository link
- Quick access to data directory
- Quick access to log directory

---

## Common Operation Tips

### Keyboard Shortcuts

- `Ctrl/Cmd + S`: Save current configuration

### Configuration Management Tips

- Create multiple configurations for quick switching between different scenarios
- Configuration files are in standard JSON format, can be exported and shared with team members
- Configuration format and compatibility are automatically validated during import

### Interception Notes

- Modified rules need to be re-saved to take effect
- If configuration is already enabled, saving will automatically reload rules
- Switching attached targets will automatically stop interception, need to re-enable
- Disconnecting will stop interception and clear event list

### Common Connection Issues

- **Cannot Connect**: Check DevTools URL format (`http://localhost:port`)
- **Connection Dropped**: Browser closed or network exception will cause connection interruption
- **Port Occupied**: If port 9222 is occupied, the tool will automatically select another port

---

## Next Steps

🎉 Congratulations! You have mastered the basic usage of cdpnetool. Next you can:

- 📚 View [Rule Configuration Reference](./03-rule-reference.md) to learn all matching conditions and action types
- ❓ Have questions? Check [FAQ](./04-faq.md) for solutions
