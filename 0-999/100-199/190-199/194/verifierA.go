package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func solveA(r *bufio.Reader) string {
	var n, k int
	if _, err := fmt.Fscan(r, &n, &k); err != nil {
		return ""
	}
	for x2 := 0; x2 <= n; x2++ {
		minSum := 3*n - x2
		maxSum := 5*n - 3*x2
		if k >= minSum && k <= maxSum {
			return fmt.Sprintf("%d", x2)
		}
	}
	return ""
}

func generateCaseA(rng *rand.Rand) string {
	n := rng.Intn(50) + 1
	x2 := rng.Intn(n + 1)
	sum := 2 * x2
	for i := 0; i < n-x2; i++ {
		sum += rng.Intn(3) + 3
	}
	return fmt.Sprintf("%d %d\n", n, sum)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseA(rng)
		expect := solveA(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
