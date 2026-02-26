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
	"time"
)

func buildOracle() (string, error) {
	exe := "oracleC2.bin"
	cmd := exec.Command("go", "build", "-o", exe, "1354C2.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return "./" + exe, nil
}

func deterministicCases() []string {
	return []string{
		"1\n3\n",
		"1\n5\n",
		"1\n7\n",
	}
}

func randomCase(rng *rand.Rand) string {
	n := 2*(rng.Intn(100)) + 1
	if n < 3 {
		n = 3
	}
	return fmt.Sprintf("1\n%d\n", n)
}

func runCase(bin, oracle, input string) error {
	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	expected := strings.TrimSpace(outO.String())

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errb.String())
	}
	result := strings.TrimSpace(out.String())
	
	expFloat, errE := strconv.ParseFloat(expected, 64)
	resFloat, errR := strconv.ParseFloat(result, 64)
	
	if errE != nil || errR != nil {
		if result != expected {
			return fmt.Errorf("expected %s got %s", expected, result)
		}
	} else {
		diff := math.Abs(expFloat - resFloat)
		if diff > 1e-6 && diff/math.Max(1.0, math.Abs(expFloat)) > 1e-6 {
			return fmt.Errorf("expected %s got %s", expected, result)
		}
	}
	
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC2.go /path/to/binary")
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
	cases := deterministicCases()
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, input := range cases {
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
