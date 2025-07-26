package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func minimalShift(chars []byte) int {
	m := len(chars)
	for d := 1; d <= m; d++ {
		if m%d != 0 {
			continue
		}
		ok := true
		for i := 0; i < m; i++ {
			if chars[i] != chars[(i+d)%m] {
				ok = false
				break
			}
		}
		if ok {
			return d
		}
	}
	return m
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
			p[i]--
		}
		visited := make([]bool, n)
		var ans int64 = 1
		for i := 0; i < n; i++ {
			if visited[i] {
				continue
			}
			cycle := []int{}
			j := i
			for !visited[j] {
				visited[j] = true
				cycle = append(cycle, j)
				j = p[j]
			}
			m := len(cycle)
			chars := make([]byte, m)
			for idx := 0; idx < m; idx++ {
				chars[idx] = s[cycle[idx]]
			}
			d := minimalShift(chars)
			ans = lcm(ans, int64(d))
		}
		fmt.Fprintln(writer, ans)
	}
}
