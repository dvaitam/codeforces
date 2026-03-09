package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	oracle := filepath.Join(os.TempDir(), "oracleE2.bin")
	cmd := exec.Command("go", "build", "-o", oracle, "1744E2.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("oracle build failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin string, input string) (string, error) {
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

func generate() string {
	const T = 100
	rng := rand.New(rand.NewSource(6))
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", T)
	for i := 0; i < T; i++ {
		a := rng.Int63n(1_000_000) + 1
		b := rng.Int63n(1_000_000) + 1
		c := a + rng.Int63n(1_000_000) + 1
		d := b + rng.Int63n(1_000_000) + 1
		fmt.Fprintf(&sb, "%d %d %d %d\n", a, b, c, d)
	}
	return sb.String()
}

func parseLines(s string) []string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	result := make([]string, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			result = append(result, l)
		}
	}
	return result
}

func parsePair(s string) (int64, int64, error) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected 2 numbers, got %q", s)
	}
	x, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	y, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return x, y, nil
}

func check(input, oracleOut, candOut string) error {
	inLines := parseLines(input)
	oracleLines := parseLines(oracleOut)
	candLines := parseLines(candOut)

	// first line of input is t
	t, err := strconv.Atoi(inLines[0])
	if err != nil {
		return fmt.Errorf("bad input: %v", err)
	}
	if len(oracleLines) != t {
		return fmt.Errorf("oracle output has %d lines, expected %d", len(oracleLines), t)
	}
	if len(candLines) != t {
		return fmt.Errorf("candidate output has %d lines, expected %d", len(candLines), t)
	}

	for i := 0; i < t; i++ {
		parts := strings.Fields(inLines[i+1])
		a, _ := strconv.ParseInt(parts[0], 10, 64)
		b, _ := strconv.ParseInt(parts[1], 10, 64)
		c, _ := strconv.ParseInt(parts[2], 10, 64)
		d, _ := strconv.ParseInt(parts[3], 10, 64)

		ox, oy, err := parsePair(oracleLines[i])
		if err != nil {
			return fmt.Errorf("bad oracle output line %d: %v", i+1, err)
		}
		cx, cy, err := parsePair(candLines[i])
		if err != nil {
			return fmt.Errorf("bad candidate output line %d: %v", i+1, err)
		}

		oracleImpossible := ox == -1 && oy == -1

		if oracleImpossible {
			if cx != -1 || cy != -1 {
				return fmt.Errorf("case %d (a=%d b=%d c=%d d=%d): oracle says impossible but candidate output %d %d",
					i+1, a, b, c, d, cx, cy)
			}
			continue
		}

		// Solution exists; validate candidate's answer
		if cx == -1 && cy == -1 {
			return fmt.Errorf("case %d (a=%d b=%d c=%d d=%d): candidate says impossible but oracle found %d %d",
				i+1, a, b, c, d, ox, oy)
		}
		if cx <= a || cx > c {
			return fmt.Errorf("case %d (a=%d b=%d c=%d d=%d): x=%d not in (%d,%d]",
				i+1, a, b, c, d, cx, a, c)
		}
		if cy <= b || cy > d {
			return fmt.Errorf("case %d (a=%d b=%d c=%d d=%d): y=%d not in (%d,%d]",
				i+1, a, b, c, d, cy, b, d)
		}
		ab := a * b
		if (cx*cy)%ab != 0 {
			return fmt.Errorf("case %d (a=%d b=%d c=%d d=%d): x*y=%d not divisible by a*b=%d",
				i+1, a, b, c, d, cx*cy, ab)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	input := generate()
	exp, err := run(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
		os.Exit(1)
	}
	got, err := run(cand, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := check(input, exp, got); err != nil {
		fmt.Println("wrong answer:", err)
		fmt.Println("input:\n" + input)
		fmt.Println("expected:\n" + exp)
		fmt.Println("got:\n" + got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
