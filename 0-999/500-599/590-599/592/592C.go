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

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t, w, b uint64
	if _, err := fmt.Fscan(in, &t, &w, &b); err != nil {
		return
	}
	if w > b {
		w, b = b, w
	}
	g := gcd(w, b)
	l := b / g
	var L uint64
	if l > 0 && w > t/l {
		L = t + 1
	} else {
		L = l * w
		if L > t {
			L = t + 1
		}
	}
	q := t / L
	r := t % L
	num := q*w + min(r, w-1)
	den := t
	g2 := gcd(num, den)
	fmt.Printf("%d/%d\n", num/g2, den/g2)
}
