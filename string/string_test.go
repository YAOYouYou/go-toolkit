package string

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByte2String(t *testing.T) {
	b := []byte{'1', '2', '3', '4'}
	s := "1234"
	assert.Equal(t, s, string(b))
}

func TestString2Byte(t *testing.T) {
	s := "1234"
	b := String2Bytes(s)
	assert.Equal(t, b, []byte{'1', '2', '3', '4'})
}

func BenchmarkByte2String(b *testing.B) {
	// b := []byte("jin tian shi zhou ri. cha bu duo ke yi shui jiao")
	x := []byte("hahahahhahhahahahahaha")
	for i := 0; i < b.N; i++ {
		Bytes2String(x)
	}
}

func BenchmarkNormalByte2String(b *testing.B) {
	// b := []byte("jin tian shi zhou ri. cha bu duo ke yi shui jiao")
	x := []byte("hahahahhahhahahahahaha")
	for i := 0; i < b.N; i++ {
		_ = string(x)
	}
}

func BenchmarkString2Bytes(b *testing.B) {
	s := "hahahahhahhahahahahaha"
	for i:=0;i<b.N;i++{
		String2Bytes(s)
	}
}

func BenchmarkNormalString2Bytes(b *testing.B) {
	s := "hahahahhahhahahahahaha"
	for i:=0;i<b.N;i++{
		_ = []byte(s)
	}
}