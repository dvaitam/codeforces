package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type segment struct{ l, r int }

type segHeap []segment

func (h segHeap) Len() int { return len(h) }
func (h segHeap) Less(i, j int) bool {
	lenI := h[i].r - h[i].l + 1
	lenJ := h[j].r - h[j].l + 1
	if lenI == lenJ {
		return h[i].l < h[j].l
	}
	return lenI > lenJ
}
func (h segHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *segHeap) Push(x interface{}) { *h = append(*h, x.(segment)) }
func (h *segHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type testD struct{ n int }

func genTestsD() []testD {
	rand.Seed(4353)
	tests := make([]testD, 0, 100)
	tests = append(tests, testD{1}, testD{2}, testD{3}, testD{10})
	for len(tests) < 100 {
		n := rand.Intn(100) + 1
		tests = append(tests, testD{n})
	}
	return tests[:100]
}

func expectedD(n int) []int {
	ans := make([]int, n+1)
	h := &segHeap{{1, n}}
	heap.Init(h)
	for i := 1; i <= n; i++ {
		seg := heap.Pop(h).(segment)
		l, r := seg.l, seg.r
		var mid int
		length := r - l + 1
		if length%2 == 1 {
			mid = (l + r) / 2
		} else {
			mid = (l + r - 1) / 2
		}
		ans[mid] = i
		if l <= mid-1 {
			heap.Push(h, segment{l, mid - 1})
		}
		if mid+1 <= r {
			heap.Push(h, segment{mid + 1, r})
		}
	}
	return ans[1:]
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTestsD()
	for idx, tc := range tests {
		input := fmt.Sprintf("1\n%d\n", tc.n)
		expArr := expectedD(tc.n)
		exp := make([]string, len(expArr))
		for i, v := range expArr {
			exp[i] = fmt.Sprintf("%d", v)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", idx+1, err, got)
			os.Exit(1)
		}
		gotFields := strings.Fields(got)
		if len(gotFields) != len(expArr) {
			fmt.Printf("test %d failed: expected %d numbers got %d\n", idx+1, len(expArr), len(gotFields))
			os.Exit(1)
		}
		for i, v := range expArr {
			if gotFields[i] != fmt.Sprintf("%d", v) {
				fmt.Printf("test %d failed at pos %d: expected %d got %s\n", idx+1, i+1, v, gotFields[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
