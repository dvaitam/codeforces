package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type server struct {
	cap int64
	idx int
}

func attempt(x1, x2 int64, servers []server) (bool, []int, []int) {
	n := len(servers)
	prefix := make([]int, n+1)
	idx := -1
	for i := 1; i <= n; i++ {
		if idx == -1 && int64(servers[i-1].cap)*int64(i) >= x1 {
			idx = i
		}
		prefix[i] = idx
	}
	for t := 1; t <= n; t++ {
		i1 := prefix[t]
		if i1 != -1 && i1 < t {
			j := t - i1
			if int64(servers[t-1].cap)*int64(j) >= x2 {
				a := make([]int, i1)
				b := make([]int, j)
				for k := 0; k < i1; k++ {
					a[k] = servers[k].idx
				}
				for k := 0; k < j; k++ {
					b[k] = servers[i1+k].idx
				}
				return true, a, b
			}
		}
	}
	return false, nil, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var x1, x2 int64
	fmt.Fscan(reader, &n, &x1, &x2)
	servers := make([]server, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &servers[i].cap)
		servers[i].idx = i + 1
	}
	sort.Slice(servers, func(i, j int) bool { return servers[i].cap > servers[j].cap })
	if ok, a, b := attempt(x1, x2, servers); ok {
		fmt.Println("Yes")
		fmt.Println(len(a), len(b))
		for i, v := range a {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(v)
		}
		fmt.Println()
		for i, v := range b {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(v)
		}
		fmt.Println()
		return
	}
	if ok, b, a := attempt(x2, x1, servers); ok {
		fmt.Println("Yes")
		fmt.Println(len(a), len(b))
		for i, v := range a {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(v)
		}
		fmt.Println()
		for i, v := range b {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(v)
		}
		fmt.Println()
		return
	}
	fmt.Println("No")
}
