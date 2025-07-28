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

type testCaseC struct {
	n int
	k int
}

func generateCaseC(rng *rand.Rand) (string, testCaseC) {
	n := rng.Intn(29) + 2 // 2..30
	maxK := n * (n + 1) / 2
	k := rng.Intn(maxK + 1)
	input := fmt.Sprintf("1\n%d %d\n", n, k)
	return input, testCaseC{n: n, k: k}
}

func solveCaseC(n, k int) []int {
	a := make([]int, n)
	for i := range a {
		a[i] = -1000
	}
	temp := 0
	for i := 1; i <= n; i++ {
		if i*(i+1)/2 <= k {
			temp = i
		}
	}
	remK := k - temp*(temp+1)/2
	for i := 0; i < temp; i++ {
		a[i] = i + 2
	}
	if remK > 0 {
		sum := 0
		for i := remK - 1; i < temp; i++ {
			sum += a[i]
		}
		a[temp] = -(sum - 1)
	}
	return a
}

func expectedC(tc testCaseC) string {
	arr := solveCaseC(tc.n, tc.k)
	var sb strings.Builder
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String()
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, tc := generateCaseC(rng)
		expect := expectedC(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
