package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type stone struct {
	a, b, c int
}

type testCase struct {
	stones []stone
}

func expectedCase(tc testCase) string {
	type rect struct {
		a, b, c int
		id      int
	}
	n := len(tc.stones)
	rects := make([]rect, n)
	for i, s := range tc.stones {
		a, b, c := s.a, s.b, s.c
		if a < b {
			a, b = b, a
		}
		if a < c {
			a, c = c, a
		}
		if b < c {
			b, c = c, b
		}
		rects[i] = rect{a: a, b: b, c: c, id: i + 1}
	}
	best := 0
	p1 := 1
	for _, r := range rects {
		if r.c > best {
			best = r.c
			p1 = r.id
		}
	}
	sort.Slice(rects, func(i, j int) bool {
		if rects[i].a != rects[j].a {
			return rects[i].a > rects[j].a
		}
		if rects[i].b != rects[j].b {
			return rects[i].b > rects[j].b
		}
		return rects[i].c > rects[j].c
	})
	flag := false
	p2 := 0
	for i := 1; i < n; i++ {
		if rects[i].a == rects[i-1].a && rects[i].b == rects[i-1].b {
			sumC := rects[i].c + rects[i-1].c
			cand := sumC
			if cand > rects[i].b {
				cand = rects[i].b
			}
			if cand > best {
				best = cand
				p1 = rects[i].id
				p2 = rects[i-1].id
				flag = true
			}
		}
	}
	if !flag {
		return fmt.Sprintf("1\n%d\n", p1)
	}
	return fmt.Sprintf("2\n%d %d\n", p1, p2)
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.stones)))
	for _, s := range tc.stones {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", s.a, s.b, s.c))
	}
	return sb.String()
}

func normalizeOutput(s string) (int, []int, error) {
	lines := strings.Fields(strings.TrimSpace(s))
	if len(lines) < 2 {
		return 0, nil, fmt.Errorf("too few tokens in output")
	}
	k, err := strconv.Atoi(lines[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid k: %q", lines[0])
	}
	if k != 1 && k != 2 {
		return 0, nil, fmt.Errorf("k must be 1 or 2, got %d", k)
	}
	if len(lines) != k+1 {
		return 0, nil, fmt.Errorf("expected %d stone indices, got %d", k, len(lines)-1)
	}
	ids := make([]int, k)
	for i := 0; i < k; i++ {
		ids[i], err = strconv.Atoi(lines[i+1])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid index: %q", lines[i+1])
		}
	}
	sort.Ints(ids)
	return k, ids, nil
}

func runCase(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func check(got, exp string) error {
	expK, expIDs, err := normalizeOutput(exp)
	if err != nil {
		return fmt.Errorf("bad expected output: %v", err)
	}
	gotK, gotIDs, err := normalizeOutput(got)
	if err != nil {
		return fmt.Errorf("bad candidate output: %v", err)
	}
	if gotK != expK {
		return fmt.Errorf("expected k=%d got k=%d", expK, gotK)
	}
	for i := range expIDs {
		if gotIDs[i] != expIDs[i] {
			return fmt.Errorf("expected %q got %q", strings.TrimSpace(exp), strings.TrimSpace(got))
		}
	}
	return nil
}

func buildTests(rng *rand.Rand) []testCase {
	var tests []testCase
	add := func(ss ...stone) {
		tests = append(tests, testCase{stones: ss})
	}

	// Manual cases from the problem
	add(stone{5, 5, 5}, stone{3, 2, 4}, stone{1, 4, 1}, stone{2, 1, 3}, stone{3, 2, 4}, stone{3, 3, 4})
	add(stone{12, 18, 10}, stone{18, 4, 12})

	// Single stone
	add(stone{1, 1, 1})
	add(stone{1000000000, 1000000000, 1000000000})
	add(stone{3, 5, 7})

	// Two identical stones (can always glue)
	add(stone{4, 6, 8}, stone{4, 6, 8})
	add(stone{2, 3, 5}, stone{2, 3, 5})

	// Exhaustive small: n=1..3, edges from small set
	vals := []int{1, 2, 3, 4}
	for a := range vals {
		for b := range vals {
			for c := range vals {
				add(stone{vals[a], vals[b], vals[c]})
			}
		}
	}
	for n := 2; n <= 3; n++ {
		for trial := 0; trial < 80; trial++ {
			ss := make([]stone, n)
			for i := range ss {
				ss[i] = stone{vals[rng.Intn(len(vals))], vals[rng.Intn(len(vals))], vals[rng.Intn(len(vals))]}
			}
			tests = append(tests, testCase{stones: ss})
		}
	}

	// Random larger cases
	for trial := 0; trial < 200; trial++ {
		n := rng.Intn(10) + 1
		ss := make([]stone, n)
		for i := range ss {
			ss[i] = stone{rng.Intn(20) + 1, rng.Intn(20) + 1, rng.Intn(20) + 1}
		}
		tests = append(tests, testCase{stones: ss})
	}

	// Large value random
	for trial := 0; trial < 50; trial++ {
		n := rng.Intn(5) + 1
		ss := make([]stone, n)
		for i := range ss {
			ss[i] = stone{rng.Intn(1000000000) + 1, rng.Intn(1000000000) + 1, rng.Intn(1000000000) + 1}
		}
		tests = append(tests, testCase{stones: ss})
	}

	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := buildTests(rng)

	for i, tc := range tests {
		input := buildInput(tc)
		exp := expectedCase(tc)
		got, err := runCase(exe, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput: %s", i+1, err, input)
			os.Exit(1)
		}
		if err := check(got, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput: %s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
