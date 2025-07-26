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

type point struct{ x, y int }

type testCase struct {
	input    string
	expected string
}

func inside(n, d int, p point) string {
	if p.x+p.y >= d && p.x+p.y <= 2*n-d && p.y-p.x >= -d && p.y-p.x <= d {
		return "YES\n"
	}
	return "NO\n"
}

func buildCase(n, d int, pts []point) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, d))
	sb.WriteString(fmt.Sprintf("%d\n", len(pts)))
	var expSB strings.Builder
	for _, pt := range pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", pt.x, pt.y))
		expSB.WriteString(inside(n, d, pt))
	}
	return testCase{input: sb.String(), expected: expSB.String()}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 1
	d := rng.Intn(n + 1)
	m := rng.Intn(20) + 1
	pts := make([]point, m)
	for i := range pts {
		pts[i] = point{rng.Intn(n + 1), rng.Intn(n + 1)}
	}
	return buildCase(n, d, pts)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if out.String() != tc.expected {
		return fmt.Errorf("expected %q got %q", tc.expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		buildCase(7, 2, []point{{2, 4}, {4, 5}}),
		buildCase(1, 0, []point{{0, 0}}),
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
