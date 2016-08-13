#include <stdio.h>
#include "irslinger.h"

int main(int argc, char *argv[])
{
	uint32_t outPin = 22;            // The Broadcom pin number the signal will be sent on
	int frequency = 38000;           // The frequency of the IR signal in Hz
	double dutyCycle = 0.5;          // The duty cycle of the IR signal. 0.5 means for every cycle,
	                                 // the LED will turn on for half the cycle time, and off the other half
    double arr[101] = {6125, 7400, 640, 3350, 640, 3350, 640, 3350, 640, 3350, 640, 1400, 640, 3350, 640, 3350, 640, 3350, 640, 1400, 640, 1400, 640, 1400, 640, 1400, 640, 3350, 640, 1400, 640, 1400, 640, 1400, 640, 1400, 640, 3350, 640, 3350, 640, 1400, 640, 1400, 640, 1400, 640, 3350, 640, 1400, 640, 3350, 640, 1400, 640, 1400, 640, 3350, 640, 3350, 640, 3350, 640, 1400, 640, 3350, 640, 1400, 640, 3350, 640, 1400, 640, 3350, 640, 1400, 640, 3350, 640, 1400, 640, 1400, 640, 3350, 640, 1400, 640, 3350, 640, 1400, 640, 3350, 640, 1400, 640, 3350, 640, 3350, 640, 7400, 640};

	int result = irSling(
        outPin, frequency, dutyCycle,
        &arr, 101    );
	
	return result;
}