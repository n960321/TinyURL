package base58

import (
	"bytes"
	"errors"
)

var charList = [58]byte{
	'1', '2', '3', '4',
	'5', '6', '7', '8',
	'9',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
}

var charMap = map[byte]float64{
	'1': 0, '2': 1, '3': 2, '4': 3,
	'5': 4, '6': 5, '7': 6, '8': 7,
	'9': 8,
	'A': 9, 'B': 10, 'C': 11, 'D': 12, 'E': 13, 'F': 14, 'G': 15, 'H': 16, 'J': 17, 'K': 18, 'L': 19, 'M': 20, 'N': 21, 'P': 22, 'Q': 23, 'R': 24, 'S': 25, 'T': 26, 'U': 27, 'V': 28, 'W': 29, 'X': 30, 'Y': 31, 'Z': 32,
	'a': 33, 'b': 34, 'c': 35, 'd': 36, 'e': 37, 'f': 38, 'g': 39, 'h': 40, 'i': 41, 'j': 42, 'k': 43, 'm': 44, 'n': 45, 'o': 46, 'p': 47, 'q': 48, 'r': 49, 's': 50, 't': 51, 'u': 52, 'v': 53, 'w': 54, 'x': 55, 'y': 56, 'z': 57,
}

var powList = [11]float64{
	1,
	58,
	58 * 58,
	58 * 58 * 58,
	58 * 58 * 58 * 58,
	58 * 58 * 58 * 58 * 58,
	58 * 58 * 58 * 58 * 58 * 58,
	58 * 58 * 58 * 58 * 58 * 58 * 58,
	58 * 58 * 58 * 58 * 58 * 58 * 58 * 58,
	58 * 58 * 58 * 58 * 58 * 58 * 58 * 58 * 58,
	58 * 58 * 58 * 58 * 58 * 58 * 58 * 58 * 58 * 58,
}

func EncodeFromInt(i int) string {
	var buf bytes.Buffer
	l := make([]byte, 0)
	for i > 0 {
		r := i % 58
		l = append(l, charList[r])
		i /= 58
	}

	for i := len(l) - 1; i >= 0; i-- {
		buf.WriteByte(l[i])
	}

	return buf.String()
}

func DecodeToInt(s string) (int, error) {
	var r float64
	for i := 0; i < len(s); i++ {
		if _, exist := charMap[s[i]]; !exist {
			return 0, errors.New("the string can not decode")
		}
		r += powList[len(s)-i-1] * charMap[s[i]]
	}
	return int(r), nil
}
