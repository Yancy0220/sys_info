package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
)

/**
@author: yangkexian
@date: 2023/4/7 14:57
*/

// 获取CPU
func GetCpuInfo(c *fiber.Ctx) error {
	cpuPercent, _ := cpu.Percent(time.Second, true)
	//fmt.Printf("CPU使用率: %.3f%% \n", cpuPercent[0])
	cpuNumber, _ := cpu.Counts(true)
	//fmt.Printf("CPU核心数: %v \n", cpuNumber)
	return c.JSON(fiber.Map{
		"msg":        "cpuPercent:CPU使用率,cpuNumber:CPU核心数",
		"cpuPercent": fmt.Sprintf("%.3f%% ", cpuPercent[0]),
		"cpuNumber":  cpuNumber,
	})
}

// 获取内存
func GetMemInfo(c *fiber.Ctx) error {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("get memory info fail. err： ", err)
	}
	// 获取总内存大小，单位GB
	memTotal := memInfo.Total / 1024 / 1024
	// 获取已用内存大小，单位MB
	memUsed := memInfo.Used / 1024 / 1024
	// 可用内存大小
	memAva := memInfo.Available / 1024 / 1024
	// 内存可用率
	memUsedPercent := memInfo.UsedPercent
	//fmt.Printf("总内存: %v GB, 已用内存: %v MB, 可用内存: %v MB, 内存使用率: %.3f %% \n",memTotal,memUsed,memAva,memUsedPercent)
	return c.JSON(fiber.Map{
		"msg":            "memTotal:总内存,memUsed:已用内存,memAva:可用内存,memUsedPercent:内存使用率",
		"memTotal":       memTotal,
		"memUsed":        memUsed,
		"memAva":         memAva,
		"memUsedPercent": fmt.Sprintf("%.3f %%", memUsedPercent),
	})
}

// 获取系统平均负载
func GetSysLoad(c *fiber.Ctx) error {
	loadInfo, err := load.Avg()
	if err != nil {
		fmt.Println("get average load fail. err: ", err)
	}
	//fmt.Printf("系统平均负载: %v \n",loadInfo)
	return c.JSON(fiber.Map{
		"msg":      "loadInfo:系统平均负载",
		"loadInfo": loadInfo,
	})
}

type Data struct {
	Path string `json:"path" xml:"path" form:"path"`
}

// 获取硬盘
func GetDiskInfo(c *fiber.Ctx) error {
	p := new(Data)
	if err := c.BodyParser(p); err != nil {
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	if p.Path == "" {
		return c.JSON(fiber.Map{
			"msg": "path不能为空",
		})
	}
	diskPart, err := disk.Partitions(false)
	var info map[string]string
	if err != nil {
		return c.JSON(fiber.Map{
			"msg":      "Total:分区总大小,UsedPercent:分区使用率,InodesUsedPercent:分区inode使用率",
			"loadInfo": info,
		})
	}
	for _, dp := range diskPart {
		fmt.Println(dp.Device)
		if dp.Device == p.Path {
			diskUsed, _ := disk.Usage(dp.Mountpoint)
			//fmt.Printf("分区总大小: %d GB \n", diskUsed.Total/1024/1024/1024)
			//fmt.Printf("已用分区: %d GB \n", diskUsed.Used/1024/1024/1024)
			//fmt.Printf("剩余分区: %d GB \n", diskUsed.Free/1024/1024/1024)
			//fmt.Printf("分区使用率: %.3f %% \n", diskUsed.UsedPercent)
			return c.JSON(fiber.Map{
				"msg":         "Total:分区总大小,UsedPercent:分区使用率,Used:已用分区,Free:剩余分区",
				"Total":       fmt.Sprintf("%d GB", diskUsed.Total/1024/1024/1024),
				"UsedPercent": fmt.Sprintf("%.3f %%", diskUsed.UsedPercent),
				"Used":        fmt.Sprintf("%d GB", diskUsed.Used/1024/1024/1024),
				"Free":        fmt.Sprintf("%d GB", diskUsed.Free/1024/1024/1024),
			})
		}
	}
	return c.JSON(fiber.Map{
		"msg":      "Total:分区总大小,UsedPercent:分区使用率,InodesUsedPercent:分区inode使用率",
		"loadInfo": info,
	})

}
