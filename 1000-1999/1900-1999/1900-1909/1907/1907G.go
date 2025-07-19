package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
testsLoop:
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		// read state array
		s := make([]bool, n)
		line, _ := reader.ReadString('\n')
		if len(line) == 0 || line[0] == '\n' {
			line, _ = reader.ReadString('\n')
		}
		fields := strings.Fields(line)
		if len(fields) == 1 && len(fields[0]) == n {
			for i, ch := range fields[0] {
				s[i] = (ch == '1')
			}
		} else {
			for i := 0; i < n; i++ {
				s[i] = (fields[i] == "1")
			}
		}
		// read directs
		direct := make([]int, n)
		deg := make([]int, n)
		for i := 0; i < n; i++ {
			var k int
			fmt.Fscan(reader, &k)
			k--
			direct[i] = k
			deg[k]++
		}
		// process acyclic nodes
		used := make([]bool, n)
		ans := make([]int, 0, n)
		stack := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if deg[i] == 0 {
				stack = append(stack, i)
			}
		}
		for len(stack) > 0 {
			f := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			x := direct[f]
			used[f] = true
			if s[f] {
				ans = append(ans, f)
				s[f] = false
				s[x] = !s[x]
			}
			deg[x]--
			if deg[x] == 0 {
				stack = append(stack, x)
			}
		}
		// process cycles
		s1 := make([]bool, n)
		copy(s1, s)
		for i := 0; i < n; i++ {
			if used[i] {
				continue
			}
			// find cycle starting at i
			f0 := i
			tCount := 0
			cur := i
			e := -1
			for {
				used[cur] = true
				x := direct[cur]
				if s[cur] && e == -1 {
					e = cur
				}
				if s[cur] {
					tCount++
				}
				if x == f0 {
					break
				}
				cur = x
			}
			f := e
			if tCount%2 == 1 {
				fmt.Fprintln(writer, -1)
				continue testsLoop
			}
			if f == -1 {
				continue
			}
			// two strategies
			ans1 := make([]int, 0)
			cur = f
			for {
				x := direct[cur]
				if s[cur] {
					s[cur] = !s[cur]
					s[x] = !s[x]
					ans1 = append(ans1, cur)
				}
				cur = x
				if cur == f {
					break
				}
			}
			// second strategy from next
			ans2 := make([]int, 0)
			cur = direct[f]
			for {
				x := direct[cur]
				if s1[cur] {
					s1[cur] = !s1[cur]
					s1[x] = !s1[x]
					ans2 = append(ans2, cur)
				}
				cur = x
				if cur == direct[f] {
					break
				}
			}
			if len(ans1) < len(ans2) {
				ans = append(ans, ans1...)
			} else {
				ans = append(ans, ans2...)
			}
		}
		// output
		fmt.Fprintln(writer, len(ans))
		for i, v := range ans {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprintf(writer, "%d", v+1)
		}
		writer.WriteByte('\n')
	}
}
