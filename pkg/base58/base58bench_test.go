package base58_test

import (
	"strconv"
	"testing"
	"tinyurl/pkg/base58"

	btcbase58 "github.com/btcsuite/btcd/btcutil/base58"
)

/*
Result:
goos: darwin
goarch: arm64
pkg: tinyurl/pkg/base58
BenchmarkBase58Decode
BenchmarkBase58Decode-8          	25230484	        46.79 ns/op
BenchmarkBase58BtcutilDecode
BenchmarkBase58BtcutilDecode-8   	19115382	        61.88 ns/op
BenchmarkBase58Encode
BenchmarkBase58Encode-8          	20997159	        56.93 ns/op
BenchmarkBase58BtcutilEncode
BenchmarkBase58BtcutilEncode-8   	15248187	        78.08 ns/op
PASS
ok  	tinyurl/pkg/base58	6.072s

如果我們只針對 Number 去做 Encode 那可以用自己寫的
*/
func BenchmarkBase58Decode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		base58.DecodeToInt("2XNGAK")
	}
}

func BenchmarkBase58BtcutilDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		btcbase58.Decode("2XNGAK")
	}
}

func BenchmarkBase58Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		base58.EncodeFromInt(825241648)
	}
}

func BenchmarkBase58BtcutilEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		btcbase58.Encode([]byte(strconv.Itoa(1000)))
	}
}
