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

type event struct {
	x   int64
	h   int64
	typ int // 1 start, 0 end
}

type maxHeap []int64

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *maxHeap) Pop() interface{} {
	a := *h
	v := a[len(a)-1]
	*h = a[:len(a)-1]
	return v
}

func envelope(blds [][3]int64) [][2]int64 {
	events := make([]event, 0, len(blds)*2)
	for _, b := range blds {
		h, l, r := b[0], b[1], b[2]
		events = append(events, event{x: l, h: h, typ: 1})
		events = append(events, event{x: r, h: h, typ: 0})
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].x != events[j].x {
			return events[i].x < events[j].x
		}
		if events[i].typ != events[j].typ {
			return events[i].typ > events[j].typ
		}
		if events[i].typ == 1 {
			return events[i].h > events[j].h
		}
		return events[i].h < events[j].h
	})
	hq := &maxHeap{}
	heap.Init(hq)
	heap.Push(hq, 0)
	rem := make(map[int64]int)
	prev := int64(0)
	var keys []event
	for i := 0; i < len(events); {
		x := events[i].x
		j := i
		for j < len(events) && events[j].x == x {
			e := events[j]
			if e.typ == 1 {
				heap.Push(hq, e.h)
			} else {
				rem[e.h]++
			}
			j++
		}
		for hq.Len() > 0 {
			top := (*hq)[0]
			if c, ok := rem[top]; ok && c > 0 {
				heap.Pop(hq)
				rem[top]--
			} else {
				break
			}
		}
		cur := (*hq)[0]
		if cur != prev {
			keys = append(keys, event{x: x, h: cur})
			prev = cur
		}
		i = j
	}
	var verts [][2]int64
	if len(keys) > 0 {
		prevX := keys[0].x
		prevH := int64(0)
		verts = append(verts, [2]int64{prevX, 0})
		for _, k := range keys {
			if k.h != prevH {
				if k.x != prevX {
					verts = append(verts, [2]int64{k.x, prevH})
				}
				verts = append(verts, [2]int64{k.x, k.h})
				prevX, prevH = k.x, k.h
			}
		}
	}
	return verts
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	blds := make([][3]int64, n)
	for i := 0; i < n; i++ {
		h := int64(rng.Intn(10) + 1)
		l := int64(rng.Intn(10))
		r := l + int64(rng.Intn(5)+1)
		blds[i] = [3]int64{h, l, r}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, b := range blds {
		fmt.Fprintf(&sb, "%d %d %d\n", b[0], b[1], b[2])
	}
	input := sb.String()
	verts := envelope(blds)
	var exp strings.Builder
	fmt.Fprintf(&exp, "%d\n", len(verts))
	for _, v := range verts {
		fmt.Fprintf(&exp, "%d %d\n", v[0], v[1])
	}
	return input, strings.TrimSpace(exp.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases [][2]string
	cases = append(cases, [2]string{"1\n1 0 1\n", "2\n0 0\n1 1"})
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, [2]string{in, exp})
	}
	for i, tc := range cases {
		out, err := run(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc[0])
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc[1] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
