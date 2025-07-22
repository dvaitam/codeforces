package main

import (
	"bytes"
	"container/heap"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var n int
	fmt.Fscan(rdr, &n)
	costsMap := make(map[int][]int)
	c0 := 0
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(rdr, &a, &b)
		if a == 0 {
			c0++
		} else {
			costsMap[a] = append(costsMap[a], b)
		}
	}
	maxCount := 0
	buckets := make([][]int, n+1)
	for _, costs := range costsMap {
		cnt := len(costs)
		if cnt > maxCount {
			maxCount = cnt
		}
		sort.Ints(costs)
		buckets[cnt] = append(buckets[cnt], costs...)
	}
	h := &IntHeap{}
	heap.Init(h)
	bribes := 0
	var costSum int64
	ans := int64(1<<63 - 1)
	for d := maxCount; d >= 0; d-- {
		for _, c := range buckets[d] {
			heap.Push(h, c)
		}
		needed := d - (c0 + bribes) + 1
		for needed > 0 && h.Len() > 0 {
			minCost := heap.Pop(h).(int)
			costSum += int64(minCost)
			bribes++
			needed--
		}
		if c0+bribes > d && costSum < ans {
			ans = costSum
		}
	}
	return fmt.Sprintf("%d", ans)
}

func genRandomCase() string {
	n := rand.Intn(6) + 1
	var buf bytes.Buffer
	fmt.Fprintln(&buf, n)
	for i := 0; i < n; i++ {
		a := rand.Intn(3)
		b := rand.Intn(20) + 1
		fmt.Fprintf(&buf, "%d %d\n", a, b)
	}
	return buf.String()
}

func generateCases() []testCase {
	rand.Seed(3)
	cases := []testCase{}
	fixed := []string{
		"1\n0 5\n",
		"3\n0 1\n1 2\n2 3\n",
	}
	for _, f := range fixed {
		cases = append(cases, testCase{f, compute(f)})
	}
	for len(cases) < 100 {
		inp := genRandomCase()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierC.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%sexpected:%s\nactual:%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
