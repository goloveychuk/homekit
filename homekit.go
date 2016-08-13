package main

import (
	"fmt"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/log"
	"github.com/goloveychuk/homekit/cond"
)

const (
	MIN  = 18.0
	MAX  = 30.0
	OFF  = 0
	HEAT = 1
	COOL = 2
	AUTO = 3
)

func turnLightOn() {
	log.Println("Turn Light On")
}

func turnLightOff() {
	log.Println("Turn Light Off")
}

func sendState(therm *accessory.Thermostat) {
	temp := therm.Thermostat.TargetTemperature.GetValue()
	if temp > MAX {
		temp = MAX
	}
	if temp < MIN {
		temp = MIN
	}
	coolingState := therm.Thermostat.TargetHeatingCoolingState.GetValue()
	mode := cond.COLD
	enabled := cond.ON
	switch coolingState {
	case OFF:
		enabled = cond.OFF
	case COOL:
		mode = cond.COLD
	case HEAT:
		mode = cond.HEAT
	case AUTO:
		mode = cond.WAVE
	}
	fmt.Println(enabled, mode)
	msg := cond.Encode(enabled, mode, int64(temp))
	resp := cond.Serialize(msg)
	cond.Send(resp)
	therm.Thermostat.CurrentTemperature.SetValue(temp)
	therm.Thermostat.CurrentHeatingCoolingState.SetValue(coolingState)
	fmt.Println(resp)

}
func main() {
	log.Verbose = true
	log.Info = true
	info := accessory.Info{
		Name:         "Air conditioner2",
		Manufacturer: "Matthias",
	}

	acc := accessory.NewThermostat(info, 23, MIN, MAX, 1)

	acc.Thermostat.TargetTemperature.OnValueRemoteUpdate(func(temp float64) {
		sendState(acc)
		log.Println("updating target temp", temp)
	})
	acc.Thermostat.TargetHeatingCoolingState.OnValueRemoteUpdate(func(state int) {
		sendState(acc)
		log.Println("new cooling state", state)
	})

	t, err := hc.NewIPTransport(hc.Config{Pin: "00000000"}, acc.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	hc.OnTermination(func() {
		t.Stop()
	})

	t.Start()
}
