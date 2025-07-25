package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Block struct{ a, b, c, d int }

func countPaths(n, m int, blocks map[[4]int]bool) *big.Int {
	dp := make([][]*big.Int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]*big.Int, m+1)
		for j := 0; j <= m; j++ {
			dp[i][j] = new(big.Int)
		}
	}
	dp[1][1].SetInt64(1)
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if i == 1 && j == 1 {
				continue
			}
			if i > 1 {
				if !blocks[[4]int{i - 1, j, i, j}] {
					dp[i][j].Add(dp[i][j], dp[i-1][j])
				}
			}
			if j > 1 {
				if !blocks[[4]int{i, j - 1, i, j}] {
					dp[i][j].Add(dp[i][j], dp[i][j-1])
				}
			}
		}
	}
	return dp[n][m]
}

func generateCase(rng *rand.Rand) *big.Int {
	// random T up to about 1e6
	val := rng.Int63n(1_000_000) + 1
	return big.NewInt(val)
}

func runCase(bin string, T *big.Int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%s\n", T.String()))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	vals := []int{}
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("non-integer output")
		}
		vals = append(vals, v)
	}
	if len(vals) < 3 || (len(vals)-2)%4 != 0 {
		return fmt.Errorf("invalid output format")
	}
	n := vals[0]
	m := vals[1]
	if n <= 0 || m <= 0 || n > 50 || m > 50 {
		return fmt.Errorf("invalid grid size")
	}
	k := vals[2]
	if k < 0 || k*4 != len(vals)-3 {
		return fmt.Errorf("k mismatch")
	}
	blocks := make(map[[4]int]bool)
	idx := 3
	for i := 0; i < k; i++ {
		a, b, c, d := vals[idx], vals[idx+1], vals[idx+2], vals[idx+3]
		idx += 4
		if a < 1 || a > n || c < 1 || c > n || b < 1 || b > m || d < 1 || d > m {
			return fmt.Errorf("invalid coordinates")
		}
		if !((a == c && abs(b-d) == 1) || (b == d && abs(a-c) == 1)) {
			return fmt.Errorf("non-adjacent rooms")
		}
		blocks[[4]int{a, b, c, d}] = true
		blocks[[4]int{c, d, a, b}] = true
	}
	paths := countPaths(n, m, blocks)
	if paths.Cmp(T) != 0 {
		return fmt.Errorf("expected %s paths got %s", T.String(), paths.String())
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		T := generateCase(rng)
		if err := runCase(bin, T); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
