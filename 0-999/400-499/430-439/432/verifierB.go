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

type testCase struct {
	n  int
	xs []int
	ys []int
}

func expected(tc testCase) [][2]int {
	count := make(map[int]int)
	for _, x := range tc.xs {
		count[x]++
	}
	res := make([][2]int, tc.n)
	for i := 0; i < tc.n; i++ {
		conflicts := count[tc.ys[i]]
		res[i][0] = (tc.n - 1) + conflicts
		res[i][1] = (tc.n - 1) - conflicts
	}
	return res
}

func runCase(exe string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.xs[i], tc.ys[i]))
	}
	input := sb.String()
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	r := strings.Fields(out.String())
	if len(r) != 2*tc.n {
		return fmt.Errorf("expected %d numbers, got %d", 2*tc.n, len(r))
	}
	res := expected(tc)
	for i := 0; i < tc.n; i++ {
		var gotHome, gotAway int
		fmt.Sscan(r[2*i], &gotHome)
		fmt.Sscan(r[2*i+1], &gotAway)
		if gotHome != res[i][0] || gotAway != res[i][1] {
			return fmt.Errorf("team %d expected %d %d got %d %d", i+1, res[i][0], res[i][1], gotHome, gotAway)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(50) + 2
		xs := make([]int, n)
		ys := make([]int, n)
		for j := 0; j < n; j++ {
			xs[j] = rng.Intn(100000) + 1
			for {
				ys[j] = rng.Intn(100000) + 1
				if ys[j] != xs[j] {
					break
				}
			}
		}
		tc := testCase{n: n, xs: xs, ys: ys}
		if err := runCase(exe, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
