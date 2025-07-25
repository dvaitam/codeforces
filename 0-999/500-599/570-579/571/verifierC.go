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
	exe := filepath.Join(os.TempDir(), "oracleC")
	cmd := exec.Command("go", "build", "-o", exe, "571C.go")
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
	nc := rng.Intn(3) + 1 // number of clauses
	n := rng.Intn(4) + 1  // number of variables
	// prepare slices
	clauses := make([][]int, nc)
	for i := 1; i <= n; i++ {
		occ := rng.Intn(2) + 1 // 1 or 2 occurrences
		cats := rng.Perm(nc)[:occ]
		for _, c := range cats {
			val := i
			if rng.Intn(2) == 0 {
				val = -val
			}
			clauses[c] = append(clauses[c], val)
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", nc, n)
	for _, cl := range clauses {
		fmt.Fprintf(&sb, "%d", len(cl))
		for _, v := range cl {
			fmt.Fprintf(&sb, " %d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
