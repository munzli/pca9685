package pca9685

import (
	"log"
	"math"
	"time"

	"github.com/go-daq/smbus"
)

const (
	ADDR           uint8 = 0x40
	// registers/etc.
	__SUBADR1      uint8 = 0x02
	__SUBADR2      uint8 = 0x03
	__SUBADR3      uint8 = 0x04
	__MODE1        uint8 = 0x00
	__PRESCALE     uint8 = 0xFE
	__LED0_ON_L    uint8 = 0x06
	__LED0_ON_H    uint8 = 0x07
	__LED0_OFF_L   uint8 = 0x08
	__LED0_OFF_H   uint8 = 0x09
	__ALLLED_ON_L  uint8 = 0xFA
	__ALLLED_ON_H  uint8 = 0xFB
	__ALLLED_OFF_L uint8 = 0xFC
	__ALLLED_OFF_H uint8 = 0xFD
)

// device is a PCA9685 based device.
type Device struct {
	conn *smbus.Conn // connection to smbus
	addr uint8       // address
}

func Open() (*Device, error) {
	c, err := smbus.Open(1, ADDR)
	if err != nil {
		log.Fatalf("open error: %v\n", err)
	}

	dev := Device{
		conn: c,
		addr: ADDR,
	}

	time.Sleep(50 * time.Millisecond) // wait required time
	return &dev, err
}

func (dev *Device) Write(reg, value uint8) {
	// writes an 8-bit value to the specified register/address
	_ = dev.conn.WriteReg(dev.addr, reg, value)
	log.Printf("I2C: write %d to register %d\n", value, reg)
}

func (dev *Device) Read(reg uint8) uint8 {
	// read an unsigned byte from the I2C device
	result, _ := dev.conn.ReadReg(dev.addr, reg)
	log.Printf("I2C: device %d returned %d from reg %d", dev.addr, result & 0xFF, reg)
	return result
}

func (dev *Device) Close() error {
	return dev.conn.Close()
}

func (dev *Device) SetPWMFrequency(frequency int) {
	prescaleval := 25000000.0 // 25MHz
	prescaleval /= 4096.0     // 12-bit
	prescaleval /= float64(frequency)
	prescaleval -= 1.0

	log.Printf("setting PWM frequency to %d Hz\n", frequency)
	log.Printf("estimated pre-scale: %f\n", prescaleval)

	prescale := math.Floor(prescaleval + 0.5)

	log.Printf("final pre-scale: %f\n", prescale)

	oldmode := dev.Read(__MODE1)
	newmode := (oldmode & 0x7F) | 0x10 // sleep
	dev.Write(__MODE1, newmode)  // go to sleep
	dev.Write(__PRESCALE, uint8(math.Floor(prescale)))
	dev.Write(__MODE1, oldmode)
	time.Sleep(5)
	dev.Write(__MODE1, oldmode | 0x80)
}

func (dev *Device) setPWM(channel uint8, on, off int) {
	// sets a single PWM channel
	dev.Write(__LED0_ON_L + 4 * channel, uint8(on & 0xFF))
	dev.Write(__LED0_ON_H + 4 * channel, uint8(on >> 8))
	dev.Write(__LED0_OFF_L + 4 * channel, uint8(off & 0xFF))
	dev.Write(__LED0_OFF_H + 4 * channel, uint8(off >> 8))

	log.Printf("channel: %d  LED_ON: %d LED_OFF: %d\n", channel, on, off)
}

func (dev *Device) SetServoPulse(channel uint8, pulse int) {
	// sets the servo pulse
	pulse = pulse * 4096 / 20000 // PWM frequency is 50HZ, the period is 20000us
	dev.setPWM(channel, 0 , pulse)
}