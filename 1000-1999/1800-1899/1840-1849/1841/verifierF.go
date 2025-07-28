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

type point struct {
	x, y int64
}

func quadrant(p point) int {
	if p.x > 0 && p.y >= 0 {
		return 1
	} else if p.x <= 0 && p.y > 0 {
		return 2
	} else if p.x < 0 && p.y <= 0 {
		return 3
	}
	return 4
}

type testCaseF struct {
	input    string
	expected string
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveF(groups [][4]int64) float64 {
	var v []point
	var sx, sy int64
	for _, g := range groups {
		a, b, c, d := g[0], g[1], g[2], g[3]
		dx := a - b
		dy := c - d
		if dx == 0 && dy == 0 {
			continue
		}
		v = append(v, point{dx, dy})
		v = append(v, point{-dx, -dy})
		if c < d || (c == d && a < b) {
			sx += dx
			sy += dy
		}
	}
	sort.Slice(v, func(i, j int) bool {
		pi, pj := v[i], v[j]
		qi, qj := quadrant(pi), quadrant(pj)
		if qi != qj {
			return qi < qj
		}
		return pi.x*pj.y > pi.y*pj.x
	})
	ans := float64(sx*sx + sy*sy)
	for _, p := range v {
		sx += p.x
		sy += p.y
		val := float64(sx*sx + sy*sy)
		if val > ans {
			ans = val
		}
	}
	return ans
}

func generateCaseF(rng *rand.Rand) testCaseF {
	n := rng.Intn(6) + 1
	var in strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", n))
	groups := make([][4]int64, n)
	for i := 0; i < n; i++ {
		for j := 0; j < 4; j++ {
			groups[i][j] = int64(rng.Intn(5))
		}
		in.WriteString(fmt.Sprintf("%d %d %d %d\n", groups[i][0], groups[i][1], groups[i][2], groups[i][3]))
	}
	exp := fmt.Sprintf("%.10f\n", solveF(groups))
	return testCaseF{input: in.String(), expected: exp}
}

func runCaseF(bin string, tc testCaseF) error {
	got, err := runCandidate(bin, tc.input)
	if err != nil {
		return err
	}
	got = strings.TrimSpace(got)
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCaseF{generateCaseF(rng)}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseF(rng))
	}
	for i, tc := range cases {
		if err := runCaseF(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
