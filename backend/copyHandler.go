package backend

import (
	"fmt"
	"github.com/JesperY/GoCopyUSB/config"
	"github.com/JesperY/GoCopyUSB/copylogger"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// copyFile 向硬盘中拷贝文件
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

// 目录检查，通过返回 true，继续备份，未通过返回 false，跳过
func checkDir(targetPath string, sourcePath string) bool {
	// 1. 检查源目录是否为白名单目录
	// 2. 检查目标目录是否已存在
	_, err := os.Stat(targetPath)
	// 当目标路径不存在且源路径不在白名单时，通过检查
	// (!os.IsExist(err) || err == nil) 共同判定 targetPath 存在
	// 这是因为 targetPath 存在时可能返回 nil，也可能返回非 ErrNotExist 的 err
	// 但是 os.IsExist(nil) 返回 false 判定为不存在，因此需要额外判断 err == nil
	if !IsWhiteDir(sourcePath) && !(err == nil || os.IsExist(err)) { // 注意短路，先判断 err == nil
		return true
	} else {
		return false
	}
}

func checkFile(targetPath string, sourcePath string) bool {

	// 获取后缀
	suffix := filepath.Ext(targetPath)
	filename := filepath.Base(sourcePath)
	filenameWithoutExt := strings.TrimSuffix(filename, suffix)
	// 当文件名或文件后缀在白名单或源文件未更新，任一不通过检查，返回 false
	if IsWhiteFilename(filenameWithoutExt) || IsWhiteSuffix(suffix) || !isUpdatedFile(targetPath, sourcePath) {
		return false
	}
	return true
}

func doCopy(instance *ole.IDispatch) error {
	// 获取 instance 的 DeviceID 属性并转为 String
	deviceId := oleutil.MustGetProperty(instance, "DeviceID").ToString()
	fmt.Printf("USB Drive inserted: %s\n", deviceId)

	sourcePath := deviceId + `\` // Assume the USB is mounted with a drive letter.
	// 从 json 读取目标路径配置
	targetPath := config.ConfigPtr.TargetDir

	/*
		Copy all files and directories from USB drive to target directory
		filepath.Walk 可以遍历指定目录下的所有文件和目录
		第一个参数为要遍历的目录，第二个参数为回调函数，每遍历到一个文件或目录就调用一次
		此处使用匿名函数，其功能为将遍历到的文件或目录复制到目标路径下
	*/
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			copylogger.SugarLogger.Errorf("Failed walk dir: %v\n", err)
			return err
		}
		// 构建 targetFilePath 作为复制的目标路径
		// TrimPrefix 将 sourcePath 从 path 中去除，只保留 sourcePath 之后的路径字符串，然后拼接到 targetPath
		targetFilePath := filepath.Join(targetPath, strings.TrimPrefix(path, sourcePath))
		// 如果遍历到的是目录，则创建对应目录，保持原权限
		if info.IsDir() {
			// 执行备份前检查
			if !checkDir(targetFilePath, sourcePath) {
				return filepath.SkipDir
			}
			return os.MkdirAll(targetFilePath, info.Mode())
		} else { // 遍历到文件则执行复制操作
			// 备份前检查
			if !checkFile(targetPath, sourcePath) {
				return filepath.SkipDir
			}
			return copyFile(path, targetFilePath)
		}
	})

	if err != nil {
		copylogger.SugarLogger.Errorf("Failed copy file: %v\n", err)
		//log.Println("Error copying files:", err)
	} else {
		copylogger.SugarLogger.Infof("\"All files copied successfully from %s", deviceId)
		//fmt.Println("All files copied successfully from", deviceId)
	}
	// TODO:优化 error 处理
	return err
}

// HandleEvent 插入U盘时的处理逻辑
func HandleEvent(eventSource *ole.IDispatch) {
	// 调用 SWbemEventSource 的 NextEvent 方法，获取下一个事件
	// 返回一个事件的 COM 对象
	eventRaw, err := oleutil.CallMethod(eventSource, "NextEvent", nil)
	if err != nil {
		//fmt.Println("Error getting next event:", err)
		copylogger.SugarLogger.Errorf("Error getting next USB event: %v", err)
		return
	}
	event := eventRaw.ToIDispatch()
	defer event.Release()

	/*
		MustGetProperty 用于获取 COM 对象的属性
		此处尝试从 event 中获取指定的 TargetInstance 属性
		该方法如果获取失败将会引发 panic
		如果不能保证一定可以获取到属性，应该考虑使用 oleutil.GetProperty 并适当处理错误
		对于诸如 __InstanceCreationEvent、__InstanceDeletionEvent 等事件，TargetInstance 属性通常包含了引发事件的实例
		此处，如果创建了新的逻辑磁盘，则会指向 Win32_LogicalDisk 实例
	*/
	targetInstance := oleutil.MustGetProperty(event, "TargetInstance")
	instance := targetInstance.ToIDispatch()
	defer instance.Release()
	if config.ConfigPtr.DelayMinutes > 0 { // 延迟指定的时间后再进行备份
		time.Sleep(time.Duration(config.ConfigPtr.DelayMinutes) * time.Minute)
	}
	err = doCopy(instance)
	if err != nil {
		//fmt.Println("Error copying files:", err)
		copylogger.SugarLogger.Errorf("Error copying files: %v", err)
		return
	}

}
