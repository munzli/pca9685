package main

import (
	"github.com/munzli/pca9685"
	"time"
)

func main() {
	device, _ := pca9685.Open()

	device.SetPWMFrequency(50)

	for i := 800; i < 2000; i+=5 {
		device.SetServoPulse(2,i)
		time.Sleep(2 * time.Millisecond)
	}
	for i := 2000; i > 800; i-=5 {
		device.SetServoPulse(2,i)
		time.Sleep(2 * time.Millisecond)
	}

	device.Close()
}
