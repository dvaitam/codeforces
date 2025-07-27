package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 10000000

var ans []int32

func precompute() {
	ans = make([]int32, N+1)
	lp := make([]int32, N+1)
	pw := make([]int32, N+1)
	sumP := make([]int32, N+1)
	sigma := make([]int32, N+1)
	primes := make([]int32, 0, 700000)

	pw[1] = 1
	sumP[1] = 1
	sigma[1] = 1
	ans[1] = 1

	for i := 2; i <= N; i++ {
		if lp[i] == 0 {
			lp[i] = int32(i)
			primes = append(primes, int32(i))
			pw[i] = int32(i)
			sumP[i] = 1 + int32(i)
			sigma[i] = sumP[i]
		}
		li := lp[i]
		for _, p32 := range primes {
			p := int(p32)
			if p > int(li) || i*p > N {
				break
			}
			idx := i * p
			lp[idx] = p32
			if p32 == li {
				pw[idx] = pw[i] * p32
				sumP[idx] = sumP[i] + pw[idx]
				sigma[idx] = sigma[i/int(pw[i])] * sumP[idx]
			} else {
				pw[idx] = p32
				sumP[idx] = 1 + p32
				sigma[idx] = sigma[i] * sumP[idx]
			}
		}
		s := sigma[i]
		if int(s) <= N && ans[s] == 0 {
			ans[s] = int32(i)
		}
	}
}

func main() {
	precompute()
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var c int
		fmt.Fscan(reader, &c)
		if ans[c] == 0 {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, ans[c])
		}
	}
}
