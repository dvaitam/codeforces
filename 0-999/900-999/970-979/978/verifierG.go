package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

type examInfo struct {
	s int
	d int
	c int
}

type testG struct {
	n     int
	exams []examInfo
}

func genTestsG() []testG {
	rand.Seed(48)
	tests := make([]testG, 100)
	for i := range tests {
		n := rand.Intn(10) + 2
		m := rand.Intn(n) + 1
		exams := make([]examInfo, m)
		usedDays := make(map[int]bool)
		for j := 0; j < m; j++ {
			// choose exam day unique
			var d int
			for {
				d = rand.Intn(n-1) + 2
				if !usedDays[d] {
					usedDays[d] = true
					break
				}
			}
			s := rand.Intn(d-1) + 1
			c := rand.Intn(d-s) + 1
			exams[j] = examInfo{s: s, d: d, c: c}
		}
		tests[i] = testG{n: n, exams: exams}
	}
	return tests
}

type pqExam []exam

type exam struct {
	d   int
	idx int
}

func (p pqExam) Len() int            { return len(p) }
func (p pqExam) Less(i, j int) bool  { return p[i].d < p[j].d }
func (p pqExam) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pqExam) Push(x interface{}) { *p = append(*p, x.(exam)) }
func (p *pqExam) Pop() interface{} {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[:n-1]
	return x
}

func solveG(tc testG) ([]int, bool) {
	n := tc.n
	m := len(tc.exams)
	examDay := make([]int, n+1)
	start := make([][]int, n+1)
	left := make([]int, m+1)
	for i, e := range tc.exams {
		idx := i + 1
		examDay[e.d] = idx
		start[e.s] = append(start[e.s], idx)
		left[idx] = e.c
	}
	res := make([]int, n+1)
	pq := &pqExam{}
	heap.Init(pq)
	for day := 1; day <= n; day++ {
		for _, idx := range start[day] {
			heap.Push(pq, exam{d: tc.exams[idx-1].d, idx: idx})
		}
		if examDay[day] != 0 {
			idx := examDay[day]
			if left[idx] != 0 {
				return nil, false
			}
			res[day] = m + 1
			continue
		}
		if pq.Len() > 0 {
			e := heap.Pop(pq).(exam)
			res[day] = e.idx
			left[e.idx]--
			if left[e.idx] > 0 {
				heap.Push(pq, e)
			}
		} else {
			res[day] = 0
		}
	}
	for i := 1; i <= m; i++ {
		if left[i] != 0 {
			return nil, false
		}
	}
	return res[1:], true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsG()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, len(tc.exams))
		for _, e := range tc.exams {
			fmt.Fprintf(&input, "%d %d %d\n", e.s, e.d, e.c)
		}
	}

	expected := make([][]int, len(tests))
	possible := make([]bool, len(tests))
	for i, tc := range tests {
		res, ok := solveG(tc)
		possible[i] = ok
		if ok {
			expected[i] = res
		}
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i, ok := range possible {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		if !ok {
			val := scanner.Text()
			if val != "-1" {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected -1\n", i+1)
				os.Exit(1)
			}
			continue
		}
		// expecting n integers
		outArr := make([]int, len(expected[i]))
		first, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
			os.Exit(1)
		}
		outArr[0] = first
		for j := 1; j < len(outArr); j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
				os.Exit(1)
			}
			val, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
				os.Exit(1)
			}
			outArr[j] = val
		}
		for j, v := range expected[i] {
			if outArr[j] != v {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
