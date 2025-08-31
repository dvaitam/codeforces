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

const modD = 1000000009

type testCaseD struct {
	m  int
	xs []int
	ys []int
}

func generateCase(rng *rand.Rand) (string, testCaseD) {
	m := rng.Intn(8) + 1
	xs := make([]int, m)
	ys := make([]int, m)
	pos := map[[2]int]bool{}
	xs[0], ys[0] = 0, 0
	pos[[2]int{0, 0}] = true
	for i := 1; i < m; i++ {
		for {
			base := rng.Intn(i) // pick existing cube index
			bx, by := xs[base], ys[base]
			dx := rng.Intn(3) - 1 // -1,0,1
			x := bx + dx
			y := by + 1
			if pos[[2]int{x, y}] {
				continue
			}
			// ensure at least one supporter exists (bx,by)
			xs[i] = x
			ys[i] = y
			pos[[2]int{x, y}] = true
			break
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		fmt.Fprintf(&sb, "%d %d\n", xs[i], ys[i])
	}
	return sb.String(), testCaseD{m: m, xs: xs, ys: ys}
}

// heaps

type MinHeap []int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func expected(tc testCaseD) (int64, error) {
	m := tc.m
	xs := tc.xs
	ys := tc.ys
	idx := make(map[uint64]int, m)
	for i := 0; i < m; i++ {
		key := (uint64(xs[i]) << 32) | uint64(uint32(ys[i]))
		idx[key] = i
	}
	supporters := make([][]int, m)
	dependents := make([][]int, m)
	supportCount := make([]int, m)
	for i := 0; i < m; i++ {
		x, y := xs[i], ys[i]
		for dx := -1; dx <= 1; dx++ {
			key := (uint64(x+dx) << 32) | uint64(uint32(y-1))
			if j, ok := idx[key]; ok {
				supporters[i] = append(supporters[i], j)
			}
		}
		supportCount[i] = len(supporters[i])
		for _, j := range supporters[i] {
			dependents[j] = append(dependents[j], i)
		}
	}
	inSet := make([]bool, m)
	removed := make([]bool, m)
	var minh MinHeap
	var maxh MaxHeap
	heap.Init(&minh)
	heap.Init(&maxh)
	for i := 0; i < m; i++ {
		ok := true
		for _, u := range dependents[i] {
			if !removed[u] && supportCount[u] < 2 {
				ok = false
				break
			}
		}
		if ok {
			inSet[i] = true
			heap.Push(&minh, i)
			heap.Push(&maxh, i)
		}
	}
	refillMax := func() {
		for i := 0; i < m; i++ {
			if !removed[i] && inSet[i] {
				heap.Push(&maxh, i)
			}
		}
	}
	refillMin := func() {
		for i := 0; i < m; i++ {
			if !removed[i] && inSet[i] {
				heap.Push(&minh, i)
			}
		}
	}
	popMax := func() (int, bool) {
		for {
			for maxh.Len() > 0 {
				top := heap.Pop(&maxh).(int)
				if !removed[top] && inSet[top] {
					return top, true
				}
			}
			refillMax()
			if maxh.Len() == 0 {
				return 0, false
			}
		}
	}
	popMin := func() (int, bool) {
		for {
			for minh.Len() > 0 {
				top := heap.Pop(&minh).(int)
				if !removed[top] && inSet[top] {
					return top, true
				}
			}
			refillMin()
			if minh.Len() == 0 {
				return 0, false
			}
		}
	}
	seq := make([]int, 0, m)
	for t := 0; t < m; t++ {
		var v int
		var ok bool
		if t%2 == 0 {
			v, ok = popMax()
		} else {
			v, ok = popMin()
		}
		if !ok {
			return 0, fmt.Errorf("no removable cube available at step %d", t)
		}
		removed[v] = true
		inSet[v] = false
		seq = append(seq, v)
		// Update dependents of removed cube v: always decrement, then if it drops to 1,
		// the unique remaining supporter becomes non-removable.
		for _, u := range dependents[v] {
			if removed[u] {
				continue
			}
			supportCount[u]--
			if supportCount[u] == 1 {
				for _, w := range supporters[u] {
					if !removed[w] {
						if inSet[w] {
							inSet[w] = false
						}
						break
					}
				}
			}
		}
		// Try to add supporters of v that may have become removable now.
		for _, w := range supporters[v] {
			if removed[w] || inSet[w] {
				continue
			}
			ok2 := true
			for _, u := range dependents[w] {
				if !removed[u] && supportCount[u] < 2 {
					ok2 = false
					break
				}
			}
			if ok2 {
				inSet[w] = true
				heap.Push(&minh, w)
				heap.Push(&maxh, w)
			}
		}
	}
	var res int64
	for i := 0; i < m; i++ {
		res = (res*int64(m) + int64(seq[i])) % modD
	}
	return res, nil
}

func runCase(bin string, input string, tc testCaseD) error {
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
	gotStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(gotStr, &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	want, err := expected(tc)
	if err != nil {
		return fmt.Errorf("verifier error: %v", err)
	}
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		if err := runCase(bin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
