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

type cut struct {
	t byte
	x int
}

type testCase struct {
	w, h int
	cuts []cut
}

func solve(tc testCase) string {
	v := []int{0, tc.w}
	h := []int{0, tc.h}
	out := make([]int64, len(tc.cuts))
	for i, c := range tc.cuts {
		if c.t == 'V' {
			v = append(v, c.x)
			sort.Ints(v)
		} else {
			h = append(h, c.x)
			sort.Ints(h)
		}
		maxW := 0
		for j := 1; j < len(v); j++ {
			if v[j]-v[j-1] > maxW {
				maxW = v[j] - v[j-1]
			}
		}
		maxH := 0
		for j := 1; j < len(h); j++ {
			if h[j]-h[j-1] > maxH {
				maxH = h[j] - h[j-1]
			}
		}
		out[i] = int64(maxW) * int64(maxH)
	}
	var buf strings.Builder
	for _, v := range out {
		fmt.Fprintln(&buf, v)
	}
	return strings.TrimSpace(buf.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(1))
	tests := make([]testCase, 0, 100)
	for len(tests) < 100 {
		w := rnd.Intn(50) + 2
		h := rnd.Intn(50) + 2
		n := rnd.Intn(10) + 1
		cuts := make([]cut, 0, n)
		vset := map[int]bool{0: true, w: true}
		hset := map[int]bool{0: true, h: true}
		for len(cuts) < n {
			if rnd.Intn(2) == 0 {
				x := rnd.Intn(w-1) + 1
				if !vset[x] {
					vset[x] = true
					cuts = append(cuts, cut{'V', x})
				}
			} else {
				y := rnd.Intn(h-1) + 1
				if !hset[y] {
					hset[y] = true
					cuts = append(cuts, cut{'H', y})
				}
			}
		}
		tests = append(tests, testCase{w, h, cuts})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", tc.w, tc.h, len(tc.cuts))
		for _, c := range tc.cuts {
			fmt.Fprintf(&sb, "%c %d\n", c.t, c.x)
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
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\n got:\n%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
