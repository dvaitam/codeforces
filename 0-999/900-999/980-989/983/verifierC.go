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

const base = 5

var pow5 [10]int32

func init() {
	pow5[0] = 1
	for i := 1; i < 10; i++ {
		pow5[i] = pow5[i-1] * base
	}
}

type State struct {
	t     int
	idx   int16
	floor int8
	code  int32
	size  int8
}

type PQ []State

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].t < pq[j].t }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(State)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func add(code int32, floor int8, delta int8) int32 {
	return code + int32(delta)*pow5[floor-1]
}

func get(code int32, floor int8) int8 {
	return int8((code / pow5[floor-1]) % base)
}

func makeKey(idx int16, floor int8, code int32) int64 {
	return int64(idx)<<32 | int64(floor)<<28 | int64(code)
}

func solveC(n int, a, b []int8) int {
	start := State{t: 0, idx: 0, floor: 1, code: 0, size: 0}
	pq := &PQ{start}
	heap.Init(pq)
	dist := make(map[int64]int)
	dist[makeKey(start.idx, start.floor, start.code)] = 0
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(State)
		key := makeKey(cur.idx, cur.floor, cur.code)
		if d, ok := dist[key]; ok && d < cur.t {
			continue
		}
		if int(cur.idx) == n && cur.size == 0 {
			return cur.t
		}
		if cur.floor < 9 {
			next := cur
			next.t++
			next.floor++
			k := makeKey(next.idx, next.floor, next.code)
			if d, ok := dist[k]; !ok || next.t < d {
				dist[k] = next.t
				heap.Push(pq, next)
			}
		}
		if cur.floor > 1 {
			next := cur
			next.t++
			next.floor--
			k := makeKey(next.idx, next.floor, next.code)
			if d, ok := dist[k]; !ok || next.t < d {
				dist[k] = next.t
				heap.Push(pq, next)
			}
		}
		exits := get(cur.code, cur.floor)
		idx := cur.idx
		code := cur.code
		size := cur.size
		cost := 0
		if exits > 0 {
			code = add(code, cur.floor, -exits)
			size -= exits
			cost += int(exits)
		}
		for int(idx) < n && a[int(idx)+1] == cur.floor && size < 4 {
			idx++
			code = add(code, b[int(idx)], 1)
			size++
			cost++
		}
		if cost > 0 {
			next := State{t: cur.t + cost, idx: idx, floor: cur.floor, code: code, size: size}
			k := makeKey(next.idx, next.floor, next.code)
			if d, ok := dist[k]; !ok || next.t < d {
				dist[k] = next.t
				heap.Push(pq, next)
			}
		}
	}
	return -1
}

func generateC(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	a := make([]int8, n+1)
	b := make([]int8, n+1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 1; i <= n; i++ {
		ai := int8(rng.Intn(9) + 1)
		bi := int8(rng.Intn(9) + 1)
		for bi == ai {
			bi = int8(rng.Intn(9) + 1)
		}
		a[i] = ai
		b[i] = bi
		fmt.Fprintf(&sb, "%d %d\n", ai, bi)
	}
	res := solveC(n, a, b)
	return sb.String(), fmt.Sprint(res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(44))
	for i := 0; i < 100; i++ {
		in, exp := generateC(rng)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(in)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", i+1, err, out.String())
			return
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
