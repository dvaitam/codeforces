package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type matrix [3][3]int64

func runCmd(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solveRef(t matrix, n int) int64 {
	var dp [41][3][3]int64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == j {
				dp[1][i][j] = 0
			} else {
				k := 3 - i - j
				dp[1][i][j] = min(t[i][j], t[i][k]+t[k][j])
			}
		}
	}
	for d := 2; d <= n; d++ {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if i == j {
					dp[d][i][j] = 0
					continue
				}
				k := 3 - i - j
				cost1 := dp[d-1][i][k] + t[i][j] + dp[d-1][k][j]
				cost2 := dp[d-1][i][j] + t[i][k] + dp[d-1][j][i] + t[k][j] + dp[d-1][i][j]
				dp[d][i][j] = min(cost1, cost2)
			}
		}
	}
	return dp[n][0][2]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rand.Seed(1)
	for tc := 1; tc <= 100; tc++ {
		var t matrix
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if i == j {
					t[i][j] = 0
				} else {
					t[i][j] = int64(rand.Intn(10000) + 1)
				}
			}
		}
		n := rand.Intn(40) + 1
		inputBuilder := strings.Builder{}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				inputBuilder.WriteString(fmt.Sprintf("%d ", t[i][j]))
			}
			inputBuilder.WriteByte('\n')
		}
		inputBuilder.WriteString(fmt.Sprintf("%d\n", n))
		input := inputBuilder.String()

		expected := fmt.Sprintf("%d", solveRef(t, n))
		got, err := runCmd(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", tc, err)
			return
		}
		if got != expected {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", tc, input, expected, got)
			return
		}
	}
	fmt.Println("All tests passed!")
}
