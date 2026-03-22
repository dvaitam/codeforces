package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Embedded correct solver for 1379D.
func solve(input string) string {
	r := strings.NewReader(input)
	var n, h, m, k int
	fmt.Fscan(r, &n, &h, &m, &k)

	H := m / 2
	x := make([]int, n)

	type Event struct {
		t   int
		val int
	}
	events := make([]Event, 0, 2*n+1)
	events = append(events, Event{0, 0})

	for i := 0; i < n; i++ {
		var hi, mi int
		fmt.Fscan(r, &hi, &mi)
		xi := mi % H
		x[i] = xi

		if k == 1 {
			continue
		}

		if xi+1 < H {
			if xi+k <= H {
				events = append(events, Event{xi + 1, 1})
				events = append(events, Event{xi + k, -1})
			} else {
				events = append(events, Event{xi + 1, 1})
				events = append(events, Event{0, 1})
				events = append(events, Event{xi + k - H, -1})
			}
		} else {
			events = append(events, Event{0, 1})
			events = append(events, Event{k - 1, -1})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].t < events[j].t
	})

	var merged []Event
	for _, e := range events {
		if e.t >= H {
			continue
		}
		if len(merged) > 0 && merged[len(merged)-1].t == e.t {
			merged[len(merged)-1].val += e.val
		} else {
			merged = append(merged, e)
		}
	}

	minCanceled := n + 1
	optT := 0
	curr := 0

	for _, e := range merged {
		curr += e.val
		if curr < minCanceled {
			minCanceled = curr
			optT = e.t
		}
	}

	if minCanceled == n+1 {
		minCanceled = 0
	}

	var sb strings.Builder
	fmt.Fprintln(&sb, minCanceled, optT)

	if minCanceled > 0 {
		var canceled []int
		for i := 0; i < n; i++ {
			xi := x[i]
			if k == 1 {
				continue
			}
			start := (xi + 1) % H
			end := (xi + k - 1) % H
			isCanceled := false
			if start <= end {
				isCanceled = (optT >= start && optT <= end)
			} else {
				isCanceled = (optT >= start || optT <= end)
			}
			if isCanceled {
				canceled = append(canceled, i+1)
			}
		}
		for i, idx := range canceled {
			if i > 0 {
				fmt.Fprint(&sb, " ")
			}
			fmt.Fprint(&sb, idx)
		}
		fmt.Fprintln(&sb)
	} else {
		fmt.Fprintln(&sb)
	}

	return strings.TrimSpace(sb.String())
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

type Case struct{ input string }

func genCases() []Case {
	rng := rand.New(rand.NewSource(4))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(5) + 1
		h := rng.Intn(10) + 1
		m := (rng.Intn(10) + 1) * 2
		k := rng.Intn(m/2) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d\n", n, h, m, k)
		for j := 0; j < n; j++ {
			hi := rng.Intn(h)
			mi := rng.Intn(m)
			fmt.Fprintf(&sb, "%d %d\n", hi, mi)
		}
		cases[i] = Case{sb.String()}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	cases := genCases()
	for i, c := range cases {
		expect := solve(c.input)
		got, err := runBinary(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
