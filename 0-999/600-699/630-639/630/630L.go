package main

import (
	"fmt"
)

func main() {
	var s string
	if _, err := fmt.Scan(&s); err != nil {
		return
	}
	if len(s) != 5 {
		// ensure 5 digits by padding
		for len(s) < 5 {
			s = "0" + s
		}
	}
	// rearrange digits: first, third, fifth, fourth, second
	rearr := string([]byte{s[0], s[2], s[4], s[3], s[1]})

	// convert to int
	var k int
	fmt.Sscanf(rearr, "%d", &k)
	mod := 100000
	res := 1
	for i := 0; i < 5; i++ {
		res = res * k % mod
	}
	fmt.Printf("%05d", res)
}
