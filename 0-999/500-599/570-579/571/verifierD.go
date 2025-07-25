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
	exe := filepath.Join(os.TempDir(), "oracleD")
	cmd := exec.Command("go", "build", "-o", exe, "571D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return exe, nil
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	hasQ := false
	ops := []string{"U", "M", "A", "Z", "Q"}
	for i := 0; i < m; i++ {
		op := ops[rng.Intn(len(ops))]
		switch op {
		case "U", "M":
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "%s %d %d\n", op, a, b)
		case "A":
			x := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "%s %d\n", op, x)
		case "Z":
			y := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "%s %d\n", op, y)
		case "Q":
			hasQ = true
			q := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "%s %d\n", op, q)
		}
	}
	if !hasQ {
		q := rng.Intn(n) + 1
		fmt.Fprintf(&sb, "Q %d\n", q)
	}
	return sb.String()
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
		input := genCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
