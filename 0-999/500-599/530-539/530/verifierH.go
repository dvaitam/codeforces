package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
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

func expected(xs, ys []int) float64 {
	n := len(xs)
	maxX, maxY := 0, 0
	for i := 0; i < n; i++ {
		if xs[i] > maxX {
			maxX = xs[i]
		}
		if ys[i] > maxY {
			maxY = ys[i]
		}
	}
	best := math.Inf(1)
	for A := maxX + 1; ; A++ {
		B := 0
		for i := 0; i < n; i++ {
			num := A * ys[i]
			den := A - xs[i]
			b := num / den
			if num%den != 0 {
				b++
			}
			if b > B {
				B = b
			}
		}
		area := float64(A*B) / 2.0
		if area < best {
			best = area
		}
		if float64(A*maxY)/2.0 > best {
			break
		}
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		xs := make([]int, n)
		ys := make([]int, n)
		for j := 0; j < n; j++ {
			xs[j] = rng.Intn(5) + 1
			ys[j] = rng.Intn(5) + 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			fmt.Fprintf(&sb, "%d %d\n", xs[j], ys[j])
		}
		input := sb.String()
		want := expected(xs, ys)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseFloat(strings.TrimSpace(out), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output\n", i+1)
			os.Exit(1)
		}
		if math.Abs(got-want) > 1e-4*math.Max(1, math.Abs(want)) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %s\n", i+1, want, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
