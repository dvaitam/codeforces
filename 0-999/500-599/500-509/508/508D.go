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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		fmt.Fprintln(writer, "NO")
		return
	}
	const maxv = 200 * 200
	graph := make([][]int, maxv)
	num := make([]int, maxv)

	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		u1 := int(s[0])*200 + int(s[1])
		u2 := int(s[1])*200 + int(s[2])
		graph[u1] = append(graph[u1], u2)
		num[u1]--
		num[u2]++
	}

	root := -1
	countNegLE := 0
	for i := 0; i < maxv; i++ {
		if num[i] <= -1 {
			root = i
			countNegLE++
		}
	}
	if countNegLE >= 2 || (root != -1 && num[root] < -1) {
		fmt.Fprintln(writer, "NO")
		return
	}
	if root == -1 {
		for i := 0; i < maxv; i++ {
			if len(graph[i]) > 0 {
				root = i
				break
			}
		}
	}
	// Hierholzer's algorithm
	stack := []int{root}
	ansChars := make([]byte, 0, n+1)
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		if len(graph[v]) > 0 {
			nxt := graph[v][len(graph[v])-1]
			graph[v] = graph[v][:len(graph[v])-1]
			stack = append(stack, nxt)
		} else {
			stack = stack[:len(stack)-1]
			ansChars = append(ansChars, byte(v%200))
		}
	}
	if len(ansChars) != n+1 {
		fmt.Fprintln(writer, "NO")
		return
	}
	// reverse
	for i, j := 0, len(ansChars)-1; i < j; i, j = i+1, j-1 {
		ansChars[i], ansChars[j] = ansChars[j], ansChars[i]
	}
	// build full answer
	fullAns := make([]byte, 0, n+2)
	fullAns = append(fullAns, byte(root/200))
	fullAns = append(fullAns, ansChars...)

	fmt.Fprintln(writer, "YES")
	writer.Write(fullAns)
	writer.WriteByte('\n')
}
