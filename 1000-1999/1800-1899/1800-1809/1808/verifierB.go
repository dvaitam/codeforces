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
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "1808B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(1000))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	run := func(exe string) (string, error) {
		var cmd *exec.Cmd
		if strings.HasSuffix(exe, ".go") {
			cmd = exec.Command("go", "run", exe)
		} else {
			cmd = exec.Command(exe)
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
	exp, err := run(oracle)
	if err != nil {
		return fmt.Errorf("oracle: %v", err)
	}
	got, err := run(bin)
	if err != nil {
		return err
	}
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
