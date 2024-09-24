package util

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func Unzip(zipFile string, destDir string) {
	// Open the zip file
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		fmt.Printf("\033[31m無法開啟壓縮檔: %v\033[0m\n", err)
		os.Exit(1)
	}
	defer r.Close()

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		fmt.Printf("\033[31m無法創建資料夾: %v\033[0m\n", err)
		os.Exit(1)
	}

	// Extract the files
	for _, f := range r.File {
		// Create the destination file path
		destPath := filepath.Join(destDir, f.Name)

		// Check if the destination file is a directory
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, f.FileInfo().Mode()); err != nil {
				fmt.Printf("\033[31m無法創建資料夾: %v\033[0m\n", err)
				os.Exit(1)
			}
			continue
		}

		// Create the destination file
		outFile, err := os.Create(destPath)
		if err != nil {
			fmt.Printf("\033[31m無法創建檔案: %v\033[0m\n", err)
			os.Exit(1)
		}
		defer outFile.Close()

		// Extract the file content
		rc, err := f.Open()
		if err != nil {
			fmt.Printf("\033[31m無法開啟檔案: %v\033[0m\n", err)
			os.Exit(1)
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			fmt.Printf("\033[31m無法複製檔案: %v\033[0m\n", err)
			os.Exit(1)
		}
	}
}

func IsFolderEmpty(folderPath string) bool {
	entries, err := ioutil.ReadDir(folderPath)
	if err != nil {
		fmt.Printf("\033[31m無法讀取資料夾: %v\033[0m\n", err)
		return false
	}

	return len(entries) == 0
}

func RemoveTopLevelFolderIfExists(parentDir string) {
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		fmt.Println("\033[31m資料夾不存在:", err, "\033[0m")
		os.Exit(1)
	}

	// Get the list of files and directories in the directory
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		fmt.Println("\033[31m無法讀取資料夾:", err, "\033[0m")
		os.Exit(1)
	}

	if len(entries) != 1 {
		return
	}

	childInfo, err := entries[0].Info()
	if err != nil {
		fmt.Println("\033[31m無法讀取子項目:", err, "\033[0m")
		os.Exit(1)
	}

	if !childInfo.IsDir() {
		return
	}

	childPath := path.Join(parentDir, entries[0].Name())

	childEntries, err := os.ReadDir(childPath)
	if err != nil {
		fmt.Println("\033[31m無法讀取資料夾:", err, "\033[0m")
		os.Exit(1)
	}
	// Move each entry to the parent directory
	for _, entry := range childEntries {
		srcPath := filepath.Join(childPath, entry.Name())
		dstPath := filepath.Join(parentDir, entry.Name())

		// Check if the destination file or directory already exists
		if _, err := os.Stat(dstPath); !os.IsNotExist(err) {
			fmt.Printf("\033[31m檔案: %s 以存在于目標資料夾\033[0m", dstPath)
			os.Exit(1)
		}

		// Move the entry
		err = os.Rename(srcPath, dstPath)
		if err != nil {
			fmt.Printf("\033[31m無法將 %s 移動至 %s: %v\033[0m", srcPath, dstPath, err)
			os.Exit(1)
		}
	}

	// Remove the original directory
	err = os.Remove(childPath)
	if err != nil {
		fmt.Printf("\033[31m無法刪除 %s: %v\033[0m", childPath, err)
		os.Exit(1)
	}
}

func GiveExecutePermission(filePath string) {
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("\033[31m檔案不存在:", filePath, "\033[0m")
		os.Exit(1)
	}

	// Get the current file mode
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("\033[31m無法讀取檔案資訊:", err, "\033[0m")
		os.Exit(1)
	}

	// Add execute permission for the owner, group, and others
	newMode := fileInfo.Mode() | 0111
	err = os.Chmod(filePath, newMode)
	if err != nil {
		fmt.Println("\033[31m無法設定權限:", err, "\033[0m")
		return
	}

	// The file now has execute permission
	println("\033[32m成功設定", filePath, "的權限", "\033[0m")
}
