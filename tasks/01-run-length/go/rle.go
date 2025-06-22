package rle

import (
	"fmt"
	"strings"
)

// Encode returns the run‑length encoding of UTF‑8 string s.
//
// "AAB" → "A2B1"
func Encode(s string) string {
	// if input is empty, then return an empty string
	if s == "" {
		return ""
	}

	//convert an input to a slice of rune to handle multi-byte chars like emoji
	runes := []rune(s)

	//starting count from 1 because every chars is at least appeared once
	count := 1

	// result strings as a buffer
	var result strings.Builder

	// iterate over, starting from the second char to compare with previous one
	for i := 1; i < len(runes); i++ {
		// if the current char is same as previous one, increment the count
		if runes[i] == runes[i-1] {
			count++
		} else {
			// if the current char is different, append the previous character and its count to result
			result.WriteString(fmt.Sprintf("%c%d", runes[i-1], count))
			count = 1 // reset count for the new character
		}
	}

	//return the result with concatenate the last character and its count
	return result.String() + fmt.Sprintf("%c%d", runes[len(runes)-1], count)
}
