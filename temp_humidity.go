package main

import (
	"fmt"
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var pin rpio.Pin

var rh1, rh2, temp1, temp2, crc int64 // 传感器数据

var dat int64        // 传感器缓存数据40位
var t1, t2 time.Time // 数据接收计时

var ts = make([]int64, 0)

// 读取传感器数据
// 通过计算时序周期来判断数字量
func readSensorData() bool {
	// 开始初始化请求
	pin.Output()
	pin.Low()
	time.Sleep(time.Millisecond * 25) // >= 20ms
	pin.High()
	pin.Input()

	for pin.Read() == rpio.High {
		time.Sleep(time.Nanosecond * 20)
	} // 等待数据发送的低电平信号
	for pin.Read() == rpio.Low {
		time.Sleep(time.Nanosecond * 20)
	} // 等待总线反馈高电平信号
	for pin.Read() == rpio.High {
		time.Sleep(time.Nanosecond * 20)
	} // 等待数据发送的低电平信号
	for i := 0; i < 40; i++ {
		t1 = time.Now() // 开始1bit数据采集周期
		for pin.Read() == rpio.Low {
		} // 50us低电平
		for pin.Read() == rpio.High {
			if time.Now().UnixNano()-t1.UnixNano() > 500000 {
				break
			}
		} // 26-28us 或者 70us
		t2 := time.Now() // 1bit数据结束时间
		dat = dat * 2
		t := t2.UnixNano() - t1.UnixNano()
		ts = append(ts, t)
		if t > 85000 { // 根据实际情况调整,>85us判断为1
			dat++
		}
	}
	fmt.Println("TS: ", ts)
	ts = ts[:0]
	// 数据校验
	rh1 = (dat >> 32) & 0xFF        // 湿度整数
	rh2 = (dat >> 24) & 0xFF        // 湿度小数
	temp1 = (dat >> 16) & 0xFF      // 温度整数
	temp2 = (dat >> 8) & 0xFF       // 温度小数
	crc = dat & 0xFF                // CRC
	dat = 0                         // 缓存清0
	if crc == rh1+rh2+temp1+temp2 { // CRC校验
		return true
	}
	return false
}

// DHT11温湿度模块
// 通过温湿度模块的信号时序来获取数字信息
func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalln(err)
	}
	defer rpio.Close()

	pin = rpio.Pin(27)
	fmt.Println("--- begin get sensor ---")
	for {
		time.Sleep(time.Second * 2)
		if readSensorData() {
			fmt.Println("--- finish get sensor ---")
			fmt.Println("originial data: ", rh1, rh2, temp1, temp2, crc)
			fmt.Printf("RH:%d.%d\n", rh1, rh2)
			fmt.Printf("TMP:%d.%d\n", temp1, temp2)
		} else {
			fmt.Println("--- failed get sensor ---")
			fmt.Println("originial data: ", rh1, rh2, temp1, temp2, crc)
			fmt.Println("crc failed")
		}
		rh1 = 0
		rh2 = 0
		temp1 = 0
		temp2 = 0
		crc = 0
	}
}
