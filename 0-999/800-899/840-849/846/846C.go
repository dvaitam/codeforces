package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var N int
	if _, err := fmt.Fscan(reader, &N); err != nil {
		return
	}
	a := make([]int64, N+1)
	for i := 1; i <= N; i++ {
		fmt.Fscan(reader, &a[i])
	}
	s := make([]int64, N+1)
	for i := 1; i <= N; i++ {
		s[i] = s[i-1] + a[i]
	}

	var sum, now int64
	// initial now = 0 for empty segment
	var ans int64 = -5000000000000000000
	st := 1
	St, Ed := 1, 0
	var m0, m1, m2 int
	for i := 1; i <= N; i++ {
		sum += a[i]
		if sum >= 0 {
			sum = 0
			st = i + 1
		} else {
			// potential minimal segment [st..i]
			if sum < now {
				now = sum
				St = st
				Ed = i
			}
		}
		// candidate with delim2 = i
		// res = (s[i] - now)*2 - s[N]
		val := (s[i]-now)*2 - s[N]
		if val > ans {
			ans = val
			m0 = St - 1
			m1 = Ed
			m2 = i
		}
	}
	fmt.Fprintln(writer, m0, m1, m2)
}
