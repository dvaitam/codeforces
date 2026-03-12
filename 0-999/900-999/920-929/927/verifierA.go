package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// generateTestInput creates a valid test input for 927A.
// Format: w h\nk\nx1 y1\n...\norders (t sx sy tx ty)\n-1 0 0 0 0\n
func generateTestInput(rng *rand.Rand) (string, int) {
	w := 300 + rng.Intn(100)
	h := 300 + rng.Intn(100)
	k := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", w, h)
	fmt.Fprintf(&sb, "%d\n", k)
	for i := 0; i < k; i++ {
		x := rng.Intn(w) + 1
		y := rng.Intn(h) + 1
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	q := rng.Intn(5) + 1
	t := 0
	for i := 0; i < q; i++ {
		t += rng.Intn(100) + 1
		if t > 86400 {
			q = i
			break
		}
		sx := rng.Intn(w) + 1
		sy := rng.Intn(h) + 1
		tx := rng.Intn(w) + 1
		ty := rng.Intn(h) + 1
		for tx == sx && ty == sy {
			tx = rng.Intn(w) + 1
			ty = rng.Intn(h) + 1
		}
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", t, sx, sy, tx, ty)
	}
	fmt.Fprintf(&sb, "-1 0 0 0 0\n")
	return sb.String(), q
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func main() {
	// Handle signals gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		os.Exit(1)
	}()

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	numTests := 20
	for i := 0; i < numTests; i++ {
		input, q := generateTestInput(rng)
		candOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		// This is an optimization/interactive problem.
		// The candidate output should have exactly q+2 response lines.
		// Each response line should start with a non-negative integer f (number of cars to instruct).
		candLines := strings.Split(candOut, "\n")
		expectedLines := q + 2
		if len(candLines) != expectedLines {
			fmt.Fprintf(os.Stderr, "test %d: expected %d output lines, got %d\ninput:\n%s\ncand output:\n%s\n",
				i+1, expectedLines, len(candLines), input, candOut)
			os.Exit(1)
		}
		for j, line := range candLines {
			line = strings.TrimSpace(line)
			if line == "" {
				fmt.Fprintf(os.Stderr, "test %d line %d: empty line\ninput:\n%s\n", i+1, j+1, input)
				os.Exit(1)
			}
			fields := strings.Fields(line)
			f, err := strconv.Atoi(fields[0])
			if err != nil || f < 0 {
				fmt.Fprintf(os.Stderr, "test %d line %d: invalid car count %q\ninput:\n%s\n", i+1, j+1, fields[0], input)
				os.Exit(1)
			}
			_ = f
		}
	}
	fmt.Printf("All %d tests passed\n", numTests)
}
