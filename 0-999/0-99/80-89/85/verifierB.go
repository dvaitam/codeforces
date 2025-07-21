package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	exe := "oracleB"
	cmd := exec.Command("go", "build", "-o", exe, "./0-999/0-99/80-89/85/85B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func generateCase(rng *rand.Rand) string {
	k1 := rng.Intn(3) + 1
	k2 := rng.Intn(3) + 1
	k3 := rng.Intn(3) + 1
	t1 := rng.Intn(10) + 1
	t2 := rng.Intn(10) + 1
	t3 := rng.Intn(10) + 1
	n := rng.Intn(8) + 1
	arr := make([]int, n)
	cur := rng.Intn(3) + 1
	arr[0] = cur
	for i := 1; i < n; i++ {
		cur += rng.Intn(3)
		arr[i] = cur
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", k1, k2, k3)
	fmt.Fprintf(&sb, "%d %d %d\n", t1, t2, t3)
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
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
		input := generateCase(rng)
		exp, err := runProg("./"+oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
