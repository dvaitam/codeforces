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
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "163C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	l := rng.Intn(100) + 1
	v1 := rng.Intn(20) + 1
	v2 := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(2 * l)
	}
	// sort arr ascending
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[i] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, l, v1, v2))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteString("\n")
	return sb.String()
}

func parseFloats(out string) ([]float64, error) {
	fields := strings.Fields(strings.TrimSpace(out))
	res := make([]float64, len(fields))
	for i, f := range fields {
		val, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, err
		}
		res[i] = val
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		expectedOut, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		expectedVals, err := parseFloats(expectedOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle parse error: %v\n", err)
			os.Exit(1)
		}
		gotOut, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		gotVals, err := parseFloats(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: parse output: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if len(gotVals) != len(expectedVals) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d values got %d\ninput:\n%s", i+1, len(expectedVals), len(gotVals), input)
			os.Exit(1)
		}
		for j := range gotVals {
			diff := gotVals[j] - expectedVals[j]
			if diff < 0 {
				diff = -diff
			}
			if diff > 1e-6 {
				fmt.Fprintf(os.Stderr, "case %d failed on value %d: expected %f got %f\ninput:\n%s", i+1, j, expectedVals[j], gotVals[j], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
