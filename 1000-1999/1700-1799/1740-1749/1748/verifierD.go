package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		log.Fatal("REFERENCE_SOURCE_PATH environment variable is not set")
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(bin, input string) (string, error) {
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
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) (string, uint64, uint64, uint64) {
	a := uint64(r.Uint32() & ((1 << 30) - 1))
	if a == 0 {
		a = 1
	}
	b := uint64(r.Uint32() & ((1 << 30) - 1))
	if b == 0 {
		b = 1
	}
	d := uint64(r.Uint32() & ((1 << 30) - 1))
	if d == 0 {
		d = 1
	}
	return fmt.Sprintf("1\n%d %d %d\n", a, b, d), a, b, d
}

// validateAnswer checks whether x is a valid answer for the given a, b, d.
// If oracleAnswer is -1, the only valid answer is -1.
// Otherwise, x must satisfy: 0 <= x < 2^60, (a|x) % d == 0, (b|x) % d == 0.
func validateAnswer(x int64, a, b, d uint64, oracleAnswer string) error {
	if x == -1 {
		// Candidate says no solution. Check oracle agrees.
		if oracleAnswer != "-1" {
			return fmt.Errorf("candidate returned -1 but a solution exists (oracle: %s)", oracleAnswer)
		}
		return nil
	}
	if x < 0 {
		return fmt.Errorf("invalid answer %d (negative)", x)
	}
	ux := uint64(x)
	if ux >= (1 << 60) {
		return fmt.Errorf("x=%d is >= 2^60", x)
	}
	if (a|ux)%d != 0 {
		return fmt.Errorf("(a|x) = (%d|%d) = %d is not divisible by d=%d", a, ux, a|ux, d)
	}
	if (b|ux)%d != 0 {
		return fmt.Errorf("(b|x) = (%d|%d) = %d is not divisible by d=%d", b, ux, b|ux, d)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, a, b, d := genCase(rng)
		expect, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n%s", i, err, input)
			os.Exit(1)
		}
		x, parseErr := strconv.ParseInt(got, 10, 64)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to parse candidate output %q: %v\ninput:\n%s", i, got, parseErr, input)
			os.Exit(1)
		}
		if valErr := validateAnswer(x, a, b, d, expect); valErr != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, valErr, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
