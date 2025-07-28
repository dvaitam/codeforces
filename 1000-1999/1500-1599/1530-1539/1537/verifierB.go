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

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solveCase(n, m, i, j int64) string {
	d11 := abs64(i-1) + abs64(j-1)
	dnn := abs64(i-n) + abs64(j-m)
	d1m := abs64(i-1) + abs64(j-m)
	dn1 := abs64(i-n) + abs64(j-1)
	if d11+dnn >= d1m+dn1 {
		return fmt.Sprintf("1 1 %d %d\n", n, m)
	}
	return fmt.Sprintf("1 %d %d 1\n", m, n)
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func generateCase(rng *rand.Rand) (string, string) {
	n := int64(rng.Intn(20) + 1)
	m := int64(rng.Intn(20) + 1)
	i := int64(rng.Intn(int(n)) + 1)
	j := int64(rng.Intn(int(m)) + 1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d %d\n", n, m, i, j)
	expect := solveCase(n, m, i, j)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), strings.TrimSpace(out), in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
