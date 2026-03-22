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
)

type RefItem struct {
	next int
	book int
}

type RefMaxHeap []RefItem

func (h RefMaxHeap) Len() int            { return len(h) }
func (h RefMaxHeap) Less(i, j int) bool  { return h[i].next > h[j].next }
func (h RefMaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *RefMaxHeap) Push(x interface{}) { *h = append(*h, x.(RefItem)) }
func (h *RefMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solveReference(input string) string {
	fields := strings.Fields(input)
	if len(fields) < 2 {
		return ""
	}
	n, _ := strconv.Atoi(fields[0])
	k, _ := strconv.Atoi(fields[1])
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i], _ = strconv.Atoi(fields[2+i])
	}

	next := make([]int, n)
	last := make(map[int]int)
	for i := 0; i < n; i++ {
		last[a[i]] = n
	}
	for i := n - 1; i >= 0; i-- {
		next[i] = last[a[i]]
		last[a[i]] = i
	}

	library := make(map[int]int)
	pq := &RefMaxHeap{}
	heap.Init(pq)
	cost := 0

	for i := 0; i < n; i++ {
		b := a[i]
		nxt := next[i]
		if _, ok := library[b]; ok {
			library[b] = nxt
			heap.Push(pq, RefItem{nxt, b})
			continue
		}
		cost++
		if len(library) >= k {
			for pq.Len() > 0 {
				item := heap.Pop(pq).(RefItem)
				if cur, ok := library[item.book]; ok && cur == item.next {
					delete(library, item.book)
					break
				}
			}
		}
		if k > 0 {
			library[b] = nxt
			heap.Push(pq, RefItem{nxt, b})
		}
	}

	return strconv.Itoa(cost)
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func joinInts(a []int) string {
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func genTests() []string {
	rand.Seed(2)
	tests := make([]string, 0, 102)
	for i := 0; i < 100; i++ {
		n := rand.Intn(50) + 1
		k := rand.Intn(50) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(n) + 1
		}
		tests = append(tests, fmt.Sprintf("%d %d\n%s\n", n, k, joinInts(arr)))
	}
	tests = append(tests, "1 1\n1\n")
	tests = append(tests, "3 1\n1 2 3\n")
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := genTests()
	for i, in := range tests {
		exp := solveReference(strings.TrimSpace(in))
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
