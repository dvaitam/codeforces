package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	N    int
	K    int
	A    []int64
	B    []int64
	cand []int
)

func hoge(v int64) bool {
	var num int64
	// compute required operations
	for i := 0; i < N; i++ {
		x := A[i] - int64(K)*B[i]
		if x < 0 {
			if v == 0 {
				return false
			}
			x = -x
			num += (x + v - 1) / v
			if num > int64(K) {
				return false
			}
		}
	}
	// schedule feasibility
	for i := 0; i <= K; i++ {
		cand[i] = 0
	}
	for i := 0; i < N; i++ {
		cur := A[i]
		// while still below target B[i]*K
		for cur-B[i]*int64(K) < 0 {
			// earliest time slot
			ng := int(cur/B[i] + 1)
			if ng > K {
				break
			}
			cand[ng]++
			cur += v
		}
	}
	var sum int
	for i := 1; i <= K; i++ {
		sum += cand[i]
		if sum > i {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fscan(reader, &N, &K)
	// reduce operations count
	K--
	A = make([]int64, N)
	B = make([]int64, N)
	cand = make([]int, K+1)
	for i := 0; i < N; i++ {
		fmt.Fscan(reader, &A[i])
	}
	for i := 0; i < N; i++ {
		fmt.Fscan(reader, &B[i])
	}
	// binary search on v
	var ret int64 = (int64(1) << 43) - 1
	if !hoge(ret) {
		fmt.Fprintln(writer, -1)
		return
	}
	for i := 42; i >= 0; i-- {
		d := int64(1) << i
		if ret > d && hoge(ret-d) {
			ret -= d
		}
	}
	fmt.Fprintln(writer, ret)
}
