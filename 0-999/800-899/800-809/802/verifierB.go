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
	id   int
}

type RefMaxHeap []RefItem

func (h RefMaxHeap) Len() int { return len(h) }
func (h RefMaxHeap) Less(i, j int) bool {
	if h[i].next == h[j].next {
		return h[i].id > h[j].id
	}
	return h[i].next > h[j].next
}
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
	q := n
	a := make([]int, q)
	for i := 0; i < q; i++ {
		a[i], _ = strconv.Atoi(fields[2+i])
	}

	inf := q + 1
	nextPos := make([]int, q)
	last := make(map[int]int, q*2)
	for i := q - 1; i >= 0; i-- {
		x := a[i]
		if p, ok := last[x]; ok {
			nextPos[i] = p
		} else {
			nextPos[i] = inf
		}
		last[x] = i
	}

	cur := make(map[int]int, k*2+1)
	h := &RefMaxHeap{}
	heap.Init(h)
	misses := 0

	for i, x := range a {
		nxt := nextPos[i]
		if _, ok := cur[x]; ok {
			cur[x] = nxt
			heap.Push(h, RefItem{next: nxt, id: x})
			continue
		}

		misses++
		if len(cur) == k {
			for h.Len() > 0 {
				it := heap.Pop(h).(RefItem)
				if v, ok := cur[it.id]; ok && v == it.next {
					delete(cur, it.id)
					break
				}
			}
		}
		if k > 0 {
			cur[x] = nxt
			heap.Push(h, RefItem{next: nxt, id: x})
		}
	}

	return strconv.Itoa(misses)
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
