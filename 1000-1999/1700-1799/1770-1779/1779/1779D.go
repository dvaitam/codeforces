package main

import (
	"bufio"
	"fmt"
	"os"
)

type node struct {
	val  int
	need bool
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		var m int
		fmt.Fscan(in, &m)
		cnt := make(map[int]int)
		for i := 0; i < m; i++ {
			var x int
			fmt.Fscan(in, &x)
			cnt[x]++
		}
		ok := true
		for i := 0; i < n; i++ {
			if a[i] < b[i] {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Fprintln(out, "NO")
			continue
		}
		need := make(map[int]int)
		stack := make([]node, 0)
		for i := 0; i < n; i++ {
			bi := b[i]
			diff := a[i] > b[i]
			for len(stack) > 0 && stack[len(stack)-1].val < bi {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top.need {
					need[top.val]++
				}
			}
			if len(stack) > 0 && stack[len(stack)-1].val == bi {
				if diff {
					stack[len(stack)-1].need = true
				}
			} else {
				stack = append(stack, node{val: bi, need: diff})
			}
		}
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if top.need {
				need[top.val]++
			}
		}
		for x, v := range need {
			if cnt[x] < v {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
