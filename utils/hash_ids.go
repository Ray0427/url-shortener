package utils

import (
	"fmt"

	"github.com/speps/go-hashids/v2"
)

func Encode(salt string, minLength int, num int) string {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{num})
	fmt.Println(e)
	return e
}

func Decode(salt string, minLength int, hash string) (int, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, _ := hashids.NewWithData(hd)
	numbers, err := h.DecodeWithError(hash)
	if err != nil {
		return 0, err
	}
	fmt.Println(numbers)
	return numbers[0], nil
}
