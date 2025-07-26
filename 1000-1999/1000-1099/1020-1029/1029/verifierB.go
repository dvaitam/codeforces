package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TestCase struct {
	Input  string
	Output string
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveB(a []int) string {
	if len(a) == 0 {
		return "0"
	}
	ans := 1
	cur := 1
	for i := 1; i < len(a); i++ {
		if a[i] <= 2*a[i-1] {
			cur++
		} else {
			cur = 1
		}
		if cur > ans {
			ans = cur
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []TestCase {
	rand.Seed(43)
	tests := make([]TestCase, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(20) + 1
		arr := make([]int, n)
		cur := rand.Intn(20) + 1
		arr[0] = cur
		for i := 1; i < n; i++ {
			cur += rand.Intn(20) + 1
			arr[i] = cur
		}
		inputBuilder := strings.Builder{}
		inputBuilder.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				inputBuilder.WriteByte(' ')
			}
			inputBuilder.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		inputBuilder.WriteByte('\n')
		output := solveB(arr)
		tests[t] = TestCase{Input: inputBuilder.String(), Output: output}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(bin, tc.Input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(tc.Output) {
			fmt.Fprintf(os.Stderr, "Test %d failed:\ninput:\n%s\nexpected:%s\n got:%s\n", i+1, tc.Input, tc.Output, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
