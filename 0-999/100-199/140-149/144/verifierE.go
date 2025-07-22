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

func expected(n int, segs []segment) string {
	sort.Slice(segs, func(i, j int) bool { return segs[i].l < segs[j].l })
	h := &minHeap{}
	heap.Init(h)
	ans := make([]int, 0, len(segs))
	j := 0
	for l := 1; l <= n; l++ {
		for j < len(segs) && segs[j].l == l {
			heap.Push(h, item{r: segs[j].r, idx: segs[j].idx})
			j++
		}
		for h.Len() > 0 && (*h)[0].r < l {
			heap.Pop(h)
		}
		if h.Len() > 0 {
			it := heap.Pop(h).(item)
			ans = append(ans, it.idx)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d", len(ans)))
	if len(ans) > 0 {
		sb.WriteByte('\n')
		for i, v := range ans {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
	} else {
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) (string, string) {
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
		segs[i] = segment{l: n - r + 1, r: c, idx: i + 1}
		sb.WriteString(fmt.Sprintf("%d %d\n", r, c))
	}
	exp := expected(n, segs)
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
