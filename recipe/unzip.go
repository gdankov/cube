package main

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Unzip extracts a ZIP file located at src into a specified directory. It uses package 'archive/zip' internally.
// Giving an empty string as the target dir will unzip the contents in the current directory.
func Unzip(src, targetDir string) error {
	if targetDir == "" {
		return errors.New("target directory cannot be empty")
	}

	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		destPath := filepath.Join(targetDir, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, file.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := extractFile(file, destPath); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(src *zip.File, destPath string) error {
	reader, err := src.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, reader); err != nil {
		return err
	}

	err = os.Chmod(destPath, src.Mode())
	return err
}
