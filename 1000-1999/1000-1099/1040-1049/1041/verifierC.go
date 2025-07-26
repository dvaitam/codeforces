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
)

type Item struct {
	release int
	id      int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].release < pq[j].release }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

func solveC(N, M, D int, arr []int) (int, []int) {
	type pair struct{ t, idx int }
	arr2 := make([]pair, N)
	for i := 0; i < N; i++ {
		arr2[i] = pair{t: arr[i], idx: i}
	}
	sort.Slice(arr2, func(i, j int) bool { return arr2[i].t < arr2[j].t })
	ans := make([]int, N)
	var pq PriorityQueue
	heap.Init(&pq)
	tables := 0
	for _, e := range arr2 {
		t, i := e.t, e.idx
		if pq.Len() > 0 && pq[0].release <= t {
			item := heap.Pop(&pq).(Item)
			ans[i] = item.id
			heap.Push(&pq, Item{release: t + D + 1, id: item.id})
		} else {
			tables++
			ans[i] = tables
			heap.Push(&pq, Item{release: t + D + 1, id: tables})
		}
	}
	return tables, ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(3)
	for t := 1; t <= 100; t++ {
		N := rand.Intn(20) + 1
		M := rand.Intn(100) + N
		D := rand.Intn(M) + 1
		set := make(map[int]bool)
		arr := make([]int, N)
		for i := 0; i < N; i++ {
			v := rand.Intn(M) + 1
			for set[v] {
				v = rand.Intn(M) + 1
			}
			arr[i] = v
			set[v] = true
		}
		input := fmt.Sprintf("%d %d %d\n", N, M, D)
		for i := 0; i < N; i++ {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", arr[i])
		}
		input += "\n"
		tables, ans := solveC(N, M, D, append([]int(nil), arr...))
		expect := fmt.Sprintf("%d\n", tables)
		for i, v := range ans {
			if i+1 < len(ans) {
				expect += fmt.Sprintf("%d ", v)
			} else {
				expect += fmt.Sprintf("%d", v)
			}
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed:\nexpected:\n%s\n got:\n%s\n", t, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
