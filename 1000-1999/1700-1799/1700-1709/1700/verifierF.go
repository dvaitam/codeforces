package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const numTestsF = 100

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solveF(input string) string {
	reader := strings.NewReader(input)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	a := make([][]int, 2)
	b := make([][]int, 2)
	for i := 0; i < 2; i++ {
		a[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &a[i][j])
		}
	}
	for i := 0; i < 2; i++ {
		b[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &b[i][j])
		}
	}
	sumA, sumB := 0, 0
	for i := 0; i < 2; i++ {
		for j := 0; j < n; j++ {
			sumA += a[i][j]
			sumB += b[i][j]
		}
	}
	if sumA != sumB {
		return "-1\n"
	}
	c0, c1 := 0, 0
	cost := 0
	for i := 0; i < n; i++ {
		u := c0 + a[0][i] - b[0][i]
		v := c1 + a[1][i] - b[1][i]
		arr := []int{0, u, -v}
		sort.Ints(arr)
		x := arr[1]
		cost += abs(x) + abs(u-x) + abs(v+x)
		c0 = u - x
		c1 = v + x
	}
	if c0 != 0 || c1 != 0 {
		return "-1\n"
	}
	return fmt.Sprintf("%d\n", cost)
}

func generateTestsF() []string {
	rng := rand.New(rand.NewSource(6))
	tests := make([]string, numTestsF)
	for i := 0; i < numTestsF; i++ {
		n := rng.Intn(6) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for r := 0; r < 2; r++ {
			for c := 0; c < n; c++ {
				sb.WriteString(fmt.Sprintf("%d", rng.Intn(2)))
				if c+1 < n {
					sb.WriteByte(' ')
				}
			}
			sb.WriteByte('\n')
		}
		for r := 0; r < 2; r++ {
			for c := 0; c < n; c++ {
				sb.WriteString(fmt.Sprintf("%d", rng.Intn(2)))
				if c+1 < n {
					sb.WriteByte(' ')
				}
			}
			sb.WriteByte('\n')
		}
		tests[i] = sb.String()
	}
	return tests
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
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsF()
	for i, tc := range tests {
		expected := solveF(tc)
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
