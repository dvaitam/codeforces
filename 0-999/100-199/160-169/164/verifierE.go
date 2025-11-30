package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `
2 6 11 4 9 9 1
1 6 9 2
4 3 6 4 5 9 4 10 15 3 11 16 5
5 3 4 2 7 8 1 12 12 1 14 16 1 17 17 1
5 6 10 2 11 14 3 12 14 1 13 17 2 17 22 5
3 5 5 1 10 12 1 13 18 1
3 5 7 2 6 8 3 7 9 2
2 4 5 1 6 7 2
1 6 8 3
2 4 7 1 9 14 1
5 5 10 1 8 13 2 13 15 3 14 14 1 19 21 1
2 2 4 3 5 6 2
3 3 7 5 5 5 1 7 9 2
3 4 4 1 8 13 3 13 14 1
4 5 8 3 10 12 2 12 14 1 16 21 1
4 6 7 1 7 10 2 10 15 1 14 19 5
3 6 6 1 10 15 4 15 17 1
2 3 8 6 4 4 1
3 4 7 2 8 9 1 11 12 1
3 3 3 1 8 9 1 13 13 1
3 2 7 2 8 8 3 10 10 3
2 4 2 6 7 11 1
1 4 11 2
3 4 6 1 10 12 4 14 15 1
1 5 7 2
5 3 6 2 7 10 3 12 15 3 15 16 5 18 19 5
2 5 4 1 5 10 1
2 6 3 3 12 14 4
2 2 9 1 3 6 2
4 3 7 3 13 14 1 15 16 1 19 21 5
3 6 4 1 10 16 3 19 19 1
4 6 8 1 10 14 4 16 18 2 18 20 5
4 3 7 2 9 14 3 15 18 5 18 19 5
2 6 7 1 8 9 2
5 6 7 1 8 9 5 11 13 5 13 18 4 14 20 2
2 4 10 1 15 17 4
1 2 5 3
1 3 7 2
1 2 4 4
1 5 6 3
4 4 9 2 10 12 1 16 18 4 20 21 1
1 4 4 1
2 4 8 1 9 9 4
4 4 6 2 8 10 5 11 15 4 16 16 3
2 5 9 2 15 18 4
3 2 5 4 9 12 2 15 18 3
1 3 4 1
3 6 8 1 10 10 2 12 17 4
2 6 10 1 13 13 3
5 5 8 1 9 9 1 10 15 3 16 17 1 16 19 5
2 3 3 4 6 10 3
2 1 3 1 8 12 3
2 1 1 2 2 4 1
4 1 5 2 7 12 3 15 17 1 21 23 5
2 5 8 1 14 15 1
2 1 1 1 6 11 1
5 2 5 1 9 12 4 14 15 5 15 16 1 20 25 3
5 4 7 4 11 15 5 14 14 3 18 18 4 19 19 2
4 2 5 3 10 10 3 14 14 3 21 25 5
2 1 2 3 9 14 1
5 2 2 1 3 7 3 8 10 2 13 17 3 17 22 3
2 1 1 2 4 7 2
4 4 4 3 7 8 1 12 12 4 17 17 1
4 2 6 3 7 8 2 11 14 4 14 16 3
1 3 5 2
4 5 5 3 10 11 2 11 15 3 17 19 2
2 6 6 3 8 10 5
1 5 5 2
2 1 3 3 7 9 4
2 1 5 3 9 13 3
5 4 5 4 6 8 4 8 12 2 9 9 1 10 13 3
5 2 4 3 5 10 1 10 10 4 11 16 5 13 16 3
1 6 6 1
4 3 6 2 6 10 1 10 13 2 13 17 3
3 4 5 3 8 10 2 13 16 4
2 2 6 4 8 9 5
1 6 8 1
4 6 7 3 7 8 1 10 11 1 11 16 4
5 1 4 4 6 9 2 10 11 1 12 13 3 14 15 2
5 2 5 5 7 8 4 9 10 5 13 13 5 14 17 2
4 5 9 2 10 11 1 16 18 5 21 21 3
1 2 3 3
4 3 3 1 5 5 4 5 9 5 9 14 1
5 4 6 4 8 11 5 12 17 2 17 17 5 19 19 2
5 2 7 2 11 15 4 12 13 2 16 20 5 21 24 3
1 4 4 1
1 6 7 2
3 4 7 2 8 11 3 13 18 3
3 3 5 2 6 10 4 9 10 4
2 6 9 4 14 17 3
2 1 1 2 5 10 3
3 2 4 1 8 12 2 16 16 3
2 4 5 1 7 10 1
5 1 4 1 5 6 3 8 10 3 11 15 3 17 22 5
5 1 2 3 4 6 1 7 7 2 13 17 2 17 22 5
4 3 5 1 7 7 5 13 13 5 19 22 5
2 3 6 2 6 10 1
`

type Item struct {
	t   int64
	idx int
}

type MaxHeap []Item

func (h MaxHeap) Len() int { return len(h) }
func (h MaxHeap) Less(i, j int) bool {
	if h[i].t != h[j].t {
		return h[i].t > h[j].t
	}
	return h[i].idx > h[j].idx
}
func (h MaxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(Item))
}
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

type testCase struct {
	n int
	l []int64
	r []int64
	t []int64
}

func solveCase(tc testCase) []int {
	ans := int64(0)
	h := &MaxHeap{}
	out := make([]int, tc.n)
	for i := 0; i < tc.n; i++ {
		start := ans + 1
		if start < tc.l[i] {
			start = tc.l[i]
		}
		end := start + tc.t[i] - 1
		if end <= tc.r[i] {
			ans = end
			heap.Push(h, Item{t: tc.t[i], idx: i + 1})
			out[i] = 0
		} else if h.Len() > 0 && (*h)[0].t > tc.t[i] {
			top := heap.Pop(h).(Item)
			ans -= top.t
			start = ans + 1
			if start < tc.l[i] {
				start = tc.l[i]
			}
			end = start + tc.t[i] - 1
			if end <= tc.r[i] {
				ans = end
				heap.Push(h, Item{t: tc.t[i], idx: i + 1})
				out[i] = top.idx
			} else {
				heap.Push(h, top)
				out[i] = -1
			}
		} else {
			out[i] = -1
		}
	}
	return out
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		if len(fields) != 1+3*n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, 1+3*n, len(fields))
		}
		tc := testCase{n: n, l: make([]int64, n), r: make([]int64, n), t: make([]int64, n)}
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[1+3*i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse l[%d]: %v", idx+1, i, err)
			}
			tc.l[i] = v
			v, err = strconv.ParseInt(fields[1+3*i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse r[%d]: %v", idx+1, i, err)
			}
			tc.r[i] = v
			v, err = strconv.ParseInt(fields[1+3*i+2], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse t[%d]: %v", idx+1, i, err)
			}
			tc.t[i] = v
		}
		cases = append(cases, tc)
	}
	return cases, nil
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		expected := solveCase(tc)

		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for j := 0; j < tc.n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.l[j], tc.r[j], tc.t[j]))
		}

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		var expSB strings.Builder
		for idx, v := range expected {
			if idx > 0 {
				expSB.WriteByte('\n')
			}
			expSB.WriteString(strconv.Itoa(v))
		}
		expectedStr := strings.TrimSpace(expSB.String())
		if strings.TrimSpace(got) != expectedStr {
			fmt.Printf("case %d failed\nexpected:\n%s\ngot:\n%s\n", i+1, expectedStr, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
