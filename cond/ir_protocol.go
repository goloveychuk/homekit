package cond

import "C"

const (
	COLD int64 = 2
	WAVE int64 = 3
	HEAT int64 = 1
	CIRC int64 = 0

	ON  int64 = 1
	OFF int64 = 0

	ONE     C.double = 1400
	ZERO    C.double = 3350
	SPACE   C.double = 640
	HEADER  C.double = 6125
	HEADER2 C.double = 7400
)

func Encode(enabled, mode, temp int64) int64 {
	temp = temp - 17

	var msg int64
	msg += 8 << (10 * 4) //head
	msg += 15 << (9 * 4) //hz
	msg += 7 << (8 * 4)  //fan speed

	var modeOnGroup int64 = (mode << 2) + enabled
	msg += modeOnGroup << (7 * 4)

	var num int64
	num += (temp & 1) << 3
	num += (temp & 2) << 1
	num += (temp & 4) >> 1
	num += (temp & 8) >> 3
	msg += num << (6 * 4)

	msg += ((^modeOnGroup) & 15) << (5 * 4) // reverse
	msg += ((^num) & 15) << (4 * 4)         // reverse

	msg += 10 << (3 * 4)
	msg += 11 << (2 * 4)
	msg += 5 << (1 * 4)
	msg += 4

	// fmt.Println(strconv.FormatInt(msg, 2))
	return msg
}

func Serialize(msg int64) []C.double {
	var str []C.double
	str = append(str, HEADER)
	str = append(str, HEADER2)
	str = append(str, SPACE)
	var p int64 = 1 << 47
	for n := 0; n < 48; n++ {
		if msg&p != 0 {
			str = append(str, ONE)
		} else {
			str = append(str, ZERO)
		}
		str = append(str, SPACE)
		p = p >> 1
	}
	str = append(str, HEADER2)
	str = append(str, SPACE)

	return str

}

// temp - 6 group
// (h)(h2) 0000 1110 1111 0001 1001 0101 0110 1010 1010 1011 0101 0100 (h2)  - on
// (h)(h2) 0000 1110 1111 0001 1000 0101 0111 1010 1010 1011 0101 0100 (h2) - off
// (h)(h2) 0000 1000 1111 0111 1001 0000 0110 1111 1010 1011 0101 0100 (h2) - min temp (cont)
// (h)(h2) 0000 1000 1111 0111 1001 1000 0110 0111 1010 1011 0101 0100 (h2) 18
// (h)(h2) 0000 1000 1111 0111 1001 1011 0110 0100 1010 1011 0101 0100 (h2) 30

// (h)(h2) 0000 1000 1111 0111 1001 1001 0110 0110 1010 1011 0101 0100 (h2) modes
// (h)(h2) 0000 1000 1111 0111 1101 1001 0010 0110 1010 1011 0101 0100 (h2)
// (h)(h2) 0000 1000 1111 0111 0101 1001 1010 0110 1010 1011 0101 0100 (h2)
// (h)(h2) 0000 1000 1111 0111 0001 1000 1110 0111 1010 1011 0101 0100 (h2) - mode

// (h)(h2) 0000 1000 1111 0111 0001 1000 1110 0111 1010 1011 0101 0100 (h2) -on
// (h)(h2) 0000 1000 1111 0111 0000 1000 1111 0111 1010 1011 0101 0100 (h2) -off

// (h)(h2) 0000 1000 1111 0111 1001 0001 0110 1110 1010 1011 0101 0100 (h2)
//         0000 1000 1111 0111 1001 0111 0110 0001 1010 1011 0101 0100
//         0000 1000 1111 0111 1001 0111 0110 0000 1010 1011 0101 0100

// 1000 1111 0111 1001 0111 011000001010101101010100
// 1000 1111 0111 1000 1111 011000001010101101010100
