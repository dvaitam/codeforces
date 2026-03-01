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

type state struct {
	time int
	mask int
	pos  int
}

type minHeap []state

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].time < h[j].time }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(state)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solveCaseExact(m int, top, bottom []int) int {
	n := 2 * m
	unlock := make([]int, n)
	for c := 0; c < m; c++ {
		unlock[c] = top[c]
		unlock[m+c] = bottom[c]
	}

	adj := make([][]int, n)
	for c := 0; c < m; c++ {
		uTop := c
		uBottom := m + c
		adj[uTop] = append(adj[uTop], uBottom)
		adj[uBottom] = append(adj[uBottom], uTop)
		if c > 0 {
			adj[uTop] = append(adj[uTop], c-1)
			adj[uBottom] = append(adj[uBottom], m+c-1)
		}
		if c+1 < m {
			adj[uTop] = append(adj[uTop], c+1)
			adj[uBottom] = append(adj[uBottom], m+c+1)
		}
	}

	const inf = int(1e18)
	dist := make([][]int, 1<<n)
	for i := range dist {
		dist[i] = make([]int, n)
		for j := range dist[i] {
			dist[i][j] = inf
		}
	}

	startPos := 0 // (1,1) => top row, first column
	startMask := 1 << startPos
	dist[startMask][startPos] = 0

	pq := &minHeap{{time: 0, mask: startMask, pos: startPos}}
	heap.Init(pq)

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(state)
		if cur.time != dist[cur.mask][cur.pos] {
			continue
		}
		for _, nb := range adj[cur.pos] {
			if (cur.mask>>nb)&1 == 1 {
				continue
			}
			nextMask := cur.mask | (1 << nb)
			nextTime := cur.time + 1
			if nextTime < unlock[nb] {
				nextTime = unlock[nb]
			}
			if nextTime < dist[nextMask][nb] {
				dist[nextMask][nb] = nextTime
				heap.Push(pq, state{time: nextTime, mask: nextMask, pos: nb})
			}
		}
	}

	fullMask := (1 << n) - 1
	ans := inf
	for pos := 0; pos < n; pos++ {
		if dist[fullMask][pos] < ans {
			ans = dist[fullMask][pos]
		}
	}
	return ans
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func verifyCase(bin string, m int, top, bottom []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(top[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(bottom[i]))
	}
	sb.WriteByte('\n')
	expected := fmt.Sprint(solveCaseExact(m, top, bottom))
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != expected {
		return fmt.Errorf("expected %s got %s", expected, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		m := rng.Intn(6) + 2 // exact verifier; keep small enough for 2^(2m) state space
		top := make([]int, m)
		bottom := make([]int, m)
		top[0] = 0
		for j := 1; j < m; j++ {
			top[j] = rng.Intn(1000000000)
		}
		for j := 0; j < m; j++ {
			bottom[j] = rng.Intn(1000000000)
		}
		if err := verifyCase(bin, m, top, bottom); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
