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

type person struct {
	cnt int
	idx int
}

type pHeap []person

func (h pHeap) Len() int            { return len(h) }
func (h pHeap) Less(i, j int) bool  { return h[i].cnt > h[j].cnt }
func (h pHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *pHeap) Push(x interface{}) { *h = append(*h, x.(person)) }
func (h *pHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func maxPairs(arr []int) int {
	h := &pHeap{}
	heap.Init(h)
	for i, v := range arr {
		if v > 0 {
			heap.Push(h, person{v, i})
		}
	}
	pairs := 0
	for h.Len() >= 2 {
		a := heap.Pop(h).(person)
		b := heap.Pop(h).(person)
		pairs++
		a.cnt--
		b.cnt--
		if a.cnt > 0 {
			heap.Push(h, a)
		}
		if b.cnt > 0 {
			heap.Push(h, b)
		}
	}
	return pairs
}

func runCaseD(bin string, arr []int) error {
	n := len(arr)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid k")
	}
	if len(fields) != 1+2*k {
		return fmt.Errorf("expected %d numbers got %d", 1+2*k, len(fields))
	}
	expectK := maxPairs(arr)
	if k != expectK {
		return fmt.Errorf("wrong k: expected %d got %d", expectK, k)
	}
	counts := append([]int(nil), arr...)
	idx := 1
	for i := 0; i < k; i++ {
		aIdx, _ := strconv.Atoi(fields[idx])
		bIdx, _ := strconv.Atoi(fields[idx+1])
		idx += 2
		if aIdx < 1 || aIdx > n || bIdx < 1 || bIdx > n || aIdx == bIdx {
			return fmt.Errorf("invalid pair")
		}
		counts[aIdx-1]--
		counts[bIdx-1]--
		if counts[aIdx-1] < 0 || counts[bIdx-1] < 0 {
			return fmt.Errorf("negative count")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 2
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(5)
		}
		if err := runCaseD(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
