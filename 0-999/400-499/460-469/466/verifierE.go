package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solve(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	fmt.Fscan(reader, &n, &m)
	parent := make([]int, n+1)
	addTime := make([]int, n+1)
	s := make([]int, 1, m+1)
	packetTime := make([]int, 1, m+1)
	type Query struct{ x, pid, idx int }
	var queries []Query
	ansCount := 0
	for t := 1; t <= m; t++ {
		var tp int
		fmt.Fscan(reader, &tp)
		if tp == 1 {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			parent[x] = y
			addTime[x] = t
		} else if tp == 2 {
			var x int
			fmt.Fscan(reader, &x)
			s = append(s, x)
			packetTime = append(packetTime, t)
		} else if tp == 3 {
			var x, pid int
			fmt.Fscan(reader, &x, &pid)
			queries = append(queries, Query{x: x, pid: pid, idx: ansCount})
			ansCount++
		}
	}
	children := make([][]int, n+1)
	for v := 1; v <= n; v++ {
		p := parent[v]
		if p != 0 {
			children[p] = append(children[p], v)
		}
	}
	depth := make([]int, n+1)
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	timer := 0
	type stackEntry struct{ v, idx, state int }
	var stack []stackEntry
	for v := 1; v <= n; v++ {
		if parent[v] == 0 {
			stack = append(stack, stackEntry{v: v})
			for len(stack) > 0 {
				e := &stack[len(stack)-1]
				if e.state == 0 {
					timer++
					tin[e.v] = timer
					e.state = 1
				}
				if e.idx < len(children[e.v]) {
					c := children[e.v][e.idx]
					e.idx++
					depth[c] = depth[e.v] + 1
					stack = append(stack, stackEntry{v: c})
				} else {
					tout[e.v] = timer
					stack = stack[:len(stack)-1]
				}
			}
		}
	}
	const LOG = 18
	up := make([][]int, LOG)
	timeUp := make([][]int, LOG)
	for j := 0; j < LOG; j++ {
		up[j] = make([]int, n+1)
		timeUp[j] = make([]int, n+1)
	}
	for v := 1; v <= n; v++ {
		up[0][v] = parent[v]
		timeUp[0][v] = addTime[v]
	}
	for j := 1; j < LOG; j++ {
		for v := 1; v <= n; v++ {
			u := up[j-1][v]
			up[j][v] = up[j-1][u]
			t1 := timeUp[j-1][v]
			t2 := timeUp[j-1][u]
			if t2 > t1 {
				timeUp[j][v] = t2
			} else {
				timeUp[j][v] = t1
			}
		}
	}
	answers := make([]string, ansCount)
	for qi, q := range queries {
		x := q.x
		pid := q.pid
		u := s[pid]
		t0 := packetTime[pid]
		if !(tin[x] <= tin[u] && tout[u] <= tout[x]) {
			answers[qi] = "NO"
			continue
		}
		diff := depth[u] - depth[x]
		cur := u
		maxT := 0
		for j := 0; j < LOG; j++ {
			if diff&(1<<j) != 0 {
				if timeUp[j][cur] > maxT {
					maxT = timeUp[j][cur]
				}
				cur = up[j][cur]
			}
		}
		if maxT <= t0 {
			answers[qi] = "YES"
		} else {
			answers[qi] = "NO"
		}
	}
	var out bytes.Buffer
	for i := 0; i < ansCount; i++ {
		out.WriteString(answers[i])
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if len(lines) == 0 {
				continue
			}
			idx++
			input := strings.Join(lines, "\n") + "\n"
			expect := solve(input)
			cmd := exec.Command(bin)
			cmd.Stdin = strings.NewReader(input)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Test %d: runtime error: %v\n", idx, err)
				os.Exit(1)
			}
			got := strings.TrimSpace(string(out))
			if got != expect {
				fmt.Printf("Test %d failed: expected %s got %s\n", idx, expect, got)
				os.Exit(1)
			}
			lines = lines[:0]
			continue
		}
		lines = append(lines, strings.TrimSpace(line))
	}
	if len(lines) > 0 {
		idx++
		input := strings.Join(lines, "\n") + "\n"
		expect := solve(input)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expect {
			fmt.Printf("Test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
