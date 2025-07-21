package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	x int64
	t int64
}

func maxEventsStartZero(events []Event, V int64) int {
	n := len(events)
	sortEvents(events)
	size := 1 << n
	dp := make([][]bool, size)
	for i := range dp {
		dp[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		if abs64(events[i].x) <= V*events[i].t {
			dp[1<<i][i] = true
		}
	}
	best := 0
	for mask := 1; mask < size; mask++ {
		for last := 0; last < n; last++ {
			if !dp[mask][last] {
				continue
			}
			cnt := bitsCount(mask)
			if cnt > best {
				best = cnt
			}
			for nxt := last + 1; nxt < n; nxt++ {
				if mask&(1<<nxt) != 0 {
					continue
				}
				dt := events[nxt].t - events[last].t
				dist := abs64(events[nxt].x - events[last].x)
				if dist <= V*dt {
					dp[mask|1<<nxt][nxt] = true
				}
			}
		}
	}
	return best
}

func maxEventsFreeStart(events []Event, V int64) int {
	n := len(events)
	sortEvents(events)
	size := 1 << n
	dp := make([][]bool, size)
	for i := range dp {
		dp[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		dp[1<<i][i] = true
	}
	best := 0
	for mask := 1; mask < size; mask++ {
		for last := 0; last < n; last++ {
			if !dp[mask][last] {
				continue
			}
			cnt := bitsCount(mask)
			if cnt > best {
				best = cnt
			}
			for nxt := last + 1; nxt < n; nxt++ {
				if mask&(1<<nxt) != 0 {
					continue
				}
				dt := events[nxt].t - events[last].t
				dist := abs64(events[nxt].x - events[last].x)
				if dist <= V*dt {
					dp[mask|1<<nxt][nxt] = true
				}
			}
		}
	}
	return best
}

func bitsCount(x int) int {
	c := 0
	for x > 0 {
		x &= x - 1
		c++
	}
	return c
}

func sortEvents(ev []Event) {
	for i := 0; i < len(ev); i++ {
		for j := i + 1; j < len(ev); j++ {
			if ev[j].t < ev[i].t {
				ev[i], ev[j] = ev[j], ev[i]
			}
		}
	}
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solve(data string) string {
	fields := strings.Fields(data)
	if len(fields) == 0 {
		return ""
	}
	idx := 0
	n, _ := strconv.Atoi(fields[idx])
	idx++
	events := make([]Event, n)
	for i := 0; i < n; i++ {
		x, _ := strconv.ParseInt(fields[idx], 10, 64)
		idx++
		t, _ := strconv.ParseInt(fields[idx], 10, 64)
		idx++
		events[i] = Event{x: x, t: t}
	}
	V, _ := strconv.ParseInt(fields[idx], 10, 64)
	a1 := maxEventsStartZero(append([]Event(nil), events...), V)
	a2 := maxEventsFreeStart(events, V)
	return fmt.Sprintf("%d %d\n", a1, a2)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	events := make([]Event, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		x := int64(rng.Intn(41) - 20)
		t := int64(rng.Intn(10) + i + 1)
		sb.WriteString(fmt.Sprintf("%d %d\n", x, t))
		events[i] = Event{x: x, t: t}
	}
	V := int64(rng.Intn(5) + 1)
	sb.WriteString(fmt.Sprintf("%d\n", V))
	input := sb.String()
	expected := solve(input)
	return input, expected
}

func runCase(exe string, in, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
