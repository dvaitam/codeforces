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

const MOD int64 = 998244353

type testCaseG struct {
	n int
	k int64
	a []int64
}

func generateCaseG(rng *rand.Rand) (string, testCaseG) {
	n := rng.Intn(5) + 2
	k := int64(rng.Intn(10))
	a := make([]int64, n)
	cur := int64(0)
	for i := 0; i < n; i++ {
		cur += int64(rng.Intn(5))
		a[i] = cur
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), testCaseG{n, k, a}
}

func solveCaseG(n int, k int64, a []int64) int64 {
	if n == 1 {
		return 1
	}
	groups := []int{}
	start := 0
	for i := 1; i < n; i++ {
		if a[i]-a[i-1] > k {
			groups = append(groups, i-start)
			start = i
		}
	}
	groups = append(groups, n-start)
	if len(groups) == 1 {
		return 0
	}
	fact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	sumBad := int64(0)
	for _, s := range groups {
		sumBad = (sumBad + int64(s)*(int64(s)-1)) % MOD
	}
	ans := (fact[n] - fact[n-2]*sumBad) % MOD
	if ans < 0 {
		ans += MOD
	}
	return ans
}

func expectedG(tc testCaseG) string {
	ans := solveCaseG(tc.n, tc.k, tc.a)
	return fmt.Sprintf("%d", ans)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, tc := generateCaseG(rng)
		expect := expectedG(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
