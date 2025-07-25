package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// IntHeap is a min-heap of ints
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

func run(bin, input string) (string, error) {
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

func expected(n int, t int64, s string) string {
	dot := strings.IndexByte(s, '.')
	intStr := s[:dot]
	fracStr := s[dot+1:]
	intPart := make([]int, len(intStr))
	for i := range intStr {
		intPart[i] = int(intStr[i] - '0')
	}
	fracPart := make([]int, len(fracStr))
	for i := range fracStr {
		fracPart[i] = int(fracStr[i] - '0')
	}
	fracLen := len(fracPart)
	h := &IntHeap{}
	heap.Init(h)
	for i, v := range fracPart {
		if v >= 5 {
			heap.Push(h, i)
		}
	}
	for t > 0 && h.Len() > 0 {
		i := heap.Pop(h).(int)
		if i >= fracLen {
			continue
		}
		t--
		rp := i - 1
		if rp < 0 {
			carry := 1
			for k := len(intPart) - 1; k >= 0 && carry > 0; k-- {
				v := intPart[k] + carry
				intPart[k] = v % 10
				carry = v / 10
			}
			if carry > 0 {
				intPart = append([]int{carry}, intPart...)
			}
			fracLen = 0
			break
		}
		carry := 1
		p := rp
		for k := p; k >= 0 && carry > 0; k-- {
			v := fracPart[k] + carry
			fracPart[k] = v % 10
			carry = v / 10
			p = k
		}
		if carry > 0 {
			carryInt := 1
			for k := len(intPart) - 1; k >= 0 && carryInt > 0; k-- {
				v := intPart[k] + carryInt
				intPart[k] = v % 10
				carryInt = v / 10
			}
			if carryInt > 0 {
				intPart = append([]int{carryInt}, intPart...)
			}
			fracLen = 0
			break
		}
		fracLen = p + 1
		if fracPart[p] >= 5 {
			heap.Push(h, p)
		}
	}
	for fracLen > 0 && fracPart[fracLen-1] == 0 {
		fracLen--
	}
	var sb strings.Builder
	for _, d := range intPart {
		sb.WriteByte(byte('0' + d))
	}
	if fracLen > 0 {
		sb.WriteByte('.')
		for i := 0; i < fracLen; i++ {
			sb.WriteByte(byte('0' + fracPart[i]))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcases:", err)
		os.Exit(1)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	idx := 0
	for {
		if !sc.Scan() {
			break
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var n int
		var t int64
		fmt.Sscan(line, &n, &t)
		if !sc.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing grade\n", idx+1)
			os.Exit(1)
		}
		s := strings.TrimSpace(sc.Text())
		idx++
		exp := expected(n, t, s)
		input := fmt.Sprintf("%d %d\n%s\n", n, t, s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
