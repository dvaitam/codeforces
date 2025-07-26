package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(x string) string {
	bs := []byte(x)
	n := len(bs)
	// Look for the rightmost pair with sum >= 10
	for i := n - 1; i >= 1; i-- {
		d1 := int(bs[i-1] - '0')
		d2 := int(bs[i] - '0')
		sum := d1 + d2
		if sum >= 10 {
			res := make([]byte, 0, n)
			res = append(res, bs[:i-1]...)
			res = append(res, byte('0'+sum/10))
			res = append(res, byte('0'+sum%10))
			res = append(res, bs[i+1:]...)
			return string(res)
		}
	}
	// Otherwise use the first pair
	sum := int(bs[0] - '0' + bs[1] - '0')
	res := make([]byte, 0, n-1)
	res = append(res, byte('0'+sum))
	res = append(res, bs[2:]...)
	return string(res)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		fmt.Println(solve(s))
	}
}
