package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	lcms := []uint64{1}
	for i := 1; ; i++ {
		prev := lcms[len(lcms)-1]
		g := gcd(prev, uint64(i))
		l := prev / g * uint64(i)
		if l > 1e18 {
			break
		}
		lcms = append(lcms, l)
	}

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n uint64
		fmt.Fscan(reader, &n)
		ans := 0
		for i := 1; i < len(lcms); i++ {
			if n%lcms[i] == 0 {
				ans = i
			} else {
				break
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
