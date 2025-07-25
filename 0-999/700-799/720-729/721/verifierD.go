package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n, m int
	x    int64
	arr  []int64
}

type item struct {
	v   int64
	pos int
}

type maxHeap []item

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].v > h[j].v }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.x))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expected(tc testCase) string {
	a := append([]int64(nil), tc.arr...)
	h := &maxHeap{}
	heap.Init(h)
	num := 0
	for i, v := range a {
		if v < 0 {
			num++
		}
		heap.Push(h, item{v: abs(v), pos: i})
	}
	for k := 0; k < tc.m; k++ {
		it := heap.Pop(h).(item)
		p := it.pos
		if a[p] < 0 {
			num--
		}
		if num&1 == 1 {
			a[p] += tc.x
		} else {
			a[p] -= tc.x
		}
		if a[p] < 0 {
			num++
		}
		heap.Push(h, item{v: abs(a[p]), pos: p})
	}
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return strings.TrimSpace(sb.String())
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	x := int64(rng.Intn(5) + 1)
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(21) - 10)
	}
	return testCase{n: n, m: m, x: x, arr: arr}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{}
	for len(cases) < 105 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		input := buildInput(tc)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		exp := expected(tc)
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
