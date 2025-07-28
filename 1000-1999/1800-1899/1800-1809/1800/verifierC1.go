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

func solveC1(arr []int) string {
	h := &MaxHeap{}
	heap.Init(h)
	var ans int64
	for _, x := range arr {
		if x == 0 {
			if h.Len() > 0 {
				ans += int64(h.Peek())
				heap.Pop(h)
			}
		} else {
			heap.Push(h, x)
		}
	}
	return fmt.Sprint(ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := range arr {
		if rng.Intn(2) == 0 {
			arr[i] = 0
		} else {
			arr[i] = rng.Intn(100) + 1
		}
	}
	sb := strings.Builder{}
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprint(n))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	expected := solveC1(arr)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
