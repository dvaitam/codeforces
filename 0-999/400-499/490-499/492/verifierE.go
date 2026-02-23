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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type TestCase struct {
	Input string
	N     int64
	Dx    int64
	Dy    int64
	Pts   [][2]int64
}

func generateCase(rng *rand.Rand) TestCase {
	n := int64(rng.Intn(20) + 2)
	var dx, dy int64
	for {
		dx = int64(rng.Intn(int(n-1)) + 1)
		if gcd(n, dx) == 1 {
			break
		}
	}
	for {
		dy = int64(rng.Intn(int(n-1)) + 1)
		if gcd(n, dy) == 1 {
			break
		}
	}
	m := int64(rng.Intn(20) + 1)
	pts := make([][2]int64, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, m, dx, dy)
	for i := int64(0); i < m; i++ {
		x := int64(rng.Intn(int(n)))
		y := int64(rng.Intn(int(n)))
		pts[i] = [2]int64{x, y}
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	return TestCase{
		Input: sb.String(),
		N:     n,
		Dx:    dx,
		Dy:    dy,
		Pts:   pts,
	}
}

func getMaxApples(n, dx, dy int64, pts [][2]int64) int {
	counts := make(map[int64]int)
	maxCount := 0
	for _, p := range pts {
		inv := (p[0]*dy - p[1]*dx) % n
		if inv < 0 {
			inv += n
		}
		counts[inv]++
		if counts[inv] > maxCount {
			maxCount = counts[inv]
		}
	}
	return maxCount
}

func getApplesForStart(n, dx, dy, startX, startY int64, pts [][2]int64) int {
	counts := 0
	targetInv := (startX*dy - startY*dx) % n
	if targetInv < 0 {
		targetInv += n
	}
	for _, p := range pts {
		inv := (p[0]*dy - p[1]*dx) % n
		if inv < 0 {
			inv += n
		}
		if inv == targetInv {
			counts++
		}
	}
	return counts
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		out, err := runCandidate(bin, tc.Input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.Input)
			os.Exit(1)
		}
		
		var gotX, gotY int64
		_, err = fmt.Sscanf(out, "%d %d", &gotX, &gotY)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: could not parse output %q\ninput:\n%s", i+1, out, tc.Input)
			os.Exit(1)
		}

		maxPossible := getMaxApples(tc.N, tc.Dx, tc.Dy, tc.Pts)
		gotApples := getApplesForStart(tc.N, tc.Dx, tc.Dy, gotX, gotY, tc.Pts)
		
		if gotApples != maxPossible {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d apples, got %d (started at %d %d)\ninput:\n%s", i+1, maxPossible, gotApples, gotX, gotY, tc.Input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
