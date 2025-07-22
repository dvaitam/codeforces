package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct{ n int }

func runCase(bin string, n int) ([]string, error) {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	return lines, nil
}

func parseCoords(lines []string, n int) ([][2]float64, error) {
	if len(lines) != n {
		return nil, fmt.Errorf("expected %d lines, got %d", n, len(lines))
	}
	pts := make([][2]float64, n)
	for i, line := range lines {
		var x, y float64
		if _, err := fmt.Sscan(line, &x, &y); err != nil {
			return nil, fmt.Errorf("parse line %d: %v", i+1, err)
		}
		if math.Abs(x) > 1e6 || math.Abs(y) > 1e6 {
			return nil, fmt.Errorf("coordinate out of bounds")
		}
		pts[i] = [2]float64{x, y}
	}
	return pts, nil
}

func checkPolygon(n int, pts [][2]float64) error {
	edges := make([][2]float64, n)
	lens := make([]float64, n)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		dx := pts[j][0] - pts[i][0]
		dy := pts[j][1] - pts[i][1]
		edges[i] = [2]float64{dx, dy}
		lens[i] = math.Hypot(dx, dy)
		if lens[i] < 1-1e-3 || lens[i] > 1000+1e-3 {
			return fmt.Errorf("side length out of range")
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if math.Abs(lens[i]-lens[j]) < 1e-3 {
				return fmt.Errorf("side lengths must differ")
			}
		}
	}
	expectCos := math.Cos(2 * math.Pi / float64(n))
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		k := (i + 2) % n
		v1 := edges[j]
		v2 := edges[k]
		dot := v1[0]*v2[0] + v1[1]*v2[1]
		cos := dot / (lens[j] * lens[k])
		if math.Abs(cos-expectCos) > 1e-3 {
			return fmt.Errorf("angles not equal")
		}
		cross := v1[0]*v2[1] - v1[1]*v2[0]
		if cross <= 0 {
			return fmt.Errorf("polygon not convex ccw")
		}
	}
	return nil
}

func generateCases() []testCase {
	cases := []testCase{{3}, {4}}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 100 {
		n := rng.Intn(98) + 3
		cases = append(cases, testCase{n})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		lines, err := runCase(bin, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if tc.n < 5 {
			if len(lines) != 1 || strings.TrimSpace(lines[0]) != "No solution" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected 'No solution'\n", i+1)
				os.Exit(1)
			}
			continue
		}
		pts, err := parseCoords(lines, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkPolygon(tc.n, pts); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
