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
	oracle := filepath.Join(dir, "oracleG")
	cmd := exec.Command("go", "build", "-o", oracle, "1009G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(bin, input string) (string, error) {
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

func genCaseG(rng *rand.Rand) (string, int, [][2]interface{}) {
	n := rng.Intn(10) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(6))
	}
	m := rng.Intn(n*2 + 1)
	ops := make([][2]interface{}, m)
	for i := 0; i < m; i++ {
		pos := rng.Intn(n) + 1
		mask := 0
		letters := rng.Intn(1 << 6)
		// ensure some letters chosen
		mask = letters
		ops[i][0] = pos
		ops[i][1] = mask
	}
	return string(b), m, ops
}

func formatOps(ops [][2]interface{}) string {
	var sb strings.Builder
	for i, op := range ops {
		if i > 0 {
			sb.WriteByte('\n')
		}
		pos := op[0].(int)
		mask := op[1].(int)
		sb.WriteString(strconv.Itoa(pos))
		sb.WriteByte(' ')
		if mask == 0 {
			sb.WriteByte('\n')
			continue
		}
		for j := 0; j < 6; j++ {
			if mask>>j&1 == 1 {
				sb.WriteByte(byte('a' + j))
			}
		}
	}
	if len(ops) > 0 {
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCaseG(bin, oracle, s string, m int, ops [][2]interface{}) error {
	var sb strings.Builder
	sb.WriteString(s)
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(m))
	sb.WriteByte('\n')
	sb.WriteString(formatOps(ops))
	input := sb.String()
	exp, err := runProg(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	got, err := runProg(bin, input)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
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
		s, m, ops := genCaseG(rng)
		if err := runCaseG(bin, oracle, s, m, ops); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
