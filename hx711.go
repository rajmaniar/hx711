package hx711

import (
	"github.com/mrmorphic/hwio"
)

const (
	GAIN_A_128 = 1
	GAIN_B_32  = 2
	GAIN_A_64  = 3
)

type HX711 struct {
	Clock   string
	Data    string
	Gain    int
	clkPin  hwio.Pin
	dataPin hwio.Pin
}

//New instantiates a new object
func New(data string, clock string) (*HX711, error) {
	var err error
	var clkPin, dataPin hwio.Pin
	if clkPin, err = hwio.GetPin(clock); err != nil {
		return &HX711{}, err
	}
	if dataPin, err = hwio.GetPin(data); err != nil {
		return &HX711{}, err
	}
	if err = hwio.PinMode(dataPin, hwio.INPUT); err != nil {
		return &HX711{}, err
	}
	if err = hwio.PinMode(clkPin, hwio.OUTPUT); err != nil {
		return &HX711{}, err
	}
	return &HX711{Clock: clock, Data: data, Gain: GAIN_A_128, clkPin: clkPin, dataPin: dataPin}, err
}

//OnReady Blocks until the chip is ready to send data
func (h *HX711) OnReady() error {
	if err := h.clockLow(); err != nil {
		return err
	}
	ready := false
	for !ready {
		r, err := h.readBit()
		if err != nil {
			return err
		}
		if r == hwio.LOW {
			ready = true
		}
	}
	return nil
}

//Sleep the chip until you need
func (h *HX711) Sleep() error {
	return h.clockHigh()
}

//Reset / Wakeup  the chip
func (h *HX711) Reset() error {
	if err := h.clockHigh(); err != nil {
		return err
	}
	hwio.DelayMicroseconds(60)
	return h.clockLow()
}

//SetGain sets the gain after data has been read
func (h *HX711) SetGain() error {
	for i := 0; i < h.Gain; i++ {
		if err := h.tick(); err != nil {
			return err
		}
	}
	return nil
}

//ReadData gets a 24bit signed int from the chip
func (h *HX711) ReadData() (int32, error) {
	c := int32(0)
	if err := h.OnReady(); err != nil {
		return 0, err
	}
	for i := 0; i < 24; i++ {
		h.tick()
		b, err := hwio.DigitalRead(h.dataPin)
		if err != nil {
			return 0, err
		}
		c = c << 1
		if b == hwio.HIGH {
			c++
		}
	}

	return twosComp(c), h.SetGain()
}

func (h *HX711) tick() error {
	err := h.clockHigh()
	if err != nil {
		return err
	}
	return h.clockLow()
}

func (h *HX711) clockHigh() error {
	return hwio.DigitalWrite(h.clkPin, hwio.HIGH)
}

func (h *HX711) clockLow() error {
	return hwio.DigitalWrite(h.clkPin, hwio.LOW)
}

func (h *HX711) readBit() (int, error) {
	return hwio.DigitalRead(h.dataPin)
}

func twosComp(i int32) int32 {
	if (i & 0x800000) > 0 {
		i |= ^0xffffff
	}
	return i
}
