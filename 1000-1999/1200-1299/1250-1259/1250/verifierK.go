package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type event struct {
	s, e int
	kind int // 0 lecture, 1 seminar
	idx  int
}

type item struct {
	end int
	idx int
}

type pq []item

func (p pq) Len() int            { return len(p) }
func (p pq) Less(i, j int) bool  { return p[i].end < p[j].end }
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x interface{}) { *p = append(*p, x.(item)) }
func (p *pq) Pop() interface{} {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[:n-1]
	return x
}

// Embedded copy of testcasesK.txt so the verifier is self-contained.
const testcasesRaw = `100
2 2 5 3
3 8
5 5
4 7
4 8
2 2 3 2
1 2
2 5
4 5
4 6
3 3 5 5
1 4
3 4
1 4
3 4
1 6
2 7
2 2 3 3
2 3
5 6
4 4
4 8
3 3 5 2
3 3
2 6
3 5
3 7
3 7
4 9
2 1 3 1
2 4
3 4
1 3
3 3 1 3
3 7
1 1
3 4
4 6
3 3
5 5
3 1 5 1
1 5
2 6
1 2
5 7
3 1 1 2
5 6
3 4
3 7
4 9
3 2 2 5
1 3
2 4
2 5
4 6
3 3
3 1 3 2
3 5
2 3
3 4
5 5
1 2 3 5
1 2
3 4
1 3
1 1 1 3
3 8
3 6
2 1 5 5
1 5
4 7
3 5
3 2 1 4
4 4
4 4
3 6
3 4
4 8
1 3 4 2
2 3
5 5
2 3
3 5
2 2 4 5
5 5
3 5
3 6
5 8
3 1 3 4
1 2
1 2
3 8
5 10
1 1 4 1
5 9
2 2
2 3 1 5
4 9
4 4
4 4
1 3
2 4
2 1 2 3
3 4
2 3
4 5
2 1 2 3
3 8
1 2
3 8
1 3 4 5
5 6
4 9
2 5
5 5
3 2 1 4
3 7
5 5
4 6
2 2
2 2
1 2 1 3
5 7
1 6
2 4
2 3 1 5
2 4
4 9
2 7
2 2
3 5
1 1 3 4
1 5
3 7
1 3 1 2
2 7
5 8
1 2
4 7
2 3 4 3
5 10
4 7
1 3
1 4
2 7
3 1 1 3
1 4
4 9
4 5
3 6
1 2 5 2
3 5
5 10
1 2
1 1 5 3
3 6
5 7
2 3 2 2
3 6
3 4
1 3
2 4
2 7
3 3 4 1
4 7
4 7
4 4
1 2
3 3
2 6
3 3 2 2
4 5
2 5
5 9
1 2
1 1
5 10
1 2 3 3
5 8
2 2
1 1
2 1 5 5
1 6
3 8
5 5
3 2 2 5
5 5
4 7
4 8
5 10
5 8
3 3 4 2
5 5
5 9
1 5
1 1
4 8
4 7
1 1 3 1
4 9
5 6
2 3 1 5
1 3
5 8
2 7
3 4
1 6
3 1 4 5
2 7
3 8
1 6
2 5
2 2 1 4
3 3
3 4
2 5
5 10
1 2 2 5
1 3
4 7
1 2
1 2 2 4
5 5
5 5
4 8
2 2 1 2
5 6
4 5
4 5
4 6
3 1 4 2
2 5
1 2
4 8
1 2
2 2 3 3
5 5
1 6
4 7
4 5
1 2 1 1
3 8
2 4
1 2
3 1 5 4
1 3
1 1
5 8
5 5
2 3 5 5
1 6
5 5
1 2
1 6
4 9
3 3 5 1
4 6
5 10
4 9
1 4
5 7
3 6
3 1 3 1
3 5
2 3
3 4
5 5
3 1 2 4
4 4
1 5
2 3
1 2
2 3 5 2
4 7
2 3
1 3
2 6
1 4
2 2 5 5
3 7
5 10
4 7
3 7
2 2 1 5
1 4
1 2
3 5
1 4
3 2 4 2
3 7
3 7
1 5
3 7
2 5
1 2 1 2
3 7
3 3
5 8
3 1 5 3
1 3
1 2
2 5
3 5
1 1 2 3
3 6
1 3
3 1 5 5
3 8
5 7
2 5
1 6
3 3 5 5
4 5
4 7
4 8
3 8
4 5
4 6
1 2 4 2
5 7
3 4
1 2
1 3 5 4
1 6
4 9
3 8
3 5
3 2 1 1
1 3
3 6
5 6
4 5
5 10
1 3 1 3
5 5
4 6
3 6
2 4
2 3 5 1
5 8
3 5
3 3
5 6
4 5
1 1 4 2
1 6
5 6
2 3 1 3
2 3
1 5
3 6
5 5
5 9
2 1 1 3
1 1
4 7
4 8
2 2 5 5
5 5
2 7
3 4
5 5
3 3 5 4
1 5
1 1
1 2
3 6
4 8
1 4
1 3 4 2
3 6
2 2
4 4
4 7
3 1 1 4
5 5
1 5
2 7
4 5
3 3 1 1
1 4
3 5
1 3
4 7
5 5
1 2
1 1 4 1
4 7
5 6
3 3 1 2
4 4
1 5
5 6
2 6
1 1
1 1
3 3 3 2
3 3
3 4
3 4
5 5
5 9
4 4
3 1 5 3
4 5
5 6
4 9
3 4
1 3 2 3
3 5
2 4
5 8
3 8
2 3 1 4
3 4
3 4
1 4
3 8
4 4
1 2 2 5
5 6
4 5
4 4
1 2 1 3
2 6
2 6
4 7
2 2 2 5
1 6
1 4
3 5
4 7
2 2 4 4
4 5
1 1
3 5
5 5
2 2 3 4
1 5
5 6
1 5
1 4
1 1 3 2
3 3
4 7
2 1 1 2
3 7
1 1
5 9
2 1 2 2
2 4
1 1
2 5
1 1 4 2
3 6
4 9
2 2 2 5
5 5
2 5
2 7
1 3
3 1 5 3
4 9
4 9
2 5
2 2
2 1 1 3
3 5
2 3
5 5
1 3 1 3
1 2
2 5
5 8
3 8
1 2 5 1
4 9
4 9
5 5
1 1 4 4
3 4
4 6
1 1 1 4
1 3
2 7
1 2 2 2
1 4
2 3
2 5
2 3 5 1
1 6
4 4
4 6
3 5
5 9`

// solveAll implements the logic from 1250K.go for all testcases in the input.
func solveAll(input string) (string, error) {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return "", err
	}
	var out strings.Builder
	for ; t > 0; t-- {
		var n, m, x, y int
		if _, err := fmt.Fscan(in, &n, &m, &x, &y); err != nil {
			return "", err
		}
		lectures := make([][2]int, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(in, &lectures[i][0], &lectures[i][1]); err != nil {
				return "", err
			}
		}
		seminars := make([][2]int, m)
		for i := 0; i < m; i++ {
			if _, err := fmt.Fscan(in, &seminars[i][0], &seminars[i][1]); err != nil {
				return "", err
			}
		}
		events := make([]event, 0, n+m)
		for i := 0; i < n; i++ {
			events = append(events, event{s: lectures[i][0], e: lectures[i][1], kind: 0, idx: i})
		}
		for i := 0; i < m; i++ {
			events = append(events, event{s: seminars[i][0], e: seminars[i][1], kind: 1, idx: i})
		}
		sort.Slice(events, func(i, j int) bool {
			if events[i].s == events[j].s {
				return events[i].kind < events[j].kind
			}
			return events[i].s < events[j].s
		})

		availableHD := make([]int, x)
		for i := 0; i < x; i++ {
			availableHD[i] = i + 1
		}
		availableOrd := make([]int, y)
		for i := 0; i < y; i++ {
			availableOrd[i] = x + i + 1
		}
		var hdpq, ordpq pq
		assignmentsL := make([]int, n)
		assignmentsS := make([]int, m)
		possible := true

		for _, ev := range events {
			for len(hdpq) > 0 && hdpq[0].end <= ev.s {
				it := heap.Pop(&hdpq).(item)
				availableHD = append(availableHD, it.idx)
			}
			for len(ordpq) > 0 && ordpq[0].end <= ev.s {
				it := heap.Pop(&ordpq).(item)
				availableOrd = append(availableOrd, it.idx)
			}

			if ev.kind == 0 {
				if len(availableHD) == 0 {
					possible = false
					break
				}
				idx := availableHD[len(availableHD)-1]
				availableHD = availableHD[:len(availableHD)-1]
				assignmentsL[ev.idx] = idx
				heap.Push(&hdpq, item{end: ev.e, idx: idx})
			} else {
				if len(availableOrd) > 0 {
					idx := availableOrd[len(availableOrd)-1]
					availableOrd = availableOrd[:len(availableOrd)-1]
					assignmentsS[ev.idx] = idx
					heap.Push(&ordpq, item{end: ev.e, idx: idx})
				} else if len(availableHD) > 0 {
					idx := availableHD[len(availableHD)-1]
					availableHD = availableHD[:len(availableHD)-1]
					assignmentsS[ev.idx] = idx
					heap.Push(&hdpq, item{end: ev.e, idx: idx})
				} else {
					possible = false
					break
				}
			}
		}

		if possible {
			fmt.Fprintln(&out, "YES")
			for i := 0; i < n; i++ {
				if i > 0 {
					fmt.Fprint(&out, " ")
				}
				fmt.Fprint(&out, assignmentsL[i])
			}
			for j := 0; j < m; j++ {
				fmt.Fprint(&out, " ", assignmentsS[j])
			}
			fmt.Fprintln(&out)
		} else {
			fmt.Fprintln(&out, "NO")
		}
	}
	return strings.TrimSpace(out.String()), nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}

	input := strings.TrimSpace(testcasesRaw) + "\n"
	expect, err := solveAll(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "solver failed: %v\n", err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	got, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expect) {
		fmt.Println("output mismatch")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
