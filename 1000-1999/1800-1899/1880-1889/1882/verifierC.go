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

func runCandidate(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type MinHeap []int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func solveCase(arr []int) string {
	h := &MinHeap{}
	heap.Init(h)
	posOdd := 0
	var score int64
	for i, v := range arr {
		if i%2 == 0 { // odd index (1-based)
			if v >= 0 {
				score += int64(v)
				posOdd++
			} else {
				heap.Push(h, -v)
			}
		} else {
			if v > 0 {
				if posOdd > 0 {
					score += int64(v)
					posOdd--
				} else if h.Len() > 0 && v > (*h)[0] {
					score += int64(v - (*h)[0])
					heap.Pop(h)
				}
			}
		}
	}
	return fmt.Sprintf("%d", score)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(21) - 10
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), solveCase(arr)
}

func fixedCases() [][2]string {
	return [][2]string{
		{"1\n1\n5\n", "5"},
		{"1\n3\n1 -2 3\n", "4"},
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range fixedCases() {
		out, err := runCandidate(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "fixed case %d failed: %v\ninput:\n%s", idx+1, err, tc[0])
			os.Exit(1)
		}
		if out != tc[1] {
			fmt.Fprintf(os.Stderr, "fixed case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
