package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveD(arr []int) int {
	n := len(arr) - 1
	can := make([][]bool, n+2)
	for i := range can {
		can[i] = make([]bool, n+2)
	}
	for i := 1; i <= n+1; i++ {
		can[i][i-1] = true
	}
	for l := 1; l <= n; l++ {
		freq := make([]int, n+1)
		maxf := 0
		for r := l; r <= n; r++ {
			x := arr[r]
			freq[x]++
			if freq[x] > maxf {
				maxf = freq[x]
			}
			length := r - l + 1
			if length%2 == 0 && maxf <= length/2 {
				can[l][r] = true
			}
		}
	}
	dp := make([]int, n+1)
	ans := 0
	for i := 1; i <= n; i++ {
		if can[1][i-1] {
			dp[i] = 1
		}
		for j := 1; j < i; j++ {
			if arr[i] == arr[j] && dp[j] > 0 && can[j+1][i-1] {
				if dp[j]+1 > dp[i] {
					dp[i] = dp[j] + 1
				}
			}
		}
		if can[i+1][n] && dp[i] > ans {
			ans = dp[i]
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 1
	input := fmt.Sprintf("1\n%d\n", n)
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		arr[i] = rng.Intn(n) + 1
		input += fmt.Sprintf("%d", arr[i])
		if i < n {
			input += " "
		}
	}
	input += "\n"
	exp := solveD(arr)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output\ninput:\n%soutput:\n%s", i+1, in, out)
			os.Exit(1)
		}
		if val != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, val, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
