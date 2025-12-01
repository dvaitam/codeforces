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

const refSourceB = "./2113B.go"

type testCase struct {
	w, h int64
	a, b int64
	x1   int64
	y1   int64
	x2   int64
	y2   int64
	name string
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "case %d (%s) mismatch: expected %s got %s\ninput:\n%s", i+1, tc.name, refAns[i], candAns[i], formatCase(tc))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	outPath := "./ref_2113B.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSourceB)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(name string, w, h, a, b, x1, y1, x2, y2 int64) {
		tests = append(tests, testCase{name: name, w: w, h: h, a: a, b: b, x1: x1, y1: y1, x2: x2, y2: y2})
	}

	// Hand-crafted examples
	add("example1", 6, 5, 2, 3, -1, -2, 5, 4)
	add("vertical_shift", 10, 9, 3, 2, 0, 0, 0, 6)
	add("horizontal_shift", 10, 9, 3, 2, 0, 0, 6, 0)
	add("edge_touch", 4, 4, 2, 2, 0, 0, 2, 2)
	add("far_apart", 5, 5, 2, 2, -1, -1, 4, -1)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 180 {
		w := int64(rng.Intn(50) + 1)
		h := int64(rng.Intn(50) + 1)
		if rng.Intn(5) == 0 {
			w = int64(rng.Intn(1_000_000_000) + 1)
			h = int64(rng.Intn(1_000_000_000) + 1)
		}
		a := int64(rng.Intn(20) + 1)
		b := int64(rng.Intn(20) + 1)
		if rng.Intn(4) == 0 {
			a = int64(rng.Intn(1000) + 1)
			b = int64(rng.Intn(1000) + 1)
		}
		x1, y1, x2, y2 := randomNonOverlap(w, h, a, b, rng)
		add(fmt.Sprintf("random_%d", len(tests)), w, h, a, b, x1, y1, x2, y2)
	}

	return tests
}

func randomNonOverlap(w, h, a, b int64, rng *rand.Rand) (int64, int64, int64, int64) {
	minX := -a + 1
	maxX := w - 1
	minY := -b + 1
	maxY := h - 1
	for attempt := 0; attempt < 10000; attempt++ {
		x1 := minX + int64(rng.Int63n(maxX-minX+1))
		y1 := minY + int64(rng.Int63n(maxY-minY+1))
		x2 := minX + int64(rng.Int63n(maxX-minX+1))
		y2 := minY + int64(rng.Int63n(maxY-minY+1))
		if nonOverlap(x1, y1, x2, y2, a, b) {
			return x1, y1, x2, y2
		}
		// Try separating deterministically
		if attempt%2 == 0 {
			x2 = x1 + a
			if x2 > maxX {
				x2 = x1 - a
			}
		} else {
			y2 = y1 + b
			if y2 > maxY {
				y2 = y1 - b
			}
		}
		if x2 < minX {
			x2 = minX
		}
		if x2 > maxX {
			x2 = maxX
		}
		if y2 < minY {
			y2 = minY
		}
		if y2 > maxY {
			y2 = maxY
		}
		if nonOverlap(x1, y1, x2, y2, a, b) {
			return x1, y1, x2, y2
		}
	}
	// Guaranteed non-overlap by placing far apart
	return minX, minY, maxX, maxY
}

func nonOverlap(x1, y1, x2, y2, a, b int64) bool {
	return x1+a <= x2 || x2+a <= x1 || y1+b <= y2 || y2+b <= y1
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.w, tc.h, tc.a, tc.b))
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.x1, tc.y1, tc.x2, tc.y2))
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(lines))
	}
	res := make([]string, expected)
	for i, s := range lines {
		up := strings.ToUpper(s)
		if up != "YES" && up != "NO" {
			return nil, fmt.Errorf("invalid token %q (case %d)", s, i+1)
		}
		res[i] = up
	}
	return res, nil
}

func formatCase(tc testCase) string {
	return fmt.Sprintf("%d %d %d %d\n%d %d %d %d\n", tc.w, tc.h, tc.a, tc.b, tc.x1, tc.y1, tc.x2, tc.y2)
}
