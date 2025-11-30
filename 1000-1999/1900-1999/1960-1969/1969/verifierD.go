package main

import (
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type item struct{ a, b int64 }

type maxHeap []int64

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

const testcasesD = `4 4 18 5 12 20 16 19 3 20
1 1 9 18
4 1 16 18 18 16 13 5 8 5
7 0 3 6 19 2 10 1 9 16 20 13 14 13 19 15
3 2 4 2 5 16 7 9
7 4 14 17 13 19 12 18 19 14 19 8 11 1 9 20
3 2 18 19 19 4 7 19
5 2 4 3 16 16 3 12 3 14 5 1
5 3 14 4 2 20 20 2 13 19 11 18
5 4 8 2 10 1 3 4 20 18 2 7
7 4 20 9 5 2 11 11 12 5 13 13 15 17 13 20
2 2 17 9 14 8
5 3 9 17 10 18 11 1 14 19 11 1
7 2 2 11 15 12 12 20 9 16 1 19 2 1 12 9
8 4 19 20 11 6 12 6 11 12 20 9 10 13 4 1 19 5
5 4 8 9 8 11 6 14 4 4 20 11
6 5 8 15 6 3 11 7 19 15 9 8 4 2
4 2 19 6 9 11 3 20 12 19
3 3 10 17 9 15 12 14
5 3 19 14 2 14 5 7 1 16 20 17
7 3 2 15 17 10 18 11 8 3 19 10 4 8 2 2
4 3 19 2 1 16 4 6 17 10
4 0 17 18 14 2 20 4 11 5
5 4 16 2 12 8 7 4 18 4 6 8
5 1 1 16 19 13 2 9 8 9 20 17
7 0 16 11 1 2 5 2 4 2 3 16 2 3 17 17
8 5 6 11 3 12 13 13 19 10 12 9 7 11 14 4 5 18
1 1 3 19
3 0 12 15 20 18 13 2
7 0 12 16 11 14 14 15 1 8 7 18 9 19 3 14
4 3 5 1 11 12 18 9 4 15
2 2 17 13 4 11
2 2 1 16 5 8
7 0 17 3 19 4 13 6 1 11 4 1 4 16 10 19
5 0 2 19 17 17 8 4 18 4 18 2
6 6 19 6 3 8 6 8 15 20 13 9 12 20
7 5 18 14 3 13 17 8 14 6 14 19 19 17 16 5
7 2 6 4 16 16 17 15 19 6 5 9 7 5 19 17
6 1 18 10 14 20 19 19 9 7 10 1 9 16
7 3 6 19 12 8 11 16 5 14 16 20 7 15 19 18
1 1 3 13
1 1 8 8
2 0 9 8 7 9
3 1 20 2 9 6 2 11
3 3 3 3 4 3 9 10
1 1 15 19
6 0 1 11 11 14 13 16 3 7 19 16 13 5
6 0 9 3 14 4 15 17 9 4 17 12 12 15
5 5 9 4 11 19 18 17 4 16 17 12
1 1 19 6
3 1 12 15 4 4 18 5
6 5 20 14 18 10 6 15 16 10 6 3 4 6
7 5 4 9 9 13 2 5 2 16 17 9 8 17 12 11
7 7 18 3 12 16 4 5 9 19 4 4 19 4 6 7
7 6 5 19 20 5 13 7 18 17 6 19 6 7 9 12
5 0 15 14 13 11 18 19 10 16 17 10
8 0 20 7 1 4 8 16 6 17 15 7 7 17 7 2 17 15
2 2 10 5 5 15
2 2 2 1 12 20
4 4 3 16 18 1 11 11 11 12
3 0 20 2 3 11 7 3
4 3 8 16 11 4 2 14 3 7
3 3 16 16 3 18 14 7
8 4 1 15 15 13 15 6 15 2 9 12 12 15 17 12 20 13
4 0 7 9 12 5 15 18 7 6
4 0 6 19 13 17 6 1 5 4
3 3 16 6 2 1 13 15
6 3 2 2 8 13 2 13 16 1 8 8 4 13
8 3 6 11 20 4 12 4 20 2 10 9 15 10 16 8 18 9
1 1 12 11
2 0 14 3 19 20
1 0 1 3
1 0 17 2
8 0 7 17 11 7 16 11 16 12 2 13 10 20 13 3 10 6
7 1 17 13 18 11 18 13 6 13 18 12 6 12 14 15
4 3 16 12 9 6 17 20 13 16
1 0 6 1
8 1 4 11 8 20 2 20 2 15 15 11 12 1 3 7 13 4
6 4 10 4 15 3 7 8 2 5 5 19 1 4
4 2 7 8 18 17 14 17 20 11
4 3 6 20 3 2 4 20 1 4
4 2 3 4 15 13 8 20 4 16
6 3 20 15 4 10 20 15 13 7 4 18 1 15
5 5 3 11 12 7 16 3 18 12 14 3
4 1 12 2 11 8 14 15 3 9
4 2 6 7 7 20 15 18 14 12
4 4 14 16 14 16 19 2 10 1
3 0 1 5 10 17 17 2
8 0 7 7 9 16 14 2 12 15 7 10 5 4 15 10 14 15
2 0 5 16 10 13
6 1 14 10 15 16 17 18 8 12 10 10 1 15
6 2 10 8 17 1 1 5 17 5 18 1 6 2
1 0 15 12
6 4 2 16 6 8 1 9 14 11 2 20 18 4
8 4 9 8 16 14 9 11 2 1 14 2 6 19 8 5 14 17
6 4 5 9 1 6 2 1 16 2 15 15 17 20
7 5 17 6 10 6 3 5 18 4 14 12 15 15 9 9
8 4 17 5 19 11 5 17 2 14 16 8 15 19 20 9 1 11
2 1 5 9 9 4
7 1 12 2 17 16 15 7 10 12 6 13 13 11 2 9`

func expected(n, k int, a, b []int64) int64 {
	items := make([]item, n)
	for i := 0; i < n; i++ {
		items[i] = item{a[i], b[i]}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].b == items[j].b {
			return items[i].a < items[j].a
		}
		return items[i].b > items[j].b
	})

	suff := make([]int64, n+1)
	for i := n - 1; i >= 0; i-- {
		suff[i] = suff[i+1]
		diff := items[i].b - items[i].a
		if diff > 0 {
			suff[i] += diff
		}
	}

	h := &maxHeap{}
	sum := int64(0)
	ans := int64(0)
	for i := 0; i < n; i++ {
		if len(*h) >= k {
			profit := suff[i] - sum
			if profit > ans {
				ans = profit
			}
		}
		heap.Push(h, items[i].a)
		sum += items[i].a
		if len(*h) > k {
			sum -= heap.Pop(h).(int64)
		}
	}
	return ans
}

func parseLine(line string) (int, int, []int64, []int64, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0, 0, nil, nil, fmt.Errorf("not enough fields")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, nil, nil, err
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, nil, nil, err
	}
	if len(fields) != 2+2*n {
		return 0, 0, nil, nil, fmt.Errorf("expected %d numbers, got %d", 2+2*n, len(fields))
	}
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		v, err := strconv.ParseInt(fields[2+i], 10, 64)
		if err != nil {
			return 0, 0, nil, nil, err
		}
		a[i] = v
	}
	for i := 0; i < n; i++ {
		v, err := strconv.ParseInt(fields[2+n+i], 10, 64)
		if err != nil {
			return 0, 0, nil, nil, err
		}
		b[i] = v
	}
	return n, k, a, b, nil
}

func buildInput(n, k int, a, b []int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	lines := strings.Split(testcasesD, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, k, a, b, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", i+1, err)
			os.Exit(1)
		}
		want := strconv.FormatInt(expected(n, k, a, b), 10)
		input := buildInput(n, k, a, b)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
