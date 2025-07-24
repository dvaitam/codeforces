package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	const MaxV = 100000
	freq := make([]int, MaxV+1)
	maxV := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		freq[x]++
		if x > maxV {
			maxV = x
		}
	}

	maxLCM := 0
	for g := 1; g <= maxV; g++ {
		first, second := 0, 0
		for m := (maxV / g) * g; m >= g && (first == 0 || second == 0); m -= g {
			if freq[m] == 0 {
				continue
			}
			count := freq[m]
			for c := 0; c < count && (first == 0 || second == 0); c++ {
				if first == 0 {
					first = m
				} else {
					second = m
				}
			}
		}
		if second != 0 {
			l := first / gcd(first, second) * second
			if l > maxLCM {
				maxLCM = l
			}
		}
	}

	fmt.Println(maxLCM)
}
