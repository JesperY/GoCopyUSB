package backend

import (
	"io"
	"os"
)

// 执行具体的操作行为，该模块仅负责复制操作，同时复制操作过程中的问题应在该模块内部解决

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	// 打开 src 指定的文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	// 延迟关闭
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {

		}
	}(sourceFile)

	// 创建目标文件
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	// 延迟关闭
	defer func(destinationFile *os.File) {
		err := destinationFile.Close()
		if err != nil {

		}
	}(destinationFile)

	// 执行复制
	_, err = io.Copy(destinationFile, sourceFile)
	return err
}
