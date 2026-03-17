package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func floatEqual(a, b string) bool {
	fa, err1 := strconv.ParseFloat(strings.TrimSpace(a), 64)
	fb, err2 := strconv.ParseFloat(strings.TrimSpace(b), 64)
	if err1 != nil || err2 != nil {
		return strings.TrimSpace(a) == strings.TrimSpace(b)
	}
	if fb == 0 {
		return math.Abs(fa) < 1e-6
	}
	return math.Abs(fa-fb)/math.Max(1.0, math.Abs(fb)) < 1e-6
}

func buildOracle() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	content, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("read reference: %v", err)
	}
	oracle := filepath.Join(os.TempDir(), "oracleA")
	if strings.Contains(string(content), "#include") {
		cppSrc := filepath.Join(os.TempDir(), "oracleA.cpp")
		if err := os.WriteFile(cppSrc, content, 0644); err != nil {
			return "", err
		}
		cmd := exec.Command("g++", "-O2", "-o", oracle, cppSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
		}
	} else {
		cmd := exec.Command("go", "build", "-o", oracle, src)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
		}
	}
	return oracle, nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprint(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		y := r.Intn(21) - 10
		if y == 0 {
			y = 1
		}
		sb.WriteString(fmt.Sprintf("0 %d\n", y))
	}
	for i := 0; i < n; i++ {
		x := r.Intn(21) - 10
		if x == 0 {
			x = 1
		}
		sb.WriteString(fmt.Sprintf("%d 0\n", x))
	}
	return sb.String()
}

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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
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
		input := genCase(rng)
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if !floatEqual(got, expect) {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%sexpected: %s\ngot: %s\n", i, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
