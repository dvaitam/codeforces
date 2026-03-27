package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ── embedded reference solver (correct, based on accepted leftist-heap approach) ──

var hv, htag []int64
var hl, hr, hdist []int
var hcnt int

func initHeap(capacity int) {
	hv = make([]int64, capacity)
	htag = make([]int64, capacity)
	hl = make([]int, capacity)
	hr = make([]int, capacity)
	hdist = make([]int, capacity)
	hcnt = 0
}

func newHeapNode(v int64) int {
	hcnt++
	hv[hcnt] = v
	htag[hcnt] = 0
	hl[hcnt] = 0
	hr[hcnt] = 0
	hdist[hcnt] = 1
	return hcnt
}

func applyTag(x int, d int64) {
	if x == 0 {
		return
	}
	hv[x] += d
	htag[x] += d
}

func pushDown(x int) {
	if x == 0 || htag[x] == 0 {
		return
	}
	d := htag[x]
	applyTag(hl[x], d)
	applyTag(hr[x], d)
	htag[x] = 0
}

func mergeHeaps(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	if hv[a] < hv[b] {
		a, b = b, a
	}
	pushDown(a)
	hr[a] = mergeHeaps(hr[a], b)
	if hdist[hl[a]] < hdist[hr[a]] {
		hl[a], hr[a] = hr[a], hl[a]
	}
	hdist[a] = hdist[hr[a]] + 1
	return a
}

func popMaxHeap(x int) (int, int) {
	pushDown(x)
	a, b := hl[x], hr[x]
	hl[x], hr[x] = 0, 0
	hdist[x] = 1
	htag[x] = 0
	return mergeHeaps(a, b), x
}

func referenceSolve(n, k int, parent []int) int64 {
	head := make([]int, n+1)
	to := make([]int, n+1)
	nxt := make([]int, n+1)
	edge := 1
	for i := 2; i <= n; i++ {
		to[edge] = i
		nxt[edge] = head[parent[i-2]]
		head[parent[i-2]] = edge
		edge++
	}

	order := make([]int, 0, n)
	stack := make([]int, 0, n)
	stack = append(stack, 1)
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for e := head[u]; e != 0; e = nxt[e] {
			stack = append(stack, to[e])
		}
	}

	initHeap(n + 5)
	rootHeap := make([]int, n+1)

	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		h := 0
		for e := head[u]; e != 0; e = nxt[e] {
			h = mergeHeaps(h, rootHeap[to[e]])
		}
		if u != 1 {
			if h == 0 {
				h = newHeapNode(1)
			} else {
				var node int
				h, node = popMaxHeap(h)
				if h != 0 {
					applyTag(h, -1)
					if hv[h] <= 0 {
						h = 0
					}
				}
				hv[node]++
				h = mergeHeaps(h, node)
			}
		}
		rootHeap[u] = h
	}

	h := rootHeap[1]
	saved := int64(0)
	limit := int64(k) + 1
	for limit > 0 && h != 0 && hv[h] > 0 {
		var node int
		h, node = popMaxHeap(h)
		saved += hv[node]
		limit--
	}

	return int64(2*(n-1)) - saved
}

// ── verifier harness ───────────────────────────────────────────────────────

func generateCase(rng *rand.Rand) (int, int, []int) {
	n := rng.Intn(10) + 2
	k := rng.Intn(10)
	parent := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parent[i-2] = rng.Intn(i-1) + 1
	}
	return n, k, parent
}

func runCase(bin string, n, k int, parent []int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, p := range parent {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", p))
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprintf("%d", referenceSolve(n, k, parent))
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
		n, k, parent := generateCase(rng)
		if err := runCase(bin, n, k, parent); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed (n=%d k=%d parent=%v): %v\n", i+1, n, k, parent, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
