package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type traveler struct{ x, d int }

type event struct {
	pos int
	typ int
}

func maxOverlap(trs []traveler, mask int) int {
	events := make([]event, 0, 2*len(trs))
	for i := 0; i < len(trs); i++ {
		if mask&(1<<i) == 0 {
			continue
		}
		s := trs[i].x
		e := trs[i].x + trs[i].d
		events = append(events, event{pos: e, typ: -1})
		events = append(events, event{pos: s, typ: 1})
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].pos != events[j].pos {
			return events[i].pos < events[j].pos
		}
		return events[i].typ < events[j].typ // leaves before enters
	})
	cur, maxV := 0, 0
	for _, ev := range events {
		cur += ev.typ
		if cur > maxV {
			maxV = cur
		}
	}
	return maxV
}

func brute(trs []traveler, a int) (int, [][]int) {
	n := len(trs)
	best := -1
	var bestSubs [][]int
	for mask := 0; mask < (1 << n); mask++ {
		if bits.OnesCount(uint(mask)) != a {
			continue
		}
		overlap := maxOverlap(trs, mask)
		if best == -1 || overlap < best {
			best = overlap
			bestSubs = bestSubs[:0]
			idx := make([]int, 0, a)
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					idx = append(idx, i+1)
				}
			}
			bestSubs = append(bestSubs, idx)
		} else if overlap == best {
			idx := make([]int, 0, a)
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					idx = append(idx, i+1)
				}
			}
			bestSubs = append(bestSubs, idx)
		}
	}
	return best, bestSubs
}

func generateCase(rng *rand.Rand) (string, string, []traveler) {
	n := rng.Intn(7) + 1
	a := rng.Intn(n) + 1
	trs := make([]traveler, n)
	for i := range trs {
		trs[i].x = rng.Intn(20) + 1
		trs[i].d = rng.Intn(10) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, a))
	for i, v := range trs {
		sb.WriteString(fmt.Sprintf("%d %d\n", v.x, v.d))
		_ = i
	}
	best, _ := brute(trs, a)
	exp := fmt.Sprintf("%d\n", best)
	return sb.String(), exp, trs
}

func runCase(exe string, input string, exp string, trs []traveler, a int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) < 1 {
		return fmt.Errorf("no output")
	}
	seat, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("bad seat: %v", err)
	}
	if seat != int(seat) {
		// not necessary but for style
	}
	if len(fields)-1 != a {
		return fmt.Errorf("expected %d indices, got %d", a, len(fields)-1)
	}
	idxs := make([]int, 0, a)
	seen := make(map[int]bool)
	for i := 1; i < len(fields); i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil {
			return fmt.Errorf("bad index: %v", err)
		}
		if v < 1 || v > len(trs) {
			return fmt.Errorf("index out of range")
		}
		if seen[v] {
			return fmt.Errorf("duplicate index")
		}
		seen[v] = true
		idxs = append(idxs, v)
	}
	// compute seat requirement for provided subset
	mask := 0
	for _, v := range idxs {
		mask |= 1 << (v - 1)
	}
	need := maxOverlap(trs, mask)

	best, _ := brute(trs, a)
	if seat != need || seat != best {
		return fmt.Errorf("expected seat %d but got %d (need %d)", best, seat, need)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp, trs := generateCase(rng)
		a, _ := strconv.Atoi(strings.Fields(in)[1])
		if err := runCase(exe, in, exp, trs, a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
