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
	cmd := exec.Command("go", "build", "-o", exe, "./0-999/0-99/80-89/87/87B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func randomName(rng *rand.Rand) string {
	for {
		l := rng.Intn(5) + 1
		b := make([]byte, l)
		for i := range b {
			b[i] = byte('a' + rng.Intn(26))
		}
		s := string(b)
		if s != "void" && s != "errtype" {
			return s
		}
	}
}

func randomTypeExpr(rng *rand.Rand, names []string) string {
	prefix := rng.Intn(3)
	suffix := rng.Intn(3)
	base := names[rng.Intn(len(names))]
	return strings.Repeat("&", prefix) + base + strings.Repeat("*", suffix)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	names := []string{"void", "errtype"}
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			// typedef
			a := randomTypeExpr(rng, names)
			b := randomName(rng)
			fmt.Fprintf(&sb, "typedef %s %s\n", a, b)
			names = append(names, b)
		} else {
			a := randomTypeExpr(rng, names)
			fmt.Fprintf(&sb, "typeof %s\n", a)
		}
	}
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
