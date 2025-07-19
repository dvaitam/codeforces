package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	p   []int
	vec []int
	n   int
)

func getpos(x int) int {
	for i := 1; i <= n; i++ {
		if p[i] == x {
			return i
		}
	}
	panic("position not found")
}

func solve(x int) {
	vec = append(vec, x)
	// reverse prefix 1..x
	for i, j := 1, x; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		fmt.Fscan(reader, &n)
		p = make([]int, n+1)
		ok := true
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &p[i])
			if (p[i]-i)%2 != 0 {
				ok = false
			}
		}
		if !ok {
			fmt.Fprintln(writer, -1)
			continue
		}
		vec = make([]int, 0, n*5)
		for k := n / 2; k >= 1; k-- {
			// place 2*k+1 and 2*k
			pos := getpos(2*k + 1)
			solve(pos)
			pos = getpos(2 * k)
			solve(pos - 1)
			pos = getpos(2 * k)
			solve(pos + 1)
			solve(3)
			solve(2*k + 1)
		}
		// output
		fmt.Fprintln(writer, len(vec))
		for i, v := range vec {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
	}
}
