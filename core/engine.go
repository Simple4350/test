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
)

// InitEngine 初始化 BT 引擎
func InitEngine(downloadDir string) {
    mu.Lock()
    defer mu.Unlock()

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

    stats := currentTor.Stats()
    return fmt.Sprintf("%d|%d", stats.DownloadSpeed, stats.BytesReadUseful.Int64())
}

// GetFilePath 获取文件绝对路径
func GetFilePath(fileIndex int) string {
    mu.Lock()
    defer mu.Unlock()

    files := currentTor.Files()
    return filepath.Join(client.Config().DataDir, currentTor.Info().Name, files[fileIndex].Path())
}