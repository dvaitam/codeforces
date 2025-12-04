package main

import (
    "bytes"
    "fmt"
    "math"
    "math/rand"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "954E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	T := rng.Intn(50) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, T)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(50)+1)
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(50)+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 1; i <= cases; i++ {
		input := genCase(rng)
        expectStr, err := run(oracle, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
            os.Exit(1)
        }
        gotStr, err := run(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
            os.Exit(1)
        }
        var expect, got float64
        if _, err := fmt.Sscan(strings.TrimSpace(expectStr), &expect); err != nil {
            fmt.Fprintf(os.Stderr, "oracle produced invalid output on case %d: %q\n", i, expectStr)
            os.Exit(1)
        }
        if _, err := fmt.Sscan(strings.TrimSpace(gotStr), &got); err != nil {
            fmt.Fprintf(os.Stderr, "candidate produced invalid output on case %d: %q\n", i, gotStr)
            os.Exit(1)
        }
        diff := math.Abs(got - expect)
        tol := 1e-6 * math.Max(1.0, math.Abs(expect))
        if diff > tol {
            fmt.Printf("case %d failed\nexpected:\n%.10f\n\ngot:\n%.10f\n", i, expect, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", cases)
}
