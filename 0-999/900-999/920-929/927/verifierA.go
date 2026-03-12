package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
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

func locateReference() (string, error) {
	if p := os.Getenv("REFERENCE_SOURCE_PATH"); p != "" {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	candidates := []string{
		"927A.go",
		filepath.Join("0-999", "900-999", "920-929", "927", "927A.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 927A.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref927A_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
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

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	numTests := 20
	for i := 0; i < numTests; i++ {
		input, _ := generateTestInput(rng)
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		candOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		// Both outputs should have the same number of lines (q+2 lines of "0")
		refLines := strings.Split(refOut, "\n")
		candLines := strings.Split(candOut, "\n")
		if len(refLines) != len(candLines) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d output lines, got %d\ninput:\n%s\nref output:\n%s\ncand output:\n%s\n",
				i+1, len(refLines), len(candLines), input, refOut, candOut)
			os.Exit(1)
		}
		for j := range refLines {
			if strings.TrimSpace(refLines[j]) != strings.TrimSpace(candLines[j]) {
				fmt.Fprintf(os.Stderr, "test %d line %d: expected %q, got %q\ninput:\n%s\n",
					i+1, j+1, refLines[j], candLines[j], input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", numTests)
}
