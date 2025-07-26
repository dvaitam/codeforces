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

type testCaseD struct {
	input    string
	expected int64
}

func computeD(a []int, k int) int64 {
	n := len(a)
	var ans int64
	for i := 0; i < n; i++ {
		sum := 0
		prod := int64(1)
		for j := i; j < n; j++ {
			sum += a[j]
			if a[j] > 0 && prod > 2_000_000_000_000_000_000/int64(a[j]) {
				break
			}
			prod *= int64(a[j])
			if prod%int64(sum) == 0 && prod/int64(sum) == int64(k) {
				ans++
			}
		}
	}
	return ans
}

func generateCaseD(rng *rand.Rand) testCaseD {
	n := rng.Intn(6) + 1
	k := rng.Intn(5) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(5) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCaseD{input: sb.String(), expected: computeD(a, k)}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		tc := generateCaseD(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, tc.input)
			os.Exit(1)
		}
		var val int64
		if _, err := fmt.Sscan(out, &val); err != nil || val != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %s\ninput:\n%s", i, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
