package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type opD struct {
	add bool
	l   int
	r   int
}

type MaxHeap struct{ sort.IntSlice }

func (h MaxHeap) Less(i, j int) bool  { return h.IntSlice[i] > h.IntSlice[j] }
func (h *MaxHeap) Push(x interface{}) { h.IntSlice = append(h.IntSlice, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	old := h.IntSlice
	x := old[len(old)-1]
	h.IntSlice = old[:len(old)-1]
	return x
}
func (h *MaxHeap) Peek() int { return h.IntSlice[0] }

type MinHeap struct{ sort.IntSlice }

func (h *MinHeap) Push(x interface{}) { h.IntSlice = append(h.IntSlice, x.(int)) }
func (h *MinHeap) Pop() interface{} {
	old := h.IntSlice
	x := old[len(old)-1]
	h.IntSlice = old[:len(old)-1]
	return x
}
func (h *MinHeap) Peek() int { return h.IntSlice[0] }

func solveCaseD(ops []opD) []string {
	lcnt := make(map[int]int)
	rcnt := make(map[int]int)
	var lMax MaxHeap
	var rMin MinHeap
	heap.Init(&lMax)
	heap.Init(&rMin)
	segments := 0
	res := make([]string, len(ops))
	for i, op := range ops {
		if op.add {
			heap.Push(&lMax, op.l)
			heap.Push(&rMin, op.r)
			lcnt[op.l]++
			rcnt[op.r]++
			segments++
		} else {
			lcnt[op.l]--
			if lcnt[op.l] == 0 {
				delete(lcnt, op.l)
			}
			rcnt[op.r]--
			if rcnt[op.r] == 0 {
				delete(rcnt, op.r)
			}
			segments--
		}
		for lMax.Len() > 0 {
			top := lMax.Peek()
			if lcnt[top] == 0 {
				heap.Pop(&lMax)
			} else {
				break
			}
		}
		for rMin.Len() > 0 {
			top := rMin.Peek()
			if rcnt[top] == 0 {
				heap.Pop(&rMin)
			} else {
				break
			}
		}
		if segments >= 2 && lMax.Peek() > rMin.Peek() {
			res[i] = "YES"
		} else {
			res[i] = "NO"
		}
	}
	return res
}

func runCaseD(bin string, ops []opD) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ops)))
	for _, op := range ops {
		if op.add {
			sb.WriteString(fmt.Sprintf("+ %d %d\n", op.l, op.r))
		} else {
			sb.WriteString(fmt.Sprintf("- %d %d\n", op.l, op.r))
		}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotLines := strings.Fields(strings.TrimSpace(out.String()))
	exp := solveCaseD(ops)
	if len(gotLines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(gotLines))
	}
	for i := range exp {
		if strings.ToUpper(gotLines[i]) != exp[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, exp[i], gotLines[i])
		}
	}
	return nil
}

func randomCaseD(rng *rand.Rand) []opD {
	q := rng.Intn(10) + 1
	ops := make([]opD, 0, q)
	active := make([]opD, 0)
	for len(ops) < q {
		if len(active) == 0 || rng.Intn(2) == 0 {
			l := rng.Intn(20)
			r := l + rng.Intn(20)
			op := opD{add: true, l: l, r: r}
			active = append(active, op)
			ops = append(ops, op)
		} else {
			idx := rng.Intn(len(active))
			op := active[idx]
			active = append(active[:idx], active[idx+1:]...)
			ops = append(ops, opD{add: false, l: op.l, r: op.r})
		}
	}
	return ops
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := [][]opD{{{add: true, l: 1, r: 2}, {add: true, l: 3, r: 4}}}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCaseD(rng))
	}
	for idx, ops := range cases {
		if err := runCaseD(bin, ops); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
