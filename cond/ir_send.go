package cond

/*
#cgo LDFLAGS: -lpigpio -lm -pthread -lrt
#include "irslinger.h"
*/
import "C"

// ONE     string = "1400 "
// ZERO    string = "3350 "
// SPACE   string = "640 "
// HEADER  string = "6125 "
// HEADER2 string = "7400 "

func Send(d []C.double) {
	var outPin C.uint32_t = 22   // The Broadcom pin number the signal will be sent on
	var frequency C.int = 38000  // The frequency of the IR signal in Hz
	var dutyCycle C.double = 0.5 // The duty cycle of the IR signal. 0.5 means for every cycle,
	// the LED will turn on for half the cycle time, and off the other half
	// var leadingPulseDuration C.int = 6125 // The duration of the beginning pulse in microseconds
	// var leadingGapDuration C.int = 7400   // The duration of the gap in microseconds after the leading pulse
	// var onePulse C.int = 1400             // The duration of a pulse in microseconds when sending a logical 1
	// var zeroPulse C.int = 3350            // The duration of a pulse in microseconds when sending a logical 0
	// var oneGap C.int = 640                // The duration of the gap in microseconds when sending a logical 1
	// var zeroGap C.int = 640               // The duration of the gap in microseconds when sending a logical 0
	// var sendTrailingPulse C.int = 1       // 1 = Send a trailing pulse with duration equal to "onePulse"
	// 0 = Don't send a trailing pulse
	// str := C.CString(d)
	// d2 := unsafe.Pointer(&d[0])

	C.irSling(
		outPin,
		frequency,
		dutyCycle,
		&d[0], C.int(len(d)))

}
