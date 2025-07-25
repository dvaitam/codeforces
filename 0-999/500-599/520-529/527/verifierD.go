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

type point struct {
	x int64
	w int64
}

type testCase struct {
	points []point
}

func solve(tc testCase) string {
	intervals := make([][2]int64, len(tc.points))
	for i, p := range tc.points {
		intervals[i][0] = p.x - p.w
		intervals[i][1] = p.x + p.w
	}
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i][1] != intervals[j][1] {
			return intervals[i][1] < intervals[j][1]
		}
		return intervals[i][0] < intervals[j][0]
	})
	count := 0
	lastR := int64(-1 << 62)
	for _, iv := range intervals {
		if iv[0] >= lastR {
			count++
			lastR = iv[1]
		}
	}
	return fmt.Sprintf("%d", count)
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(1))
	tests := make([]testCase, 0, 100)
	for len(tests) < 100 {
		n := rnd.Intn(20) + 1
		used := map[int64]bool{}
		pts := make([]point, 0, n)
		for len(pts) < n {
			x := int64(rnd.Intn(200))
			if used[x] {
				continue
			}
			used[x] = true
			w := int64(rnd.Intn(50) + 1)
			pts = append(pts, point{x, w})
		}
		tests = append(tests, testCase{pts})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(tc.points))
		for _, p := range tc.points {
			fmt.Fprintf(&sb, "%d %d\n", p.x, p.w)
		}
		input := sb.String()
		expected := solve(tc)
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		want := strings.TrimSpace(expected)
		if got != want {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\n got: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
