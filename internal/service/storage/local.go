package storage

import (
    "fmt"
    "github.com/google/uuid"
    "io"
    "os"
    "path/filepath"
)

type LocalStorageService struct {
    dataDir string
}

func (l *LocalStorageService) CreateFile(reader io.Reader) (string, error) {
    // IDを生成
    id := fmt.Sprintf("local:%s", uuid.New().String())

    // 保存するファイルパスを作成
    filePath := filepath.Join(l.dataDir, id)

    // ファイルを作成
    file, err := os.Create(filePath)
    if err != nil {
        return "", fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    // データをファイルに書き込む
    _, err = io.Copy(file, reader)
    if err != nil {
        return "", fmt.Errorf("failed to write to file: %w", err)
    }

    return id, nil
}

func (l *LocalStorageService) GetFileUrl(id string) (string, error) {
    filePath := filepath.Join(l.dataDir, id)

    // ファイルが存在するか確認
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return "", fmt.Errorf("file not found")
    }

    // ファイルパスを返す（ローカルのURL）
    return filePath, nil
}

func (l *LocalStorageService) DeleteFile(id string) error {
    filePath := filepath.Join(l.dataDir, id)

    // ファイルを削除
    if err := os.Remove(filePath); err != nil {
        return fmt.Errorf("failed to delete file: %w", err)
    }

    return nil
}

func NewLocalStorageService() (*LocalStorageService, error) {
    dataDir := os.Getenv("DATA_DIR")
    if dataDir == "" {
        return nil, fmt.Errorf("DATA_DIR environment variable is not set")
    }

    return &LocalStorageService{dataDir: dataDir}, nil
}
