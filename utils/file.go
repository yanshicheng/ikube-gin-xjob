package utils

import (
	"os"
	"path/filepath"
)

// 1. 判断文件是否存在
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// 2. 判断文件夹是否存在
func FolderExists(foldername string) bool {
	info, err := os.Stat(foldername)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// 3. 创建文件
func CreateFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

// 4. 创建文件夹
func CreateFolder(foldername string) error {
	err := os.Mkdir(foldername, 0755)
	if err != nil {
		return err
	}
	return nil
}

// 5. 递归创建文件夹
func CreateFolderRecursive(foldername string) error {
	err := os.MkdirAll(foldername, 0755)
	if err != nil {
		return err
	}
	return nil
}

// 6. 递归创建文件夹和文件
func CreateFileInFolderRecursive(filepathString string) error {
	folder := filepath.Dir(filepathString)
	err := os.MkdirAll(folder, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filepathString)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}
