package main

import (
	"fmt"
	"strconv"

	"github.com/stianeikeland/go-rpio"
)

// 通过输入数值0-255来控制LED亮度
// 使用pwm来控制
func main() {
	pin.Mode(rpio.Pwm)
	pin.Freq(64000)
	pin.DutyCycle(0, 255)
	var value string
	for {
		fmt.Printf("Please enter value 0 - 255: ")
		fmt.Scanln(&value) //Scanln 扫描来自标准输入的文本，将空格分隔的值依次存放到后续的参数内，直到碰到换行。
		d, _ := strconv.Atoi(value)
		pin.DutyCycle(uint32(d), 255)
		fmt.Println("set value: ", d)
	}
}
