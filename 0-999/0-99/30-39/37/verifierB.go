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

type spell struct {
	pow int
	dmg int
}

type Item struct{ dmg, idx int }
type MaxHeap []Item

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].dmg > h[j].dmg }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type caseB struct {
	n      int
	maxH   int
	reg    int
	spells []spell
}

func generateCase(rng *rand.Rand) caseB {
	n := rng.Intn(5) + 1
	maxH := rng.Intn(50) + 1
	reg := rng.Intn(20) + 1
	sp := make([]spell, n)
	for i := range sp {
		sp[i] = spell{pow: rng.Intn(101), dmg: rng.Intn(20) + 1}
	}
	return caseB{n, maxH, reg, sp}
}

// solver replicating 37B.go
func solveCase(tc caseB) (bool, int, [][2]int) {
	type Spell struct {
		thr      int
		dmg, idx int
	}
	spells := make([]Spell, tc.n)
	for i, s := range tc.spells {
		spells[i] = Spell{thr: s.pow * tc.maxH, dmg: s.dmg, idx: i + 1}
	}
	sort.Slice(spells, func(i, j int) bool { return spells[i].thr > spells[j].thr })
	h := tc.maxH
	D := 0
	actions := make([][2]int, 0, tc.n)
	heapItems := &MaxHeap{}
	heap.Init(heapItems)
	j := 0
	pushAvail := func() {
		for j < len(spells) && spells[j].thr >= 100*h {
			heap.Push(heapItems, Item{spells[j].dmg, spells[j].idx})
			j++
		}
	}
	pushAvail()
	if heapItems.Len() == 0 {
		return false, 0, nil
	}
	it := heap.Pop(heapItems).(Item)
	D += it.dmg
	actions = append(actions, [2]int{0, it.idx})
	t := 1
	for {
		h -= D
		h += tc.reg
		if h > tc.maxH {
			h = tc.maxH
		}
		if h <= 0 {
			break
		}
		pushAvail()
		if heapItems.Len() > 0 {
			it = heap.Pop(heapItems).(Item)
			D += it.dmg
			actions = append(actions, [2]int{t, it.idx})
		} else if D <= tc.reg {
			return false, 0, nil
		}
		t++
		if t > 1000 {
			return false, 0, nil
		} // safety
	}
	return true, t, actions
}

func runCase(bin string, tc caseB) error {
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.maxH, tc.reg)
	for _, sp := range tc.spells {
		fmt.Fprintf(&input, "%d %d\n", sp.pow, sp.dmg)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	ok, tExp, actsExp := solveCase(tc)
	if !ok {
		if len(lines) == 0 || strings.TrimSpace(lines[0]) != "NO" {
			return fmt.Errorf("expected NO, got %s", out.String())
		}
		return nil
	}
	if len(lines) < 2 {
		return fmt.Errorf("incomplete output: %s", out.String())
	}
	if strings.TrimSpace(lines[0]) != "YES" {
		return fmt.Errorf("expected YES")
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return fmt.Errorf("missing second line")
	}
	tGot, err1 := strconv.Atoi(fields[0])
	kGot, err2 := strconv.Atoi(fields[1])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("bad integers on line2")
	}
	if tGot != tExp || kGot != len(actsExp) {
		return fmt.Errorf("expected %d %d got %d %d", tExp, len(actsExp), tGot, kGot)
	}
	if len(lines)-2 != len(actsExp) {
		return fmt.Errorf("expected %d action lines, got %d", len(actsExp), len(lines)-2)
	}
	for i, a := range actsExp {
		f := strings.Fields(lines[2+i])
		if len(f) < 2 {
			return fmt.Errorf("bad action line %d", i+1)
		}
		tg, _ := strconv.Atoi(f[0])
		idx, _ := strconv.Atoi(f[1])
		if tg != a[0] || idx != a[1] {
			return fmt.Errorf("expected %d %d got %d %d", a[0], a[1], tg, idx)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
