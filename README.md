# Reading from an hx711 24-Bit ADC on a Raspberry Pi in Golang
## 24-Bit Analog-to-Digital Converter (ADC) for Weigh Scales

This is a simple package to read the value from the HX711 load cell amplifier like this one from [sparkfun](https://www.sparkfun.com/products/13879)

It uses the primitives from [HWIO](https://github.com/mrmorphic/hwio) and the protocol from the [HX711 data sheet](https://cdn.sparkfun.com/datasheets/Sensors/ForceFlex/hx711_english.pdf)

I've only tested this on a rPi Zero ([GPIO at pin 16 & 18](http://pinout.xyz/pinout/pin16_gpio23)) and it's (obviously) missing a lot of high-level things are tare, calibrate, etc. 

Feel free to submit pull requests with new features.


### Example:

```
import "github.com/rajmaniar/hx711"

func main() {

   	clock := "gpio23"
   	data := "gpio24"
   
   	h,err := hx711.New(data,clock)
   
   	if err != nil {
   		fmt.Printf("Error: %v",err)
   	}
   
   	for err == nil {
   		var data int32
   		data, err = h.ReadData()
   		fmt.Printf("Read from HX711: %v\n",data)
   		time.Sleep(250 * time.Millisecond)
   	}
   	fmt.Printf("Stopped reading because of: %v\n",err)
}

```


### NB
* `h.Reset()` will reset the chip
* `h.Gain` is set to `hx711.GAIN_A_128` by default