package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func main() {
	err := ole.CoInitialize(0)
	if err != nil {
		log.Fatal(err)
	}
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		log.Fatal(err)
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatal(err)
	}
	defer wmi.Release()

	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		log.Fatal(err)
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	queryString := "SELECT * FROM __InstanceCreationEvent WITHIN 2 WHERE TargetInstance ISA 'Win32_LogicalDisk' AND TargetInstance.DriveType = 2"
	resultRaw, err := oleutil.CallMethod(service, "ExecNotificationQuery", queryString)
	if err != nil {
		log.Fatal(err)
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	fmt.Println("Listening for USB drive insertion events...")
	for {
		handleEvent(result)
	}
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {

		}
	}(sourceFile)

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destinationFile *os.File) {
		err := destinationFile.Close()
		if err != nil {

		}
	}(destinationFile)

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

// handleEvent 插入U盘时的处理逻辑
func handleEvent(result *ole.IDispatch) {
	eventRaw, err := oleutil.CallMethod(result, "NextEvent", nil)
	if err != nil {
		log.Fatal(err)
	}
	event := eventRaw.ToIDispatch()
	defer event.Release()

	targetInstance := oleutil.MustGetProperty(event, "TargetInstance")
	instance := targetInstance.ToIDispatch()
	defer instance.Release()

	deviceId := oleutil.MustGetProperty(instance, "DeviceID").ToString()
	fmt.Printf("USB Drive inserted: %s\n", deviceId)

	sourcePath := deviceId + `\` // Assume the USB is mounted with a drive letter.
	targetPath := `D:\TargetDirectory\`

	// Copy all files and directories from USB drive to target directory
	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		targetFilePath := filepath.Join(targetPath, strings.TrimPrefix(path, sourcePath))
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
}
