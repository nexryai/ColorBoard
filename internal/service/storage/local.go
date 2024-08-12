package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var (
    ErrFileNotFound     = errors.New("file not found")
    ErrPermissionDenied = errors.New("permission denied")
)

type LocalStorageService struct {
    dataDir string
}

func (l *LocalStorageService) CreateFile(reader io.Reader, userId string) (string, error) {
    // IDを生成
    id := fmt.Sprintf("local:%s:%s", userId, uuid.New().String())

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

func (l *LocalStorageService) GetFileUrl(id string, userId string) (string, error) {
    if !strings.HasPrefix(id, fmt.Sprintf("local:%s:", userId)) {
        return "", ErrPermissionDenied 
    }
    
    filePath := filepath.Join(l.dataDir, id)

    // ファイルが存在するか確認
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return "", ErrFileNotFound
    }

    // ファイルパスを返す（ローカルのURL）
    return filePath, nil
}

func (l *LocalStorageService) DeleteFile(id string, userId string) error {
    if !strings.HasPrefix(id, fmt.Sprintf("local:%s:", userId)) {
        return ErrPermissionDenied 
    }
    
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
