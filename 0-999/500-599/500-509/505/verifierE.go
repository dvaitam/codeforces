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

type Item struct {
	id int
	ai int64
}

type MaxHeap []Item

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].ai > h[j].ai }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func expected(n, m, k int, p int64, hArr, aArr []int64) int64 {
	hi := int64(0)
	for i := 0; i < n; i++ {
		endH := hArr[i] + aArr[i]*int64(m)
		if endH > hi {
			hi = endH
		}
	}
	can := func(H int64) bool {
		events := make([][]int, m+2)
		times := make([]int, n)
		for i := 0; i < n; i++ {
			var t0 int
			if hArr[i] > H {
				t0 = 1
			} else {
				t0 = int((H-hArr[i])/aArr[i]) + 2
			}
			if t0 <= m {
				events[t0] = append(events[t0], i)
			}
		}
		hq := &MaxHeap{}
		heap.Init(hq)
		for d := 1; d <= m; d++ {
			for _, idx := range events[d] {
				heap.Push(hq, Item{id: idx, ai: aArr[idx]})
			}
			for j := 0; j < k; j++ {
				if hq.Len() == 0 {
					break
				}
				it := heap.Pop(hq).(Item)
				id := it.id
				times[id]++
				hStart := hArr[id] + aArr[id]*int64(d) - int64(times[id])*p
				var tnext int
				if hStart > H {
					tnext = d + 1
				} else {
					tnext = d + int((H-hStart)/aArr[id]) + 2
				}
				if tnext <= m {
					events[tnext] = append(events[tnext], id)
				}
			}
			if hq.Len() > 0 {
				return false
			}
		}
		return true
	}
	lo := int64(-1)
	for lo+1 < hi {
		mid := (lo + hi) / 2
		if can(mid) {
			hi = mid
		} else {
			lo = mid
		}
	}
	return hi
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		k := rng.Intn(4) + 1
		p := int64(rng.Intn(20) + 1)
		hArr := make([]int64, n)
		aArr := make([]int64, n)
		for i := 0; i < n; i++ {
			hArr[i] = int64(rng.Intn(20))
			aArr[i] = int64(rng.Intn(10) + 1)
		}
		input := fmt.Sprintf("%d %d %d %d\n", n, m, k, p)
		for i := 0; i < n; i++ {
			input += fmt.Sprintf("%d %d\n", hArr[i], aArr[i])
		}
		exp := expected(n, m, k, p, hArr, aArr)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		if out != fmt.Sprintf("%d", exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", tc+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
