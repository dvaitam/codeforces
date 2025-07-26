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

type Event struct {
	t int64
	a int64
	b int64
}

type MinHeap []int64

type MaxHeap []int64

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solve(n int, m, c, c0 int64, events []Event) int64 {
	total := c0
	for _, e := range events {
		total += e.a
	}
	if total < m {
		return -1
	}
	sort.Slice(events, func(i, j int) bool { return events[i].t < events[j].t })
	events = append(events, Event{t: m})

	minH := &MinHeap{}
	maxH := &MaxHeap{}
	heap.Init(minH)
	heap.Init(maxH)
	counts := make(map[int64]int64)

	counts[0] = c0
	heap.Push(minH, int64(0))
	heap.Push(maxH, int64(0))

	volume := c0
	var costSum int64
	prev := int64(0)
	possible := true

	consume := func(need int64) bool {
		for need > 0 {
			if minH.Len() == 0 {
				return false
			}
			cost := heap.Pop(minH).(int64)
			cnt := counts[cost]
			if cnt == 0 {
				continue
			}
			use := need
			if cnt < use {
				use = cnt
			}
			cnt -= use
			counts[cost] = cnt
			need -= use
			volume -= use
			costSum += cost * use
			if cnt > 0 {
				heap.Push(minH, cost)
				heap.Push(maxH, cost)
			}
		}
		return true
	}

	discard := func(rem int64) {
		for rem > 0 {
			if maxH.Len() == 0 {
				break
			}
			cost := heap.Pop(maxH).(int64)
			cnt := counts[cost]
			if cnt == 0 {
				continue
			}
			use := rem
			if cnt < use {
				use = cnt
			}
			cnt -= use
			counts[cost] = cnt
			rem -= use
			volume -= use
			if cnt > 0 {
				heap.Push(minH, cost)
				heap.Push(maxH, cost)
			}
		}
	}

	for _, e := range events {
		delta := e.t - prev
		if delta > 0 {
			if volume < delta {
				possible = false
				break
			}
			if !consume(delta) {
				possible = false
				break
			}
		}
		if e.t == m {
			break
		}
		counts[e.b] += e.a
		heap.Push(minH, e.b)
		heap.Push(maxH, e.b)
		volume += e.a
		if volume > c {
			discard(volume - c)
		}
		prev = e.t
	}

	if possible {
		return costSum
	}
	return -1
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

func generateTests() []struct {
	n        int
	m, c, c0 int64
	events   []Event
} {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]struct {
		n        int
		m, c, c0 int64
		events   []Event
	}, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 1
		m := int64(rng.Intn(20) + 5)
		c := int64(rng.Intn(10) + 5)
		if c > m {
			c = m
		}
		c0 := int64(rng.Intn(int(c) + 1))
		events := make([]Event, n)
		curT := int64(0)
		for j := 0; j < n; j++ {
			curT += int64(rng.Intn(5) + 1)
			a := int64(rng.Intn(5) + 1)
			b := int64(rng.Intn(5) + 1)
			events[j] = Event{t: curT, a: a, b: b}
		}
		tests = append(tests, struct {
			n        int
			m, c, c0 int64
			events   []Event
		}{n: n, m: m, c: c, c0: c0, events: events})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("1\n%d %d %d %d\n", t.n, t.m, t.c, t.c0)
		for _, e := range t.events {
			input += fmt.Sprintf("%d %d %d\n", e.t, e.a, e.b)
		}
		want := fmt.Sprintf("%d", solve(t.n, t.m, t.c, t.c0, append([]Event(nil), t.events...)))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\ninput:%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
