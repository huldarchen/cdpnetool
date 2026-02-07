import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Info, FolderOpen, Github } from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { useToast } from '@/hooks/use-toast'
import { api } from '@/api'

interface AboutDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function AboutDialog({ open, onOpenChange }: AboutDialogProps) {
  const { t } = useTranslation()
  const { toast } = useToast()

  const [version, setVersion] = useState('')
  const [logDirectory, setLogDirectory] = useState('')
  const [dataDirectory, setDataDirectory] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    if (open) {
      loadInfo()
    }
  }, [open])

  const loadInfo = async () => {
    setIsLoading(true)
    try {
      // 获取版本
      const versionResult = await api.system.getVersion()
      if (versionResult?.success && versionResult.data?.version) {
        setVersion(versionResult.data.version)
      }

      // 获取日志目录
      const logDirResult = await api.system.getLogDirectory()
      if (logDirResult?.success && logDirResult.data?.value) {
        setLogDirectory(logDirResult.data.value)
      }

      // 获取数据目录
      const dataDirResult = await api.system.getDataDirectory()
      if (dataDirResult?.success && dataDirResult.data?.value) {
        setDataDirectory(dataDirResult.data.value)
      }
    } catch (e) {
      console.error('加载关于信息失败:', e)
    } finally {
      setIsLoading(false)
    }
  }

  const handleOpenDirectory = async (path: string) => {
    try {
      await api.system.openDirectory(path)
    } catch (e) {
      toast({
        variant: 'destructive',
        title: t('settings.about.openFolder'),
        description: String(e),
      })
    }
  }

  const handleOpenGitHub = () => {
    window.open('https://github.com/241x/cdpnetool', '_blank')
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-sm">
        <DialogHeader>
          <DialogTitle className="sr-only">{t('settings.about.title')}</DialogTitle>
        </DialogHeader>

        {isLoading ? (
          <div className="flex items-center justify-center py-8">
            <span className="text-muted-foreground">{t('common.loading')}</span>
          </div>
        ) : (
          <div className="space-y-6 py-2">
            {/* 应用信息 */}
            <div className="text-center space-y-3">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary/10">
                <Info className="w-8 h-8 text-primary" />
              </div>
              <div>
                <h2 className="text-2xl font-bold">{t('settings.about.appName')}</h2>
                <p className="text-sm text-muted-foreground mt-1">
                  {t('settings.about.version')} {version || 'v1.1.1'}
                </p>
              </div>
            </div>

            {/* 作者和仓库 */}
            <div className="space-y-2 text-center text-sm">
              <p className="text-muted-foreground">
                {t('settings.about.author')}: <span className="selectable">hemailads@163.com</span>
              </p>
              <Button
                variant="link"
                className="h-auto p-0 text-sm"
                onClick={handleOpenGitHub}
              >
                <Github className="w-3.5 h-3.5 mr-1.5" />
                241x/cdpnetool
              </Button>
            </div>

            {/* 目录操作 */}
            <div className="flex flex-col gap-2 items-center">
              <Button
                variant="outline"
                size="sm"
                className="w-48 justify-center gap-2"
                onClick={() => handleOpenDirectory(dataDirectory)}
                disabled={!dataDirectory}
              >
                <FolderOpen className="w-4 h-4" />
                {t('settings.about.dataDirectory')}
              </Button>
              <Button
                variant="outline"
                size="sm"
                className="w-48 justify-center gap-2"
                onClick={() => handleOpenDirectory(logDirectory)}
                disabled={!logDirectory}
              >
                <FolderOpen className="w-4 h-4" />
                {t('settings.about.logDirectory')}
              </Button>
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}
