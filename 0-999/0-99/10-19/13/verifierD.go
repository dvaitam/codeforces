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

type Point struct{ x, y int }

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

func cross(ax, ay, bx, by int) int {
	return ax*by - ay*bx
}

func inside(a, b, c, p Point) bool {
	s1 := cross(b.x-a.x, b.y-a.y, p.x-a.x, p.y-a.y)
	s2 := cross(c.x-b.x, c.y-b.y, p.x-b.x, p.y-b.y)
	s3 := cross(a.x-c.x, a.y-c.y, p.x-c.x, p.y-c.y)
	if (s1 > 0 && s2 > 0 && s3 > 0) || (s1 < 0 && s2 < 0 && s3 < 0) {
		return true
	}
	return false
}

func countTriangles(reds []Point, blues []Point) int {
	n := len(reds)
	cnt := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				if cross(reds[j].x-reds[i].x, reds[j].y-reds[i].y, reds[k].x-reds[i].x, reds[k].y-reds[i].y) == 0 {
					continue
				}
				ok := true
				for _, b := range blues {
					if inside(reds[i], reds[j], reds[k], b) {
						ok = false
						break
					}
				}
				if ok {
					cnt++
				}
			}
		}
	}
	return cnt
}

func generateCaseD(rng *rand.Rand) (string, int) {
	n := rng.Intn(4) + 3
	m := rng.Intn(4)
	reds := make([]Point, n)
	blues := make([]Point, m)
	for i := range reds {
		reds[i] = Point{rng.Intn(11) - 5, rng.Intn(11) - 5}
	}
	for i := range blues {
		blues[i] = Point{rng.Intn(11) - 5, rng.Intn(11) - 5}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, p := range reds {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	for _, p := range blues {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return sb.String(), countTriangles(reds, blues)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseD(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
