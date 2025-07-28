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

const numTestsD = 100

func generateTestsD() []string {
	rng := rand.New(rand.NewSource(4))
	tests := make([]string, numTestsD)
	for i := 0; i < numTestsD; i++ {
		n := rng.Intn(10) + 1
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Int63n(100) + 1
		}
		tVal := rng.Int63n(50) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[j]))
		}
		sb.WriteString("\n1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tVal))
		tests[i] = sb.String()
	}
	return tests
}

func solveD(input string) string {
	reader := strings.NewReader(input)
	var n int
	fmt.Fscan(reader, &n)
	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		pref[i] = pref[i-1] + x
	}
	maxTime := int64(0)
	for i := 1; i <= n; i++ {
		t := (pref[i] + int64(i) - 1) / int64(i)
		if t > maxTime {
			maxTime = t
		}
	}
	var q int
	fmt.Fscan(reader, &q)
	var results []string
	for ; q > 0; q-- {
		var t int64
		fmt.Fscan(reader, &t)
		if t < maxTime {
			results = append(results, "-1")
		} else {
			k := (pref[n] + t - 1) / t
			results = append(results, fmt.Sprintf("%d", k))
		}
	}
	return strings.Join(results, "\n") + "\n"
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
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsD()
	for i, tc := range tests {
		expected := solveD(tc)
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
