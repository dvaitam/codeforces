package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Update struct {
	t int64
	h int64
}

type Enemy struct {
	maxH  int64
	start int64
	regen int64
	ups   []Update
}

type Event struct {
	t     int64
	delta int64
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func addInterval(l, r int64, events *[]Event) {
	if l > r {
		return
	}
	*events = append(*events, Event{t: l, delta: 1}, Event{t: r + 1, delta: -1})
}

func processSegment(start, end, h, regen, maxH, damage int64, events *[]Event) {
	if start >= end || h > damage {
		return
	}
	if regen == 0 || maxH <= damage {
		addInterval(start, end-1, events)
		return
	}
	if damage < h {
		return
	}
	limit := (damage - h) / regen
	tEnd := start + limit
	if tEnd >= end {
		tEnd = end - 1
	}
	if tEnd >= start {
		addInterval(start, tEnd, events)
	}
}

func solveCase(n, m int, bounty, increase, damage int64, enemies []Enemy, updates [][3]int64) string {
	for i := range enemies {
		enemies[i].ups = nil
	}
	for _, u := range updates {
		t, id, h := u[0], u[1], u[2]
		enemies[id].ups = append(enemies[id].ups, Update{t: t, h: h})
	}
	const INF int64 = 1 << 60
	var events []Event
	infinite := false
	for i := 0; i < n; i++ {
		e := &enemies[i]
		sort.Slice(e.ups, func(a, b int) bool { return e.ups[a].t < e.ups[b].t })
		prevT := int64(0)
		prevH := e.start
		for _, u := range e.ups {
			processSegment(prevT, u.t, prevH, e.regen, e.maxH, damage, &events)
			prevT = u.t
			prevH = u.h
		}
		processSegment(prevT, INF, prevH, e.regen, e.maxH, damage, &events)
		if increase > 0 {
			finalH := prevH
			if e.regen > 0 {
				finalH = e.maxH
			}
			if finalH <= damage {
				infinite = true
			}
		}
	}
	if infinite {
		return "-1"
	}
	sort.Slice(events, func(i, j int) bool { return events[i].t < events[j].t })
	var curr int64
	var ans int64
	prev := int64(0)
	for i := 0; i < len(events); {
		t := events[i].t
		if curr > 0 && prev <= t-1 {
			gold := curr * (bounty + increase*(t-1))
			if gold > ans {
				ans = gold
			}
		}
		for i < len(events) && events[i].t == t {
			curr += events[i].delta
			i++
		}
		prev = t
	}
	return fmt.Sprint(ans)
}

func generateCase(rng *rand.Rand) (int, int, int64, int64, int64, []Enemy, [][3]int64) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3)
	bounty := rng.Int63n(20) + 1
	increase := rng.Int63n(5)
	damage := rng.Int63n(20) + 1
	enemies := make([]Enemy, n)
	for i := 0; i < n; i++ {
		maxH := rng.Int63n(20) + 1
		start := rng.Int63n(maxH) + 1
		regen := rng.Int63n(maxH + 1)
		enemies[i] = Enemy{maxH: maxH, start: start, regen: regen}
	}
	updates := make([][3]int64, m)
	for i := 0; i < m; i++ {
		t := rng.Int63n(10) + 1
		id := int64(rng.Intn(n))
		h := rng.Int63n(enemies[id].maxH) + 1
		updates[i] = [3]int64{t, id, h}
	}
	return n, m, bounty, increase, damage, enemies, updates
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, bounty, increase, damage, enemies, updates := generateCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		sb.WriteString(fmt.Sprintf("%d %d %d\n", bounty, increase, damage))
		for j := 0; j < n; j++ {
			e := enemies[j]
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e.maxH, e.start, e.regen))
		}
		for j := 0; j < m; j++ {
			u := updates[j]
			sb.WriteString(fmt.Sprintf("%d %d %d\n", u[0], u[1]+1, u[2]))
		}
		input := sb.String()
		expect := solveCase(n, m, bounty, increase, damage, append([]Enemy(nil), enemies...), updates)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
