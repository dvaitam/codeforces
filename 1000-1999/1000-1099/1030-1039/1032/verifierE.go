package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(input string) string {
	r := strings.NewReader(input)
	var n int
	fmt.Fscan(r, &n)

	a := make([]int, n)
	counts := make(map[int]int)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
		counts[a[i]]++
	}

	if len(counts) <= 1 {
		return fmt.Sprintf("%d", n)
	}

	dp := make([][]*big.Int, n+1)
	for i := range dp {
		dp[i] = make([]*big.Int, 10001)
		for j := range dp[i] {
			dp[i][j] = new(big.Int)
		}
	}
	dp[0][0].SetInt64(1)

	currentSum := 0
	currentCnt := 0
	for _, w := range a {
		currentCnt++
		currentSum += w
		for i := currentCnt; i >= 1; i-- {
			for j := currentSum; j >= w; j-- {
				if dp[i-1][j-w].Sign() > 0 {
					dp[i][j].Add(dp[i][j], dp[i-1][j-w])
				}
			}
		}
	}

	maxRevealed := 0
	for x, c := range counts {
		for k := 1; k <= c; k++ {
			req := k * x
			ncr := big.NewInt(1)
			for i := 1; i <= k; i++ {
				ncr.Mul(ncr, big.NewInt(int64(c-i+1)))
				ncr.Div(ncr, big.NewInt(int64(i)))
			}

			if dp[k][req].Cmp(ncr) == 0 {
				revealed := k
				if k == c && len(counts) == 2 {
					revealed = n
				}
				if revealed > maxRevealed {
					maxRevealed = revealed
				}
			}
		}
	}

	return fmt.Sprintf("%d", maxRevealed)
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", r.Intn(20)+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	r := rand.New(rand.NewSource(1))
	cases := []string{
		"1\n1\n",
		"4\n1 2 3 4\n",
		"4\n2 2 1 4\n",
	}
	for i := 0; i < 97; i++ {
		cases = append(cases, genCase(r))
	}
	for idx, input := range cases {
		want := solve(input)
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
