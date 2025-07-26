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

type testCase struct {
	input    string
	expected string
}

func expected(a [6]int) string {
	sum := 0
	for _, v := range a {
		sum += v
	}
	if sum%2 != 0 {
		return "NO"
	}
	target := sum / 2
	for i := 0; i < 6; i++ {
		for j := i + 1; j < 6; j++ {
			for k := j + 1; k < 6; k++ {
				if a[i]+a[j]+a[k] == target {
					return "YES"
				}
			}
		}
	}
	return "NO"
}

func generateCase(rng *rand.Rand) testCase {
	var arr [6]int
	for i := 0; i < 6; i++ {
		arr[i] = rng.Intn(21) // 0..20
	}
	var sb strings.Builder
	for i := 0; i < 6; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(arr[i]))
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expected: expected(arr)}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		out = strings.ToUpper(strings.TrimSpace(out))
		if out != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d: expected %s got %s\ninput:\n%s", i+1, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
