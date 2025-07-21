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

type interval struct {
	l, r int
}

type testCaseC struct {
	intervals []interval
}

func generateCaseC(rng *rand.Rand) testCaseC {
	n := rng.Intn(8) + 1
	intervals := make([]interval, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(90) + 1
		r := l + rng.Intn(10) + 1
		intervals[i] = interval{l: l, r: r}
	}
	return testCaseC{intervals: intervals}
}

func expectedC(iv []interval) []int {
	n := len(iv)
	res := []int{}
	for i := 0; i < n; i++ {
		list := make([]interval, 0, n-1)
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			list = append(list, iv[j])
		}
		sort.Slice(list, func(a, b int) bool {
			if list[a].l == list[b].l {
				return list[a].r < list[b].r
			}
			return list[a].l < list[b].l
		})
		ok := true
		for j := 1; j < len(list); j++ {
			if list[j].l < list[j-1].r {
				ok = false
				break
			}
		}
		if ok {
			res = append(res, i+1)
		}
	}
	return res
}

func runCaseC(bin string, tc testCaseC) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.intervals)))
	for _, iv := range tc.intervals {
		sb.WriteString(fmt.Sprintf("%d %d\n", iv.l, iv.r))
	}
	input := sb.String()

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := expectedC(tc.intervals)
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	var k int
	if _, err := fmt.Sscan(fields[0], &k); err != nil {
		return fmt.Errorf("failed to parse k: %v", err)
	}
	if k != len(expected) {
		return fmt.Errorf("expected k=%d got %d", len(expected), k)
	}
	if k == 0 {
		return nil
	}
	if len(fields[1:]) != k {
		return fmt.Errorf("expected %d indices got %d", k, len(fields[1:]))
	}
	got := make([]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Sscan(fields[i+1], &got[i]); err != nil {
			return fmt.Errorf("bad index: %v", err)
		}
	}
	sort.Ints(got)
	sort.Ints(expected)
	for i := range got {
		if got[i] != expected[i] {
			return fmt.Errorf("expected %v got %v", expected, got)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseC(rng)
		if err := runCaseC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
