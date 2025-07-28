package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const numTestsC = 100

func generateTestsC() []string {
	rng := rand.New(rand.NewSource(3))
	tests := make([]string, numTestsC)
	for i := 0; i < numTestsC; i++ {
		n := rng.Intn(10) + 1
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Int63n(2001) - 1000
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[j]))
		}
		sb.WriteByte('\n')
		tests[i] = sb.String()
	}
	return tests
}

func solveC(input string) string {
	var t, n int
	reader := strings.NewReader(input)
	fmt.Fscan(reader, &t)
	fmt.Fscan(reader, &n)
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	var pos, neg int64
	for i := 1; i < n; i++ {
		diff := arr[i] - arr[i-1]
		if diff > 0 {
			pos += diff
		} else {
			neg += -diff
		}
	}
	ans := pos + neg + int64(math.Abs(float64(arr[0]-neg)))
	return fmt.Sprintf("%d\n", ans)
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsC()
	for i, tc := range tests {
		expected := solveC(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("Test %d: error running binary: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s", i+1, tc, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
