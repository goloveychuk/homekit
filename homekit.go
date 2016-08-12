package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/log"
)

func turnLightOn() {
	log.Println("Turn Light On")
}

func turnLightOff() {
	log.Println("Turn Light Off")
}

func getState(therm *accessory.Thermostat) int64 {
	val := therm.Thermostat.CurrentHeatingCoolingState.GetValue()
	switch val {
	case 1:
		return COLD
	}
	return -1
}
func main() {
	log.Verbose = true
	log.Info = true
	info := accessory.Info{
		Name:         "Air conditioner2",
		Manufacturer: "Matthias",
	}
	min := 18.0
	max := 30.0
	acc := accessory.NewThermostat(info, 23, min, max, 1)

	acc.Thermostat.TargetTemperature.OnValueRemoteUpdate(func(temp float64) {
		if temp > max {
			temp = max
		}
		if temp < min {
			temp = min
		}
		state := getState(acc)
		msg := encode(ON, state, int64(temp))
		resp := serialize(msg)
		log.Println("updating target temp", resp)
		log.Println("updating target temp", temp)
	})
	acc.Thermostat.TargetHeatingCoolingState.OnValueRemoteUpdate(func(state int) {
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
