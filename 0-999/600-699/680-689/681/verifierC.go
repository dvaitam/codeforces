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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
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

const (
	opInsert = iota + 1
	opGet
	opRemove
)

type command struct {
	t int
	x int
}

func solve(cmds []command) string {
	h := &IntHeap{}
	heap.Init(h)
	var ops []string
	for _, c := range cmds {
		switch c.t {
		case opInsert:
			heap.Push(h, c.x)
			ops = append(ops, fmt.Sprintf("insert %d", c.x))
		case opGet:
			if h.Len() == 0 || (*h)[0] > c.x {
				heap.Push(h, c.x)
				ops = append(ops, fmt.Sprintf("insert %d", c.x))
			} else {
				for h.Len() > 0 && (*h)[0] < c.x {
					heap.Pop(h)
					ops = append(ops, "removeMin")
				}
				if h.Len() == 0 || (*h)[0] > c.x {
					heap.Push(h, c.x)
					ops = append(ops, fmt.Sprintf("insert %d", c.x))
				}
			}
			ops = append(ops, fmt.Sprintf("getMin %d", c.x))
		case opRemove:
			if h.Len() == 0 {
				heap.Push(h, 1)
				ops = append(ops, "insert 1")
			}
			heap.Pop(h)
			ops = append(ops, "removeMin")
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d", len(ops)))
	for _, op := range ops {
		sb.WriteByte('\n')
		sb.WriteString(op)
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	cmds := make([]command, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		typ := rng.Intn(3)
		switch typ {
		case 0:
			x := rng.Intn(50)
			cmds[i] = command{opInsert, x}
			sb.WriteString(fmt.Sprintf("insert %d\n", x))
		case 1:
			x := rng.Intn(50)
			cmds[i] = command{opGet, x}
			sb.WriteString(fmt.Sprintf("getMin %d\n", x))
		default:
			cmds[i] = command{opRemove, 0}
			sb.WriteString("removeMin\n")
		}
	}
	return sb.String(), solve(cmds)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\n\ngot:\n%s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
