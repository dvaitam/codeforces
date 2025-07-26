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

const mod int64 = 1000000007

var steps [1005]int

func popcount(x int) int {
	cnt := 0
	for x > 0 {
		cnt += x & 1
		x >>= 1
	}
	return cnt
}

func calc(x int) int {
	if steps[x] != -1 {
		return steps[x]
	}
	if x == 1 {
		steps[x] = 0
		return 0
	}
	steps[x] = calc(popcount(x)) + 1
	return steps[x]
}

func expected(nStr string, k int) int64 {
	if k == 0 {
		return 1
	}
	for i := 0; i < len(steps); i++ {
		steps[i] = -1
	}
	steps[1] = 0
	for i := 2; i < len(steps); i++ {
		calc(i)
	}
	L := len(nStr)
	comb := make([][]int64, L+1)
	for i := 0; i <= L; i++ {
		comb[i] = make([]int64, L+1)
	}
	for i := 0; i <= L; i++ {
		comb[i][0] = 1
		comb[i][i] = 1
		for j := 1; j < i; j++ {
			comb[i][j] = (comb[i-1][j-1] + comb[i-1][j]) % mod
		}
	}
	count := func(r int) int64 {
		if r < 0 {
			return 0
		}
		ones := 0
		var ans int64
		for i := 0; i < L; i++ {
			if nStr[i] == '1' {
				rem := L - i - 1
				need := r - ones
				if need >= 0 && need <= rem {
					ans = (ans + comb[rem][need]) % mod
				}
				ones++
			}
		}
		if ones == r {
			ans = (ans + 1) % mod
		}
		return ans
	}
	var ans int64
	for r := 1; r <= 1000; r++ {
		if steps[r] == k-1 {
			ans = (ans + count(r)) % mod
		}
	}
	if k == 1 {
		ans = (ans - 1 + mod) % mod
	}
	return ans
}

type testCase struct {
	n string
	k int
}

func generateCase(rng *rand.Rand) (string, testCase) {
	L := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteByte('1')
	for i := 1; i < L; i++ {
		if rng.Intn(2) == 1 {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	nStr := sb.String()
	k := rng.Intn(7) // keep k small
	input := fmt.Sprintf("%s\n%d\n", nStr, k)
	return input, testCase{n: nStr, k: k}
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, tc := generateCase(rng)
		exp := fmt.Sprintf("%d", expected(tc.n, tc.k))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
