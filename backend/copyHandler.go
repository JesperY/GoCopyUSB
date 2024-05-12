package backend

import (
	"fmt"
	"github.com/JesperY/GoCopyUSB/config"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func getInstance(EventSource *ole.IDispatch) *ole.IDispatch {
	// 调用 SWbemEventSource 的 NextEvent 方法，获取下一个事件
	// 返回一个事件的 COM 对象
	eventRaw, err := oleutil.CallMethod(EventSource, "NextEvent", nil)
	if err != nil {
		log.Fatal(err)
	}
	event := eventRaw.ToIDispatch()
	//defer event.Release()

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
	//defer instance.Release()
	return instance
}

func doCheck(targetPath string, info os.FileInfo) error {
	skip := false
	// TODO:检查源路径是否在白名单中

	// TODO:检查目标路径和源路径的修改日期，避免重复备份

	// 不满足备份要求时跳过
	if skip {
		return filepath.SkipDir
	} else {
		return nil
	}
}

func doCopy(instance *ole.IDispatch) error {
	// 获取 instance 的 DeviceID 属性并转为 String
	deviceId := oleutil.MustGetProperty(instance, "DeviceID").ToString()
	fmt.Printf("USB Drive inserted: %s\n", deviceId)

	sourcePath := deviceId + `\` // Assume the USB is mounted with a drive letter.
	// 从 json 读取目标路径配置
	targetPath := config.ConfigPtr.TargetDir

	var data, _ = os.ReadFile("config/config.json")
	os.Stdout.Write(data)

	/*
		Copy all files and directories from USB drive to target directory
		filepath.Walk 可以遍历指定目录下的所有文件和目录
		第一个参数为要遍历的目录，第二个参数为回调函数，每遍历到一个文件或目录就调用一次
		此处使用匿名函数，其功能为将遍历到的文件或目录复制到目标路径下
	*/
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 构建 targetFilePath 作为复制的目标路径
		// TrimPrefix 将 sourcePath 从 path 中去除，只保留 sourcePath 之后的路径字符串，然后拼接到 targetPath
		targetFilePath := filepath.Join(targetPath, strings.TrimPrefix(path, sourcePath))
		// todo 执行备份前检查
		err = doCheck(targetFilePath, info)
		if err != nil {
			return err
		}
		// 如果遍历到的是目录，则创建对应目录，保持原权限，否则执行赋值操作
		if info.IsDir() {
			return os.MkdirAll(targetFilePath, info.Mode())
		} else {
			return copyFile(path, targetFilePath)
		}
	})

	if err != nil {
		log.Println("Error copying files:", err)
	} else {
		fmt.Println("All files copied successfully from", deviceId)
	}
	// TODO:优化 error 处理
	return err
}

// HandleEvent 插入U盘时的处理逻辑
func HandleEvent(result *ole.IDispatch) {

	instance := getInstance(result)
	//defer instance.Release()

	err := doCopy(instance)
	if err != nil {
		log.Fatal(err)
	}

}
