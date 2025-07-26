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

type interval struct{ a, b int64 }

type event struct {
	pos   int64
	delta int
}

func runProg(bin, input string) (string, error) {
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
		return out.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(intervals []interval) int64 {
	events := make([]event, 0, len(intervals)*2)
	for _, iv := range intervals {
		events = append(events, event{iv.a, 1})
		events = append(events, event{iv.b + 1, -1})
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].pos == events[j].pos {
			return events[i].delta > events[j].delta
		}
		return events[i].pos < events[j].pos
	})
	var cnt int64
	var prev int64
	if len(events) > 0 {
		prev = events[0].pos
	}
	for _, e := range events {
		if cnt == 1 && e.pos > prev {
			return prev
		}
		cnt += int64(e.delta)
		prev = e.pos
	}
	return -1
}

func genCase(rng *rand.Rand) []interval {
	n := rng.Intn(5) + 1
	res := make([]interval, n)
	for i := 0; i < n; i++ {
		a := rng.Int63n(50) + 1
		b := a + rng.Int63n(50)
		res[i] = interval{a, b}
	}
	return res
}

func runCase(bin string, ivs []interval) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(ivs)))
	for _, iv := range ivs {
		sb.WriteString(fmt.Sprintf("%d %d\n", iv.a, iv.b))
	}
	want := fmt.Sprintf("%d", solve(ivs))
	out, err := runProg(bin, sb.String())
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	if out != want {
		return fmt.Errorf("expected %s got %s\ninput:\n%s", want, out, sb.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		ivs := genCase(rng)
		if err := runCase(bin, ivs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
