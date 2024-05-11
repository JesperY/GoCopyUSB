package backend

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 监听 USB 事件
// 按需对监听到的进行进行处理

func Listener() {
	// 初始化当前线程，设置线程的 COM 环境，这是使用 COM 对象所必须的
	// 每个线程只需要调用一次，在结束使用 COM 之前应该调用 CoUninitialize()
	// 重复调用该方法会返回 S_FALSE 错误，表明该线程已经初始化
	// 如果环境不支持 COM，也会返回错误
	// 参数 p uintptr 用于指定线程的并发模式，但是该方法总是设置为 0，因为其功能在该方法中没有实际功能
	// CoInitializeEx() 方法可以可以使用该参数
	// 该方法被调用后会初始化内部结构并分配内存，以便创建和管理 COM 对象
	// 初始化线程为 STA（单线程单元），一个线程负责所有 COM 对象
	// 注册消息循环，处理和转发消息和命令
	err := ole.CoInitialize(0)
	if err != nil {
		log.Fatal(err)
	}
	// 清理线程环境，释放资源，取消 COM 的初始化
	defer ole.CoUninitialize()

	// 创建 COM 对象，该方法简化了创建过程
	// programID 指定创建类型
	// WbemScripting.SWbemLocator 是 Windows Management Instrumentation （WMI） 的一部分，用于访问各种系统管理信息
	// 返回一个 *ole.IUnknown 和一个 err
	// *ole.IUnknown 为接口类型，用于接受创建的 COM 对象，IUnknown 是所有 COM 对象的基本接口
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		log.Fatal(err)
	}
	// 释放 COM 对象
	defer unknown.Release()

	// 该方法从一个 COM 对象中请求访问其他接口
	// unknown 对象为 COM 最基础的 ole.IUnknown 接口类型，因此需要使用该方法尝试获取其他接口类型
	// 该方法接受一个 *GUID 类型参数，是一个全局唯一标识符，COM 中每一个接口都有一个唯一标识符
	// 此处请求 ole.IID_IDispatch，表示 IDispatch 接口
	// 如果请求成功，则返回对应的接口类型，用于方法调用
	// 失败信息包含在 err 中
	// IDispatch 接口用于实现自动化和晚绑定，允许在运行时发现对象的属性和方法，是 IUnknown 接口的扩展
	// IDispatch 接口有如下四个方法：
	// 1.GetTypeInfoCount()：获取对象提供的 ITypeInfo 接口数目，ITypeInfo 提供了有关对象的方法和属性的元数据
	// 2.GetTypeInfo()：访问对象类型信息，例如属性、函数、接口等
	// 3.GetIDsOfNames()：将一组字符串名称转为对应的内部数值标识符（DISPID），用于后续调用 Invoke() 时标记哪个成员被调用
	// 4.Invoke()：IDispatch() 接口的核心，允许动态的调用方法、访问属性或触发对象事件
	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatal(err)
	}
	// 释放对象
	defer wmi.Release()

	// CallMethod 用于调用 COM 对象的方法，返回类型为 *ole.VARIANT，是一个通用的数据类型，能够包含不同类型的数据，包括整数、浮点数、字符串、布尔值，以及 COM 对象等
	// 该方法需要传入 *ole.IDispatch 类型从而通过该接口调用字符串中指定的方法
	// 此处调用 wmi 的 ConnectServer 方法，该方法属于 SWbemLocator 对象，用于连接到 WMI 服务，从而允许访问和操作系统信息
	// 通过调用 ConnectServer 方法，CallMethod 会返回一个 *SWbemServices 指针，从而进一步操作 WMI 服务
	// ConnectServer 方法也会接受参数，从而指定链接的详细信息，此处使用默认配置
	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		log.Fatal(err)
	}
	// 将 *ole.VARIANT 转为 *IDispatch，从而可以继续使用 CallMethod 调用相关方法
	service := serviceRaw.ToIDispatch()
	// 释放对象
	defer service.Release()

	// WMI Query Language(WQL)，访问和查询 Windows 管理信息
	// __InstanceCreationEvent 是一个 WMI 系统类，用于监控对象创建事件，在系统创建新对象时动态生成。
	// WITHIN 2 指定轮询间隔，单位秒，此处每 2s 询问一次
	// WHERE TargetInstance ISA 'Win32_LogicalDisk' 筛选，TargetInstance 是 __InstanceCreationEvent 的属性，指向被系统动态创建的实例
	// ISA 判断是否为指定的 Win32_LogicalDisk 类型，表示逻辑磁盘
	// AND TargetInstance.DriveType = 2 进一步过滤，TargetInstance.DriveTyp=2 指定可移动磁盘
	// 需要注意 __InstanceCreationEvent 是一个临时事件，只有在系统动态创建对象时，WMI 创建一个事件实例描述新的创建，并不会长时间存在
	queryString := "SELECT * FROM __InstanceCreationEvent WITHIN 2 WHERE TargetInstance ISA 'Win32_LogicalDisk' AND TargetInstance.DriveType = 2"
	// 通过 CallMethod 调用 ExecNotificationQuery 方法
	// ExecNotificationQuery 用于执行一个通知查询，通过指定一个 WQL 语句，订阅对应的事件
	// 返回一个 SWbemEventSource 对象，允许程序处理收到的事件
	resultRaw, err := oleutil.CallMethod(service, "ExecNotificationQuery", queryString)
	if err != nil {
		log.Fatal(err)
	}
	// 转为 IDispatch 对象
	result := resultRaw.ToIDispatch()
	// 延迟释放
	defer result.Release()

	// 永久循环执行事件处理
	fmt.Println("Listening for USB drive insertion events...")
	for {
		handleEvent(result)
	}
}

// handleEvent 插入U盘时的处理逻辑
func handleEvent(result *ole.IDispatch) {
	// 调用 SWbemEventSource 的 NextEvent 方法，获取下一个事件
	// 返回一个事件的 COM 对象
	eventRaw, err := oleutil.CallMethod(result, "NextEvent", nil)
	if err != nil {
		log.Fatal(err)
	}
	event := eventRaw.ToIDispatch()
	defer event.Release()

	// MustGetProperty 用于获取 COM 对象的属性
	// 此处尝试从 event 中获取指定的 TargetInstance 属性
	// 该方法如果获取失败将会引发 panic
	// 如果不能保证一定可以获取到属性，应该考虑使用 oleutil.GetProperty 并适当处理错误
	// 对于诸如 __InstanceCreationEvent、__InstanceDeletionEvent 等事件，TargetInstance 属性通常包含了引发事件的实例
	// 此处，如果创建了新的逻辑磁盘，则会指向 Win32_LogicalDisk 实例
	targetInstance := oleutil.MustGetProperty(event, "TargetInstance")
	instance := targetInstance.ToIDispatch()
	defer instance.Release()

	// 获取 instance 的 DeviceID 属性并转为 String
	deviceId := oleutil.MustGetProperty(instance, "DeviceID").ToString()
	fmt.Printf("USB Drive inserted: %s\n", deviceId)

	sourcePath := deviceId + `\` // Assume the USB is mounted with a drive letter.
	targetPath := `D:\TargetDirectory\`

	// Copy all files and directories from USB drive to target directory
	// filepath.Walk 可以遍历指定目录下的所有文件和目录
	// 第一个参数为要遍历的目录，第二个参数为回调函数，每遍历到一个文件或目录就调用一次
	// 此处使用匿名函数，其功能为将遍历到的文件或目录复制到目标路径下
	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 构建 targetFilePath 作为复制的目标路径
		// TrimPrefix 将 sourcePath 从 path 中去除，只保留 sourcePath 之后的路径字符串，然后拼接到 targetPath
		targetFilePath := filepath.Join(targetPath, strings.TrimPrefix(path, sourcePath))
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
}
