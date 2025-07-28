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

const inf int64 = 1 << 60

type item struct {
	t    int64
	i, j int
}

type priorityQueue []item

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].t < pq[j].t }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(item)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

func nextAfter(arr []int64, t int64) int64 {
	idx := sort.Search(len(arr), func(i int) bool { return arr[i] > t })
	if idx == len(arr) {
		return inf
	}
	return arr[idx]
}

func isShot(arr []int64, t int64) bool {
	idx := sort.Search(len(arr), func(i int) bool { return arr[i] >= t })
	return idx < len(arr) && arr[idx] == t
}

func earliest(t int64, i, j, ni, nj int, rowShots, colShots [][]int64) (int64, bool) {
	deadline := nextAfter(rowShots[i], t)
	tmp := nextAfter(colShots[j], t)
	if tmp < deadline {
		deadline = tmp
	}
	cand := t + 1
	for cand <= deadline {
		if !isShot(rowShots[ni], cand) && !isShot(colShots[nj], cand) {
			return cand, true
		}
		cand++
	}
	return 0, false
}

func solveF(n, m int, shots [][3]int64) int64 {
	rowShots := make([][]int64, n+1)
	colShots := make([][]int64, m+1)
	for _, sh := range shots {
		tt := sh[0]
		d := sh[1]
		coord := sh[2]
		if d == 1 {
			rowShots[coord] = append(rowShots[coord], tt)
		} else {
			colShots[coord] = append(colShots[coord], tt)
		}
	}
	for i := 0; i <= n; i++ {
		sort.Slice(rowShots[i], func(a, b int) bool { return rowShots[i][a] < rowShots[i][b] })
	}
	for j := 0; j <= m; j++ {
		sort.Slice(colShots[j], func(a, b int) bool { return colShots[j][a] < colShots[j][b] })
	}
	dist := make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		dist[i] = make([]int64, m+1)
		for j := 0; j <= m; j++ {
			dist[i][j] = inf
		}
	}
	dist[0][0] = 0
	pq := &priorityQueue{}
	heap.Push(pq, item{0, 0, 0})
	ans := int64(-1)
	for pq.Len() > 0 {
		it := heap.Pop(pq).(item)
		t, i, j := it.t, it.i, it.j
		if t != dist[i][j] {
			continue
		}
		if i == n && j == m {
			ans = t
			break
		}
		if i < n {
			if nt, ok := earliest(t, i, j, i+1, j, rowShots, colShots); ok && nt < dist[i+1][j] {
				dist[i+1][j] = nt
				heap.Push(pq, item{nt, i + 1, j})
			}
		}
		if j < m {
			if nt, ok := earliest(t, i, j, i, j+1, rowShots, colShots); ok && nt < dist[i][j+1] {
				dist[i][j+1] = nt
				heap.Push(pq, item{nt, i, j + 1})
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	r := rng.Intn(4)
	shots := make([][3]int64, r)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	sb.WriteString(fmt.Sprintf("%d\n", r))
	for i := 0; i < r; i++ {
		t := int64(rng.Intn(10) + 1)
		d := int64(rng.Intn(2) + 1)
		var coord int64
		if d == 1 {
			coord = int64(rng.Intn(n + 1))
		} else {
			coord = int64(rng.Intn(m + 1))
		}
		shots[i] = [3]int64{t, d, coord}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", t, d, coord))
	}
	expected := fmt.Sprintf("%d", solveF(n, m, shots))
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
