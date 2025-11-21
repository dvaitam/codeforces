package main

import (
	"bufio"
	"fmt"
	"os"
)

type queue struct {
	data []int
	head int
}

func (q *queue) push(x int) {
	q.data = append(q.data, x)
}

func (q *queue) empty() bool {
	return q.head >= len(q.data)
}

func (q *queue) pop() int {
	x := q.data[q.head]
	q.head++
	if q.head*2 >= len(q.data) {
		q.data = q.data[q.head:]
		q.head = 0
	}
	return x
}

func assignLabels(n, k int, adj [][]int) []int {
	labels := make([]int, n+1)
	visited := make([]bool, n+1)
	q := &queue{}
	q.push(1)
	visited[1] = true
	labels[1] = 0
	for !q.empty() {
		u := q.pop()
		for _, v := range adj[u] {
			expected := (labels[u] + 1) % k
			if !visited[v] {
				visited[v] = true
				labels[v] = expected
				q.push(v)
			}
		}
	}
	return labels
}

func prefixFunction(pattern []int64) []int {
	n := len(pattern)
	pi := make([]int, n)
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && pattern[i] != pattern[j] {
			j = pi[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
		}
		pi[i] = j
	}
	return pi
}

func findShifts(pattern []int64, text []int64, k int) []int {
	shifts := []int{}
	if len(pattern) == 0 {
		return shifts
	}
	pi := prefixFunction(pattern)
	j := 0
	for i := 0; i < len(text); i++ {
		for j > 0 && text[i] != pattern[j] {
			j = pi[j-1]
		}
		if text[i] == pattern[j] {
			j++
		}
		if j == len(pattern) {
			start := i - len(pattern) + 1
			if start < k {
				shifts = append(shifts, start)
			}
			j = pi[j-1]
		}
	}
	return shifts
}

func rotateArrayShift(arr []int64) []int64 {
	res := make([]int64, len(arr)*2)
	copy(res, arr)
	copy(res[len(arr):], arr)
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		var m1 int
		fmt.Fscan(in, &m1)
		adj1 := make([][]int, n+1)
		for i := 0; i < m1; i++ {
			var v, u int
			fmt.Fscan(in, &v, &u)
			adj1[v] = append(adj1[v], u)
		}

		labels1 := assignLabels(n, k, adj1)

		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		var m2 int
		fmt.Fscan(in, &m2)
		adj2 := make([][]int, n+1)
		for i := 0; i < m2; i++ {
			var v, u int
			fmt.Fscan(in, &v, &u)
			adj2[v] = append(adj2[v], u)
		}

		labels2 := assignLabels(n, k, adj2)

		outCount1 := make([]int64, k)
		inCount1 := make([]int64, k)
		outCount2 := make([]int64, k)
		inCount2 := make([]int64, k)

		for i := 0; i < n; i++ {
			label := labels1[i+1]
			if a[i] == 1 {
				outCount1[label]++
			} else {
				inCount1[label]++
			}
		}

		for i := 0; i < n; i++ {
			label := labels2[i+1]
			if b[i] == 1 {
				outCount2[label]++
			} else {
				inCount2[label]++
			}
		}

		sumOut1 := int64(0)
		sumIn2 := int64(0)
		sumOut2 := int64(0)
		sumIn1 := int64(0)
		for i := 0; i < k; i++ {
			sumOut1 += outCount1[i]
			sumIn2 += inCount2[i]
			sumOut2 += outCount2[i]
			sumIn1 += inCount1[i]
		}
		if sumOut1 != sumIn2 || sumOut2 != sumIn1 {
			fmt.Fprintln(out, "NO")
			continue
		}

		textBin := rotateArrayShift(inCount2)
		shifts1 := findShifts(outCount1, textBin, k)
		if len(shifts1) == 0 {
			fmt.Fprintln(out, "NO")
			continue
		}

		textBout := rotateArrayShift(outCount2)
		shifts2 := findShifts(inCount1, textBout, k)
		if len(shifts2) == 0 {
			fmt.Fprintln(out, "NO")
			continue
		}

		possibleShift := make([]bool, k)
		for _, s := range shifts2 {
			possibleShift[s%k] = true
		}

		ok := false
		for _, s := range shifts1 {
			t := (s - 2) % k
			if t < 0 {
				t += k
			}
			if possibleShift[t] {
				ok = true
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
