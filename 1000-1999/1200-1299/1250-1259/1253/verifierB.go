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

func solveB(events []int) ([]int, bool) {
	inside := map[int]bool{}
	seen := map[int]bool{}
	res := []int{}
	last := 0
	for i, x := range events {
		if x > 0 {
			if seen[x] {
				return nil, false
			}
			seen[x] = true
			inside[x] = true
		} else {
			id := -x
			if !inside[id] {
				return nil, false
			}
			delete(inside, id)
		}
		if len(inside) == 0 {
			res = append(res, i-last+1)
			last = i + 1
			seen = map[int]bool{}
		}
	}
	if len(inside) != 0 || last != len(events) {
		return nil, false
	}
	return res, true
}

func generateCaseB(rng *rand.Rand) []int {
	n := rng.Intn(10) + 1
	events := make([]int, n)
	for i := 0; i < n; i++ {
		id := rng.Intn(5) + 1
		if rng.Intn(2) == 0 {
			events[i] = id
		} else {
			events[i] = -id
		}
	}
	return events
}

func checkOutputB(events []int, output string, possible bool) error {
	out := strings.TrimSpace(output)
	if !possible {
		if out != "-1" {
			return fmt.Errorf("expected -1 got %q", out)
		}
		return nil
	}
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		return fmt.Errorf("expected 2 lines got %d", len(lines))
	}
	var d int
	if _, err := fmt.Sscan(lines[0], &d); err != nil {
		return fmt.Errorf("failed to parse d: %v", err)
	}
	fields := strings.Fields(lines[1])
	if len(fields) != d {
		return fmt.Errorf("expected %d segment lengths got %d", d, len(fields))
	}
	seg := make([]int, d)
	for i := 0; i < d; i++ {
		if _, err := fmt.Sscan(fields[i], &seg[i]); err != nil {
			return fmt.Errorf("failed to parse segment length: %v", err)
		}
	}
	sum := 0
	idx := 0
	for _, c := range seg {
		sum += c
		inside := map[int]bool{}
		seen := map[int]bool{}
		for j := 0; j < c; j++ {
			x := events[idx]
			if x > 0 {
				if seen[x] {
					return fmt.Errorf("employee %d enters twice in a day", x)
				}
				seen[x] = true
				inside[x] = true
			} else {
				id := -x
				if !inside[id] {
					return fmt.Errorf("employee %d leaves before entering", id)
				}
				delete(inside, id)
			}
			idx++
		}
		if len(inside) != 0 {
			return fmt.Errorf("day not empty at end")
		}
	}
	if sum != len(events) {
		return fmt.Errorf("segments sum %d != %d", sum, len(events))
	}
	return nil
}

func runCaseB(bin string, events []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(events)))
	for i, v := range events {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	seg, ok := solveB(events)
	possible := ok && seg != nil
	if err := checkOutputB(events, out.String(), possible); err != nil {
		return err
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
		events := generateCaseB(rng)
		if err := runCaseB(bin, events); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
