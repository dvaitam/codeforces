package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var k int
	var m int
	fmt.Fscan(in, &n, &k, &m)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	// build full sequence
	seq := make([]int, 0, n*m)
	for i := 0; i < m; i++ {
		seq = append(seq, arr...)
	}
	// stack of pairs value,count
	type pair struct {
		val int
		cnt int
	}
	stack := []pair{}
	for _, v := range seq {
		if len(stack) > 0 && stack[len(stack)-1].val == v {
			stack[len(stack)-1].cnt++
			if stack[len(stack)-1].cnt == k {
				stack = stack[:len(stack)-1]
			}
		} else {
			stack = append(stack, pair{val: v, cnt: 1})
		}
	}
	// compute remaining length
	res := 0
	for _, p := range stack {
		res += p.cnt
	}
	fmt.Println(res)
}
