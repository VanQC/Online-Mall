package tool

import (
	"strconv"
	"testing"
)

func TestEncryption_AesDecoding(t *testing.T) {
	en := Encryption{key: "Ek1+Ep1==Ek2+Ep2"}
	m, err := strconv.ParseFloat("", 64)
	if err != nil {
		t.Fatal("not allow null")
	}
	if m != float64(0) {
		t.Fatal("null not 0")
	}

	v := en.AesEncoding("500")
	res := en.AesDecoding(v)
	if res != "500" {
		t.Fatal(res)
	}
}

func TestParse(t *testing.T) {
	m, err := strconv.ParseFloat("0", 64)
	if err != nil {
		t.Fatal("not allow null")
	}
	if m != float64(0) {
		t.Fatal("null not 0")
	}
}
