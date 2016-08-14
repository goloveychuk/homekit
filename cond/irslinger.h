#ifndef IRSLINGER_H
#define IRSLINGER_H

#include <string.h>
#include <math.h>
#include <pigpio.h>
#include <stdio.h>

#define MAX_COMMAND_SIZE 5120
#define MAX_PULSES 120000

static inline void addPulse(uint32_t onPins, uint32_t offPins, uint32_t duration, gpioPulse_t *irSignal, int *pulseCount)
{
	int index = *pulseCount;

	irSignal[index].gpioOn = onPins;
	irSignal[index].gpioOff = offPins;
	irSignal[index].usDelay = duration;

	(*pulseCount)++;
}

// Generates a square wave for duration (microseconds) at frequency (Hz)
// on GPIO pin outPin. dutyCycle is a floating value between 0 and 1.
static inline void carrierFrequency(uint32_t outPin, double frequency, double dutyCycle, double duration, gpioPulse_t *irSignal, int *pulseCount)
{
    printf("pulse %f\n", duration);
	double oneCycleTime = 1000000.0 / frequency; // 1000000 microseconds in a second
	int onDuration = (int)round(oneCycleTime * dutyCycle);
	int offDuration = (int)round(oneCycleTime * (1.0 - dutyCycle));

	int totalCycles = (int)round(duration / oneCycleTime);
	int totalPulses = totalCycles * 2;
    
	int i;
	for (i = 0; i < totalPulses; i++)
	{
		if (i % 2 == 0)
		{
			// High pulse
			addPulse(1 << outPin, 0, onDuration, irSignal, pulseCount);
		}
		else
		{
			// Low pulse
			addPulse(0, 1 << outPin, offDuration, irSignal, pulseCount);
		}
	}
}

// Generates a low signal gap for duration, in microseconds, on GPIO pin outPin
static inline void gap(uint32_t outPin, double duration, gpioPulse_t *irSignal, int *pulseCount)
{
    printf("space %f\n", duration);    
	addPulse(0, 0, duration, irSignal, pulseCount);
}

static inline int irSling(uint32_t outPin,
	int frequency,
	double dutyCycle,
	double *codes, int len){
	if (outPin > 31)
	{
		// Invalid pin number
		return 1;
	}
    

	// // printf("code size is %zu\n", codeLen);

	// if (codeLen > MAX_COMMAND_SIZE)
	// {
	// 	// Command is too big
	// 	return 1;
	// }

	gpioPulse_t irSignal[MAX_PULSES];
	int pulseCount = 0;

	// Generate Code
	

	int flag = 1;
    int i;
    for (i = 0; i<len; i++) {
        int dur = codes[i]; 
        if (flag) {
        	carrierFrequency(outPin, frequency, dutyCycle, dur, irSignal, &pulseCount);
            flag = 0;    
        } else {
        	gap(outPin, dur, irSignal, &pulseCount);
            flag = 1;
        }
    }
	
    
	printf("pulse count is %i\n", pulseCount);
	// End Generate Code

	// Init pigpio
	if (gpioInitialise() < 0)
	{
		// Initialization failed
		printf("GPIO Initialization failed\n");
		return 1;
	}

	// Setup the GPIO pin as an output pin
	gpioSetMode(outPin, PI_OUTPUT);

	// Start a new wave
	gpioWaveClear();

	gpioWaveAddGeneric(pulseCount, irSignal);
	int waveID = gpioWaveCreate();

	if (waveID >= 0)
	{
		int result = gpioWaveTxSend(waveID, PI_WAVE_MODE_ONE_SHOT);

		printf("Result: %i\n", result);
	}
	else
	{
		printf("Wave creation failure!\n %i", waveID);
	}

	// Wait for the wave to finish transmitting
	while (gpioWaveTxBusy())
	{
		time_sleep(0.1);
	}

	// Delete the wave if it exists
	if (waveID >= 0)
	{
		gpioWaveDelete(waveID);
	}

	// Cleanup
	gpioTerminate();
	return 0;
}

#endif