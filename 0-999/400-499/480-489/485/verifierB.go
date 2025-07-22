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

func expectedArea(points [][2]int64) int64 {
	minX, maxX := points[0][0], points[0][0]
	minY, maxY := points[0][1], points[0][1]
	for i := 1; i < len(points); i++ {
		x, y := points[i][0], points[i][1]
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	dx := maxX - minX
	dy := maxY - minY
	if dx < dy {
		dx = dy
	}
	return dx * dx
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(9) + 2 // 2..10
	pts := make([][2]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		x := rng.Int63n(2_000_000_001) - 1_000_000_000
		y := rng.Int63n(2_000_000_001) - 1_000_000_000
		pts[i] = [2]int64{x, y}
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	area := expectedArea(pts)
	return sb.String(), fmt.Sprintf("%d", area)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
