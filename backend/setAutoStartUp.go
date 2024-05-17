package backend

import (
	"github.com/JesperY/GoCopyUSB/config"
	"github.com/JesperY/GoCopyUSB/logger"
	"golang.org/x/sys/windows/registry"
	"os"
)

var registryName string = "USBCopier"

func EnableAutoStartUp() {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		// todo 无法打开注册表项，弹窗
		//return fmt.Errorf("error opening registry key: %w", err)
		logger.SugarLogger.Errorf("Error opening registry key: %v", err)
	} else {
		// 确保 key 不为 nil 时调用 close()
		defer key.Close()
	}

	executablePath, err := os.Executable()
	if err != nil {
		//return fmt.Errorf("error getting executable path: %w", err)
		// todo 无法设置可执行文件路径，弹窗
		logger.SugarLogger.Errorf("Error getting executable path: %v", err)
	}

	err = key.SetStringValue(registryName, executablePath)
	if err != nil {
		//return fmt.Errorf("error setting registry value: %w", err)
		// todo 无法设置注册表项，弹窗
		logger.SugarLogger.Errorf("Error setting registry value: %v", err)
	}
}

func DisableAutoStartUp() {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		// todo 无法打开注册表项，弹窗
		//return fmt.Errorf("error opening registry key: %w", err)
		logger.SugarLogger.Errorf("Error opening registry key: %v", err)
	} else {
		defer key.Close()
	}

	err = key.DeleteValue(registryName)
	if err != nil {
		//return fmt.Errorf("error deleting registry value: %w", err)
		// todo 无法删除注册表项，弹窗
		logger.SugarLogger.Errorf("Error deleting registry value: %v", err)
	}
}

func IsAutoStartUp() bool {

	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.READ)
	if err != nil {
		//fmt.Println("Error opening registry key:", err)
		logger.SugarLogger.Errorf("Error opening registry key: %v", err)
		return false
	} else {
		defer key.Close()
	}
	//defer key.Close()

	val, _, err := key.GetStringValue(registryName)
	if err != nil {
		//fmt.Println("USBCopier is not set to auto start.")
		logger.SugarLogger.Infof("USBCopier is not set to auto start. %v", err)
		return false
	} else {
		//fmt.Println("USBCopier is set to auto start with value:", val)
		logger.SugarLogger.Infof("USBCopier is set to auto start with value: %v", val)
		return true
	}
}
func InitCheck() {
	// todo 目标目录检查
	targetPath := config.ConfigPtr.TargetDir
	fileInfo, err := os.Stat(targetPath)
	// 目标文件夹不存在，创建
	if err != nil {
		err := os.Mkdir(targetPath, 0777)
		if err != nil {
			// todo 没有创建权限，无法继续，弹窗
			logger.SugarLogger.Errorf("Permission deny, can not create target dir, %v", err)
			return
		}
	}
	if fileInfo.Mode().Perm()&os.ModePerm == 0 {
		// todo 没有修改权限，无法继续，弹窗
		logger.SugarLogger.Errorf("Permission deny, can not modify target dir")
	}

	// 开机自启动项检查
	if config.ConfigPtr.AutoStartUp {
		// 配置开机自启动，检查注册表项
		// 如果未配置自启动则配置
		if !IsAutoStartUp() {
			// 已配置自启动
			//backend.DisableAutoStartUp()
			EnableAutoStartUp()
		}

	} else {
		// 配置开机不启动，检查注册表项
		// 如果已配置则删除
		if IsAutoStartUp() {
			DisableAutoStartUp()
		}

	}

}
