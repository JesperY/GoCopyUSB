package backend

import (
	"github.com/JesperY/GoCopyUSB/copylogger"
	"golang.org/x/sys/windows/registry"
	"os"
)

var registryName string = "USBCopier"

func EnableAutoStartUp() {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		// todo 无法打开注册表项，弹窗
		//return fmt.Errorf("error opening registry key: %w", err)
		copylogger.SugarLogger.Errorf("Error opening registry key: %v", err)
	} else {
		// 确保 key 不为 nil 时调用 close()
		defer key.Close()
	}

	executablePath, err := os.Executable()
	if err != nil {
		//return fmt.Errorf("error getting executable path: %w", err)
		// todo 无法设置可执行文件路径，弹窗
		copylogger.SugarLogger.Errorf("Error getting executable path: %v", err)
	}

	err = key.SetStringValue(registryName, executablePath)
	if err != nil {
		//return fmt.Errorf("error setting registry value: %w", err)
		// todo 无法设置注册表项，弹窗
		copylogger.SugarLogger.Errorf("Error setting registry value: %v", err)
	}
}

func DisableAutoStartUp() {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		// todo 无法打开注册表项，弹窗
		//return fmt.Errorf("error opening registry key: %w", err)
		copylogger.SugarLogger.Errorf("Error opening registry key: %v", err)
	} else {
		defer key.Close()
	}

	err = key.DeleteValue(registryName)
	if err != nil {
		//return fmt.Errorf("error deleting registry value: %w", err)
		// todo 无法删除注册表项，弹窗
		copylogger.SugarLogger.Errorf("Error deleting registry value: %v", err)
	}
}

func IsAutoStartUp() bool {

	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.READ)
	if err != nil {
		//fmt.Println("Error opening registry key:", err)
		copylogger.SugarLogger.Errorf("Error opening registry key: %v", err)
		return false
	} else {
		defer key.Close()
	}
	//defer key.Close()

	val, _, err := key.GetStringValue(registryName)
	if err != nil {
		//fmt.Println("USBCopier is not set to auto start.")
		copylogger.SugarLogger.Infof("USBCopier is not set to auto start. %v", err)
		return false
	} else {
		//fmt.Println("USBCopier is set to auto start with value:", val)
		copylogger.SugarLogger.Infof("USBCopier is set to auto start with value: %v", val)
		return true
	}
}
