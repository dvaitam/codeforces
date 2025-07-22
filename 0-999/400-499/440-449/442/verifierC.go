package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

type Item struct {
	score int
	idx   int
}

type MaxHeap []Item

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].score > h[j].score }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[0 : n-1]
	return it
}

func solveC(r *bufio.Reader) string {
	var n int
	fmt.Fscan(r, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
	}
	if n < 3 {
		return "0\n"
	}
	L := make([]int, n)
	R := make([]int, n)
	removed := make([]bool, n)
	for i := 0; i < n; i++ {
		L[i] = i - 1
		R[i] = i + 1
	}
	R[n-1] = -1
	getScore := func(i int) int {
		if i < 0 || i >= n || removed[i] {
			return -1
		}
		l := L[i]
		r := R[i]
		if l < 0 || r < 0 {
			return -1
		}
		if a[l] < a[r] {
			return a[l]
		}
		return a[r]
	}
	h := &MaxHeap{}
	heap.Init(h)
	for i := 1; i < n-1; i++ {
		sc := getScore(i)
		if sc >= 0 {
			heap.Push(h, Item{score: sc, idx: i})
		}
	}
	var total int64
	for h.Len() > 0 {
		it := heap.Pop(h).(Item)
		i := it.idx
		if removed[i] {
			continue
		}
		sc := getScore(i)
		if sc != it.score {
			continue
		}
		total += int64(sc)
		removed[i] = true
		l := L[i]
		r := R[i]
		if l >= 0 {
			R[l] = r
		}
		if r >= 0 {
			L[r] = l
		}
		if sc2 := getScore(l); sc2 >= 0 {
			heap.Push(h, Item{score: sc2, idx: l})
		}
		if sc3 := getScore(r); sc3 >= 0 {
			heap.Push(h, Item{score: sc3, idx: r})
		}
	}
	return fmt.Sprintf("%d\n", total)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 3
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", rng.Intn(1000)+1)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solveC(bufio.NewReader(strings.NewReader(tc)))
		out, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %sinput:\n%s", i+1, expect, out, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
