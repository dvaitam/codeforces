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

const mod int64 = 998244353

func solveCase(n, k int) []int64 {
	dp := make([]int64, n+1)
	dp[0] = 1
	ans := make([]int64, n+1)
	for step := 0; ; step++ {
		base := k + step
		if base > n {
			break
		}
		newdp := make([]int64, n+1)
		for r := 0; r < base; r++ {
			sum := int64(0)
			for x := r; x+base <= n; x += base {
				sum += dp[x]
				if sum >= mod {
					sum -= mod
				}
				newdp[x+base] = sum
			}
		}
		has := false
		for i := base; i <= n; i++ {
			if newdp[i] != 0 {
				has = true
			}
			ans[i] += newdp[i]
			if ans[i] >= mod {
				ans[i] -= mod
			}
		}
		if !has {
			break
		}
		dp = newdp
	}
	return ans[1:]
}

func runCandidate(bin, input string) (string, error) {
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

func verifyCase(bin string, n, k int) error {
	input := fmt.Sprintf("%d %d\n", n, k)
	expectedVals := solveCase(n, k)
	expectedStr := make([]string, len(expectedVals))
	for i, v := range expectedVals {
		expectedStr[i] = fmt.Sprint(v)
	}
	expected := strings.Join(expectedStr, " ")
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != expected {
		return fmt.Errorf("expected %s got %s", expected, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(50) + 1
		k := rng.Intn(n) + 1
		if err := verifyCase(bin, n, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
