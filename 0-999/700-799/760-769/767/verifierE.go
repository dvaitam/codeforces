package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Opt struct {
	v   int64
	num int
}

type OptHeap []Opt

func (h OptHeap) Len() int            { return len(h) }
func (h OptHeap) Less(i, j int) bool  { return h[i].v < h[j].v }
func (h OptHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *OptHeap) Push(x interface{}) { *h = append(*h, x.(Opt)) }
func (h *OptHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type testCase struct {
	n  int
	m  int64
	ai []int64
	w  []int64
}

func solve(n int, m int64, ai []int64, w []int64) string {
	ans1 := make([]int64, n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		ans1[i] = ai[i] / 100
		a[i] = int(ai[i] % 100)
	}
	ans2 := make([]bool, n)
	now := int64(0)
	h := &OptHeap{}
	heap.Init(h)
	for i := 0; i < n; i++ {
		if a[i] != 0 {
			cost := int64(a[i])
			benefit := w[i] * int64(100-a[i])
			if m >= cost {
				m -= cost
				heap.Push(h, Opt{benefit, i})
			} else {
				if h.Len() == 0 {
					now += benefit
					m += int64(100 - a[i])
					ans2[i] = true
				} else {
					top := (*h)[0]
					if top.v < benefit {
						heap.Pop(h)
						heap.Push(h, Opt{benefit, i})
						now += top.v
						ans2[top.num] = true
						m += int64(100 - a[i])
					} else {
						now += benefit
						m += int64(100 - a[i])
						ans2[i] = true
					}
				}
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", now))
	for i := 0; i < n; i++ {
		if ans2[i] {
			sb.WriteString(fmt.Sprintf("%d 0\n", ans1[i]+1))
		} else {
			sb.WriteString(fmt.Sprintf("%d %d\n", ans1[i], a[i]))
		}
	}
	return strings.TrimSpace(sb.String())
}

func (tc testCase) input() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i, v := range tc.ai {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range tc.w {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	m := int64(rng.Intn(500))
	ai := make([]int64, n)
	w := make([]int64, n)
	for i := 0; i < n; i++ {
		ai[i] = int64(rng.Intn(400))
		w[i] = int64(rng.Intn(10) + 1)
	}
	return testCase{n: n, m: m, ai: ai, w: w}
}

func runProgram(bin, input string) (string, error) {
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

func runCase(bin string, tc testCase) error {
	in := tc.input()
	expected := solve(tc.n, tc.m, append([]int64(nil), tc.ai...), append([]int64(nil), tc.w...))
	got, err := runProgram(bin, in)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{randomCase(rng)}
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
