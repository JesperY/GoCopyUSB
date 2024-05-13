package backend

import (
	"github.com/JesperY/GoCopyUSB/config"
	"os"
)

//监听到 USB 事件后，进行复制操作前
//执行一系列检查操作，例如白名单、按需复制等

// IsWhiteDir 检测是否为白名单目录
func IsWhiteDir(sourcePath string) bool {
	dirList := config.ConfigPtr.WhiteListDir
	//var configPtr = &config.Config{}
	//configPtr.ReadConfigWithPath("../config/config.json")
	//dirList := configPtr.WhiteListDir
	for _, dir := range dirList {
		//fmt.Println(dir, sourcePath)
		if sourcePath == dir {
			return true
		}
	}
	return false
}

// IsWhiteFilename 检测是否为白名单文件名（忽略后缀）
func IsWhiteFilename(filenameWithoutSuffix string) bool {
	filenameList := config.ConfigPtr.WhiteListFilename
	//var configPtr = &config.Config{}
	//configPtr.ReadConfigWithPath("../config/config.json")
	//filenameList := configPtr.WhiteListFilename
	for _, whiteFilename := range filenameList {
		if filenameWithoutSuffix == whiteFilename {
			return true
		}
	}
	return false
}

// IsWhiteSuffix 检测是否为白名单后缀
func IsWhiteSuffix(suffix string) bool {
	suffixList := config.ConfigPtr.WhiteListSuffix
	//var configPtr = &config.Config{}
	//configPtr.ReadConfigWithPath("../config/config.json")
	//suffixList := configPtr.WhiteListSuffix
	for _, whiteSuffix := range suffixList {
		if suffix == whiteSuffix {
			return true
		}
	}
	return false
}

// IsExisted 已存在目录不再重复创建
//func IsExisted(targetPath string) bool {
//	//_, err := os.Stat(targetPath)
//	//if err == nil {
//	//	return true
//	//} else {
//	//	return os.IsExist(err)
//	//}
//	// 逆否命题？
//	// 实际上当 targetPath 存在时，err = nil，而 os.IsExist(nil) 返回 false
//	// 为了避免如上的额外 if 判断，采用 os.IsNotExist()
//	_, err := os.Stat(targetPath)
//	return !os.IsNotExist(err)
//}

func isUpdatedFile(targetPath string, sourcePath string) bool {
	// 目标文件不存在则返回 true，执行备份
	targetFileInfo, err := os.Stat(targetPath)
	sourceFileInfo, _ := os.Stat(sourcePath)
	if !(err != nil || !os.IsExist(err)) {
		return true
	}
	// 目标文件存在
	// 对于目标文件与源文件修改时间，目标文件时间早于源文件时执行备份
	if targetFileInfo.ModTime().Before(sourceFileInfo.ModTime()) {
		return true
	}
	return false
}
