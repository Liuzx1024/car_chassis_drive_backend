package serial

import "time"

func posixTimeoutValues(readTimeout time.Duration) (vmin uint8, vtime uint8) {
	const MAXUINT8 = 1<<8 - 1
	var minBytesToRead uint8 = 1
	var readTimeoutInDeci int64
	if readTimeout > 0 {
		minBytesToRead = 0
		readTimeoutInDeci = (readTimeout.Nanoseconds() / 1e6 / 100)
		if readTimeoutInDeci < 1 {
			readTimeoutInDeci = 1
		} else if readTimeoutInDeci > MAXUINT8 {
			readTimeoutInDeci = MAXUINT8
		}
	}
	return minBytesToRead, uint8(readTimeoutInDeci)
}
