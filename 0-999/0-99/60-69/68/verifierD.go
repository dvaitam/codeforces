package main

import (
	"bytes"
	"fmt"
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
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "68D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func generateCase(rng *rand.Rand) string {
	h := rng.Intn(5) + 1 // 1..5
	q := rng.Intn(15) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", h, q))
	maxV := 1<<(h+1) - 1
	decayCount := 0
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			v := rng.Intn(maxV) + 1
			e := rng.Intn(101)
			sb.WriteString(fmt.Sprintf("add %d %d\n", v, e))
		} else {
			sb.WriteString("decay\n")
			decayCount++
		}
	}
	if decayCount == 0 {
		sb.WriteString("decay\n")
	}
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	expectedLines := strings.Fields(strings.TrimSpace(outO.String()))

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotLines := strings.Fields(strings.TrimSpace(out.String()))
	if len(gotLines) != len(expectedLines) {
		return fmt.Errorf("expected %d lines got %d", len(expectedLines), len(gotLines))
	}
	for i := range gotLines {
		var g, e float64
		fmt.Sscan(gotLines[i], &g)
		fmt.Sscan(expectedLines[i], &e)
		if diff := g - e; diff < -1e-4 || diff > 1e-4 {
			return fmt.Errorf("line %d: expected %.4f got %.4f", i+1, e, g)
		}
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
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		if err := runCase(bin, oracle, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
