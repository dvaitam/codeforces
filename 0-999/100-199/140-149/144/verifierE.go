package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type segment struct {
	l, r int
	idx  int
}

type item struct {
	r   int
	idx int
}

type minHeap []item

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].r < h[j].r }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

func solveGreedy(n int, segs []segment) int {
	// Sort by start time
	sort.Slice(segs, func(i, j int) bool { return segs[i].l < segs[j].l })
	h := &minHeap{}
	heap.Init(h)
	count := 0
	j := 0
	for l := 1; l <= n; l++ {
		// Add intervals starting at l
		for j < len(segs) && segs[j].l == l {
			heap.Push(h, item{r: segs[j].r, idx: segs[j].idx})
			j++
		}
		// Remove expired intervals
		for h.Len() > 0 && (*h)[0].r < l {
			heap.Pop(h)
		}
		// Pick one
		if h.Len() > 0 {
			heap.Pop(h)
			count++
		}
	}
	return count
}

func validate(n int, chosenIndices []int, allSegs []segment) error {
	// Extract chosen segments
	idxMap := make(map[int]segment)
	for _, s := range allSegs {
		idxMap[s.idx] = s
	}

	chosenSegs := make([]segment, 0, len(chosenIndices))
	for _, idx := range chosenIndices {
		if s, ok := idxMap[idx]; ok {
			chosenSegs = append(chosenSegs, s)
		} else {
			return fmt.Errorf("invalid index %d", idx)
		}
	}

	// Verify if chosenSegs can be scheduled
	sort.Slice(chosenSegs, func(i, j int) bool { return chosenSegs[i].l < chosenSegs[j].l })
	h := &minHeap{}
	heap.Init(h)
	j := 0
	scheduled := 0
	for l := 1; l <= n; l++ {
		for j < len(chosenSegs) && chosenSegs[j].l == l {
			heap.Push(h, item{r: chosenSegs[j].r, idx: chosenSegs[j].idx})
			j++
		}
		// Check for deadlines: if top of heap < l, then it's expired and we failed to schedule it
		if h.Len() > 0 && (*h)[0].r < l {
			return fmt.Errorf("cannot schedule sportsman %d (deadline %d < time %d)", (*h)[0].idx, (*h)[0].r, l)
		}
		
		if h.Len() > 0 {
			heap.Pop(h)
			scheduled++
		}
	}
	if h.Len() > 0 {
		return fmt.Errorf("not enough slots for %d sportsmen", h.Len())
	}
	if scheduled != len(chosenIndices) {
		return fmt.Errorf("scheduled %d out of %d", scheduled, len(chosenIndices))
	}
	return nil
}

func genCase(rng *rand.Rand) (string, int, []segment) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	segs := make([]segment, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		var r, c int
		for {
			r = rng.Intn(n) + 1
			c = rng.Intn(n) + 1
			if r+c > n {
				break
			}
		}
		// L = n + 1 - c, R = r
		segs[i] = segment{l: n + 1 - c, r: r, idx: i + 1}
		sb.WriteString(fmt.Sprintf("%d %d\n", r, c))
	}
	return sb.String(), n, segs
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n, segs := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		
		// Parse output
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) < 1 {
			fmt.Fprintf(os.Stderr, "case %d failed: empty output\n", i+1)
			os.Exit(1)
		}
		count, err := strconv.Atoi(strings.TrimSpace(lines[0]))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid count: %v\n", i+1, err)
			os.Exit(1)
		}
		
		var indices []int
		if count > 0 {
			if len(lines) < 2 {
				fmt.Fprintf(os.Stderr, "case %d failed: missing indices line\n", i+1)
				os.Exit(1)
			}
			fields := strings.Fields(lines[1])
			if len(fields) != count {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %d indices, got %d\n", i+1, count, len(fields))
				os.Exit(1)
			}
			for _, f := range fields {
				idx, err := strconv.Atoi(f)
				if err != nil {
					fmt.Fprintf(os.Stderr, "case %d failed: invalid index: %v\n", i+1, err)
					os.Exit(1)
				}
				indices = append(indices, idx)
			}
		}

		// Check optimality
		expCount := solveGreedy(n, append([]segment(nil), segs...)) // pass copy because solveGreedy sorts
		if count != expCount {
			fmt.Fprintf(os.Stderr, "case %d failed: expected count %d, got %d\ninput:\n%s", i+1, expCount, count, in)
			os.Exit(1)
		}

		// Check validity
		if err := validate(n, indices, segs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid set: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}