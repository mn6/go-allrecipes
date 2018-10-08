package main

import "strconv"

// RoundToUint8 rounds string float to uint8
func RoundToUint8(val string) uint8 {
	new, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}
	if new < 0 {
		return uint8(new - 0.5)
	}
	return uint8(new + 0.5)
}

// RoundToUint16 rounds string float to uint8
func RoundToUint16(val string) uint16 {
	new, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}
	if new < 0 {
		return uint16(new - 0.5)
	}
	return uint16(new + 0.5)
}
