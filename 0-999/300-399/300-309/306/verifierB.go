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

type interval struct {
	x, l int
}

type testCase struct {
	n    int
	segs []interval
}

func runCase(bin string, tc testCase) ([]int, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.segs)))
	for _, s := range tc.segs {
		sb.WriteString(fmt.Sprintf("%d %d\n", s.x, s.l))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return nil, fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid k: %v", err)
	}
	if k < 0 || k > len(tc.segs) {
		return nil, fmt.Errorf("invalid k %d", k)
	}
	if len(fields[1:]) != k {
		return nil, fmt.Errorf("expected %d indices, got %d", k, len(fields)-1)
	}
	rem := make([]int, k)
	for i, f := range fields[1:] {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid index %q", f)
		}
		if v < 1 || v > len(tc.segs) {
			return nil, fmt.Errorf("index out of range: %d", v)
		}
		rem[i] = v
	}
	return rem, nil
}

func union(n int, segs []interval) []bool {
	arr := make([]bool, n+1)
	for _, s := range segs {
		end := s.x + s.l - 1
		for i := s.x; i <= end && i <= n; i++ {
			arr[i] = true
		}
	}
	return arr
}

func checkRemoval(n int, segs []interval, rem []int) error {
	removed := make([]bool, len(segs)+1)
	for _, id := range rem {
		if removed[id] {
			return fmt.Errorf("duplicate index %d", id)
		}
		removed[id] = true
	}
	kept := make([]interval, 0, len(segs))
	for i, s := range segs {
		if !removed[i+1] {
			kept = append(kept, s)
		}
	}
	want := union(n, segs)
	got := union(n, kept)
	for i := 1; i <= n; i++ {
		if want[i] != got[i] {
			return fmt.Errorf("coverage mismatch at cell %d", i)
		}
	}
	return nil
}

func minimalKeep(segs []interval) int {
	type iv struct{ x, end int }
	arr := make([]iv, len(segs))
	for i, s := range segs {
		arr[i] = iv{s.x, s.x + s.l - 1}
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].x != arr[j].x {
			return arr[i].x < arr[j].x
		}
		return arr[i].end > arr[j].end
	})
	pre, now, keep := -1, -1, 0
	for _, iv := range arr {
		if iv.x <= pre+1 {
			if now < iv.end {
				now = iv.end
			}
		} else {
			if now > pre {
				pre = now
				keep++
				now = iv.end
			} else {
				pre = iv.end
				keep++
				now = -1
			}
		}
	}
	if now > pre {
		keep++
	}
	return keep
}

func generateCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 100)
	cases = append(cases, testCase{n: 1, segs: []interval{{1, 1}}})
	for len(cases) < 100 {
		n := rng.Intn(200) + 1
		m := rng.Intn(20) + 1
		if m > n {
			m = n
		}
		segs := make([]interval, m)
		for i := 0; i < m; i++ {
			x := rng.Intn(n) + 1
			maxL := n - x + 1
			l := rng.Intn(maxL) + 1
			segs[i] = interval{x, l}
		}
		cases = append(cases, testCase{n: n, segs: segs})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		rem, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkRemoval(tc.n, tc.segs, rem); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		keepMin := minimalKeep(tc.segs)
		if len(tc.segs)-len(rem) != keepMin {
			fmt.Fprintf(os.Stderr, "case %d failed: not maximal removal\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
