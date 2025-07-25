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
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "582E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

var vars = []byte{'A', 'B', 'C', 'D', 'a', 'b', 'c', 'd', '?'}
var ops = []byte{'&', '|', '?'}

func genExpr(rng *rand.Rand, depth int) string {
	if depth == 0 || rng.Intn(3) == 0 {
		return string(vars[rng.Intn(len(vars))])
	}
	left := genExpr(rng, depth-1)
	right := genExpr(rng, depth-1)
	op := ops[rng.Intn(len(ops))]
	return "(" + left + ")" + string(op) + "(" + right + ")"
}

func generateCase(rng *rand.Rand) string {
	expr := genExpr(rng, 3)
	n := rng.Intn(4)
	var sb strings.Builder
	sb.WriteString(expr + "\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	seen := make(map[string]bool)
	for i := 0; i < n; {
		a := rng.Intn(2)
		b := rng.Intn(2)
		c := rng.Intn(2)
		d := rng.Intn(2)
		key := fmt.Sprintf("%d%d%d%d", a, b, c, d)
		if seen[key] {
			continue
		}
		seen[key] = true
		e := rng.Intn(2)
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", a, b, c, d, e)
		i++
	}
	return sb.String()
}

func run(prog, input string) (string, error) {
	cmd := exec.Command(prog)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin, oracle, input string) error {
	got, err := run(bin, input)
	if err != nil {
		return err
	}
	expect, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
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
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
