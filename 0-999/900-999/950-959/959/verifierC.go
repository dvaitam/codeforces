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

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func solveC(n int) string {
	var sb strings.Builder
	if n <= 5 {
		sb.WriteString("-1\n")
	} else {
		tmp := n - 2
		sb.WriteString("1 2\n")
		tot := tmp / 2
		for i := 1; i <= tot; i++ {
			sb.WriteString(fmt.Sprintf("1 %d\n", i+2))
		}
		for i := 1; i <= tmp-tot; i++ {
			sb.WriteString(fmt.Sprintf("2 %d\n", i+tot+2))
		}
	}
	for i := 2; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("1 %d\n", i))
	}
	return strings.TrimSpace(sb.String())
}

func genCaseC(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	return fmt.Sprintf("%d\n", n)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCaseC(rng)
		n := 0
		fmt.Sscan(strings.TrimSpace(in), &n)
		expect := solveC(n)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, expect, strings.TrimSpace(got), in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
