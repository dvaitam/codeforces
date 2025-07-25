package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func expected(a [][]int) int {
	n := len(a[0])
	dp0 := a[0][0]
	dp1 := a[0][0] + a[1][0]
	dp2 := a[0][0] + a[1][0] + a[2][0]
	for i := 1; i < n; i++ {
		ndp0 := max(dp0+a[0][i], max(dp1+a[0][i]+a[1][i], dp2+a[0][i]+a[1][i]+a[2][i]))
		ndp1 := max(dp1+a[1][i], max(dp0+a[0][i]+a[1][i], dp2+a[1][i]+a[2][i]))
		ndp2 := max(dp2+a[2][i], max(dp1+a[1][i]+a[2][i], dp0+a[0][i]+a[1][i]+a[2][i]))
		dp0, dp1, dp2 = ndp0, ndp1, ndp2
	}
	return dp2
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))

	for t := 0; t < 100; t++ {
		n := rng.Intn(10) + 1
		a := make([][]int, 3)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < 3; i++ {
			a[i] = make([]int, n)
			for j := 0; j < n; j++ {
				val := rng.Intn(21) - 10
				a[i][j] = val
				sb.WriteString(fmt.Sprintf("%d ", val))
			}
			sb.WriteString("\n")
		}
		input := sb.String()
		exp := fmt.Sprintf("%d", expected(a))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\noutput:\n%s\n", t+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
