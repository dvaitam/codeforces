package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseF struct {
	n   int
	arr []int64
	exp string
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solveF(n int, arr []int64) string {
	// Remove zeros: each zero is automatically a group of size 1 that sums to 0.
	nonzero := make([]int64, 0, n)
	for _, v := range arr {
		if v != 0 {
			nonzero = append(nonzero, v)
		}
	}
	N := len(nonzero)
	if N == 0 {
		return fmt.Sprint(0)
	}

	// A subset is "valid" (can be merged into one group) if there exists a
	// non-empty proper sub-subset whose sum equals half the total (with
	// appropriate parity). Specifically, assign +1/-1 to each element;
	// if the subset has k elements and sum S, a valid split exists iff
	// we can pick a sub-subset with sum T where 2T - S is in {-(k-1)...(k-1)}
	// with the right parity. Equivalently T in [(S-k+1)/2, (S+k-1)/2].
	valid := make([]bool, 1<<N)

	vals := make([]int64, N)
	posSuffix := make([]int64, N+1)
	negSuffix := make([]int64, N+1)

	var dfs func(idx int, currentSum int64, count int, k int, minT, maxT int64) bool
	dfs = func(idx int, currentSum int64, count int, k int, minT, maxT int64) bool {
		if currentSum >= minT && currentSum <= maxT {
			if count > 0 && count < k {
				return true
			}
		}
		if idx == k {
			return false
		}
		if currentSum+posSuffix[idx] < minT {
			return false
		}
		if currentSum+negSuffix[idx] > maxT {
			return false
		}
		if dfs(idx+1, currentSum+vals[idx], count+1, k, minT, maxT) {
			return true
		}
		if dfs(idx+1, currentSum, count, k, minT, maxT) {
			return true
		}
		return false
	}

	for mask := 1; mask < (1 << N); mask++ {
		k := 0
		var S int64 = 0
		for i := 0; i < N; i++ {
			if (mask & (1 << i)) != 0 {
				vals[k] = nonzero[i]
				S += nonzero[i]
				k++
			}
		}
		if k < 2 {
			continue
		}
		if abs64(S)%2 != int64(k-1)%2 {
			continue
		}

		minT := (S - int64(k) + 1) / 2
		maxT := (S + int64(k) - 1) / 2

		posSuffix[k] = 0
		negSuffix[k] = 0
		for i := k - 1; i >= 0; i-- {
			posSuffix[i] = posSuffix[i+1]
			negSuffix[i] = negSuffix[i+1]
			if vals[i] > 0 {
				posSuffix[i] += vals[i]
			} else {
				negSuffix[i] += vals[i]
			}
		}

		valid[mask] = dfs(0, 0, 0, k, minT, maxT)
	}

	// DP: minimum number of unmatched elements = n - max pairs merged.
	// dp[mask] = max number of valid groups (each saves at least 1 from the count).
	dp := make([]int8, 1<<N)
	for mask := 1; mask < (1 << N); mask++ {
		lsb := mask & -mask
		rem := mask ^ lsb
		res := dp[rem]

		for subRem := rem; subRem > 0; subRem = (subRem - 1) & rem {
			if valid[subRem|lsb] {
				if cand := dp[rem^subRem] + 1; cand > res {
					res = cand
				}
			}
		}

		if valid[lsb] {
			if cand := dp[rem] + 1; cand > res {
				res = cand
			}
		}
		dp[mask] = res
	}

	ans := N - int(dp[(1<<N)-1])
	return fmt.Sprint(ans)
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

func generateTests() []testCaseF {
	rng := rand.New(rand.NewSource(7))
	cases := make([]testCaseF, 100)
	for i := range cases {
		n := rng.Intn(6) + 1
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = int64(rng.Intn(10))
		}
		cases[i] = testCaseF{n: n, arr: arr, exp: solveF(n, arr)}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(tc.arr[j], 10))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
