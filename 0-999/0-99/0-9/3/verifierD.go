package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type item struct {
	diff int
	idx  int
}

type priorityQueue []item

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].diff < pq[j].diff }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(item)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

func solve(pattern string, costs [][2]int) (int, string, bool) {
	res := []rune(pattern)
	pq := &priorityQueue{}
	heap.Init(pq)
	qidx := 0
	cost := 0
	balance := 0
	for i, ch := range res {
		switch ch {
		case '(':
			balance++
		case ')':
			balance--
		case '?':
			a := costs[qidx][0]
			b := costs[qidx][1]
			qidx++
			cost += a
			res[i] = '('
			heap.Push(pq, item{diff: b - a, idx: i})
			balance++
		}
		if balance < 0 {
			if pq.Len() == 0 {
				return 0, "", false
			}
			it := heap.Pop(pq).(item)
			res[it.idx] = ')'
			cost += it.diff
			balance += 2
		}
	}
	for balance > 0 {
		if pq.Len() == 0 {
			return 0, "", false
		}
		it := heap.Pop(pq).(item)
		res[it.idx] = ')'
		cost += it.diff
		balance -= 2
	}
	if balance != 0 {
		return 0, "", false
	}
	return cost, string(res), true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		pattern := fields[0]
		costs := make([][2]int, 0)
		for _, p := range fields[1:] {
			var a, b int
			fmt.Sscanf(p, "%d,%d", &a, &b)
			costs = append(costs, [2]int{a, b})
		}
		if strings.Count(pattern, "?") != len(costs) {
			fmt.Printf("test %d invalid: mismatch costs\n", idx)
			os.Exit(1)
		}
		expCost, expStr, ok := solve(pattern, costs)
		if !ok {
			expStr = "-1"
		}
		// build input for binary
		var buf bytes.Buffer
		buf.WriteString(pattern)
		buf.WriteByte('\n')
		for i, c := range costs {
			fmt.Fprintf(&buf, "%d %d", c[0], c[1])
			if i+1 < len(costs) {
				buf.WriteByte('\n')
			}
		}
		buf.WriteByte('\n')
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if !ok {
			if len(outLines) == 0 || strings.TrimSpace(outLines[0]) != "-1" {
				fmt.Printf("Test %d failed: expected -1 got %s\n", idx, string(out))
				os.Exit(1)
			}
			continue
		}
		if len(outLines) < 2 {
			fmt.Printf("Test %d failed: output should have 2 lines\n", idx)
			os.Exit(1)
		}
		var gotCost int
		fmt.Sscan(outLines[0], &gotCost)
		gotStr := strings.TrimSpace(outLines[1])
		if gotCost != expCost || gotStr != expStr {
			fmt.Printf("Test %d failed\nexpected:\n%d\n%s\n\ngot:\n%s", idx, expCost, expStr, string(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
