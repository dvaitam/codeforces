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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(pts [][2]int) float64 {
	maxd2 := 0.0
	for i := 0; i < len(pts); i++ {
		for j := i + 1; j < len(pts); j++ {
			dx := float64(pts[i][0] - pts[j][0])
			dy := float64(pts[i][1] - pts[j][1])
			d2 := dx*dx + dy*dy
			if d2 > maxd2 {
				maxd2 = d2
			}
		}
	}
	return math.Sqrt(maxd2)
}

func runCase(bin string, pts [][2]int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", len(pts)))
	for _, p := range pts {
		input.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	out, err := run(bin, input.String())
	if err != nil {
		return err
	}
	var got float64
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(pts)
	if math.Abs(got-exp) > 1e-4 {
		return fmt.Errorf("expected %.5f got %.5f", exp, got)
	}
	return nil
}

func randomPoints(rng *rand.Rand, n int) [][2]int {
	pts := make([][2]int, n)
	for i := range pts {
		pts[i][0] = rng.Intn(101) - 50
		pts[i][1] = rng.Intn(101) - 50
	}
	return pts
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	edge := [][][2]int{
		{{0, 0}, {0, 0}},
		{{0, 0}, {3, 4}},
		{{-50, -50}, {50, 50}},
		{{-50, 50}, {50, -50}, {0, 0}},
	}
	idx := 0
	for ; idx < len(edge); idx++ {
		if err := runCase(bin, edge[idx]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		n := rng.Intn(49) + 2
		pts := randomPoints(rng, n)
		if err := runCase(bin, pts); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
