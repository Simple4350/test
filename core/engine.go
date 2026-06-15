package core

import (
    "fmt"
    "path/filepath"
    "sync"

    "github.com/anacrolix/torrent"
    "github.com/anacrolix/torrent/storage"
)

var (
    client     *torrent.Client
    currentTor *torrent.Torrent
    mu         sync.Mutex
    saveDir    string // 新增：保存下载目录，因为新版库隐藏了 client.Config
)

// InitEngine 初始化 BT 引擎
func InitEngine(downloadDir string) {
    mu.Lock()
    defer mu.Unlock()

    saveDir = downloadDir // 保存目录
    cfg := torrent.NewDefaultClientConfig()
    cfg.DefaultStorage = storage.NewFile(downloadDir)
    cfg.Seed = false

    var err error
    client, err = torrent.NewClient(cfg)
    if err != nil {
        fmt.Println("引擎初始化失败:", err)
        return
    }
}

// AddMagnet 添加磁力链并等待获取元数据
func AddMagnet(magnet string) string {
    mu.Lock()
    defer mu.Unlock()

    t, err := client.AddMagnet(magnet)
    if err != nil {
        return "失败"
    }
    currentTor = t
    <-t.GotInfo() // 阻塞等待下载到种子信息
    return "ok"
}

// GetFileList 获取文件列表
func GetFileList() string {
    mu.Lock()
    defer mu.Unlock()

    var result string
    for i, f := range currentTor.Files() {
        result += fmt.Sprintf("%d|%s|%d\n", i, f.Path(), f.Length())
    }
    return result
}

// StartSequentialDownload 顺序下载指定文件（边下边播核心）
func StartSequentialDownload(fileIndex int) string {
    mu.Lock()
    defer mu.Unlock()

    files := currentTor.Files()
    if fileIndex < 0 || fileIndex >= len(files) {
        return "越界"
    }
    files[fileIndex].SetPriority(torrent.PiecePriorityNow)
    currentTor.DownloadAll()
    return "ok"
}

// GetDownloadStats 获取下载速度
func GetDownloadStats() string {
    mu.Lock()
    defer mu.Unlock()

    if currentTor == nil || currentTor.Info() == nil {
        return "0|0"
    }
    
    stats := currentTor.Stats()
    // 新版库去掉了 DownloadSpeed，我们用读取的总字节数代替
    return fmt.Sprintf("0|%d", stats.BytesRead.Int64())
}

// GetFilePath 获取文件绝对路径
func GetFilePath(fileIndex int) string {
    mu.Lock()
    defer mu.Unlock()

    files := currentTor.Files()
    // 使用我们保存的 saveDir 替代 client.Config()
    return filepath.Join(saveDir, currentTor.Info().Name, files[fileIndex].Path())
}