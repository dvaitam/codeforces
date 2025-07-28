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

type testCase struct {
	input    string
	expected int64
}

func computeG(n int, mod int64) [][]int64 {
	g := make([][]int64, n+1)
	g[0] = []int64{0}
	arrPrev := []int64{1}
	basePrev := 0
	for L := 1; L <= n; L++ {
		maxCurr := L * (L - 1) / 2
		baseCurr := maxCurr
		arrCurr := make([]int64, 2*maxCurr+1)
		sizePrev := len(arrPrev)
		prefix0 := make([]int64, sizePrev+1)
		prefix1 := make([]int64, sizePrev+1)
		for i := 0; i < sizePrev; i++ {
			prefix0[i+1] = (prefix0[i] + arrPrev[i]) % mod
			prefix1[i+1] = (prefix1[i] + int64(i-basePrev)*arrPrev[i]) % mod
		}
		getSum := func(a, b int) (int64, int64) {
			if a < -basePrev {
				a = -basePrev
			}
			if b > basePrev {
				b = basePrev
			}
			if a > b {
				return 0, 0
			}
			idxA := a + basePrev
			idxB := b + basePrev
			sum0 := (prefix0[idxB+1] - prefix0[idxA]) % mod
			if sum0 < 0 {
				sum0 += mod
			}
			sum1 := (prefix1[idxB+1] - prefix1[idxA]) % mod
			if sum1 < 0 {
				sum1 += mod
			}
			return sum0, sum1
		}
		for d := -maxCurr; d <= maxCurr; d++ {
			s1, s2 := getSum(d-(L-1), d)
			s3, s4 := getSum(d+1, d+(L-1))
			val := (int64(L) - int64(d)) % mod * s1 % mod
			val = (val + s2) % mod
			val = (val + (int64(L)+int64(d))%mod*s3%mod) % mod
			val = (val - s4) % mod
			if val < 0 {
				val += mod
			}
			arrCurr[d+baseCurr] = val
		}
		gL := make([]int64, L+1)
		prefixPos := make([]int64, maxCurr+2)
		for d := maxCurr; d >= 0; d-- {
			prefixPos[d] = (prefixPos[d+1] + arrCurr[d+baseCurr]) % mod
		}
		for delta := 0; delta <= L; delta++ {
			if delta+1 <= maxCurr {
				gL[delta] = prefixPos[delta+1]
			} else {
				gL[delta] = 0
			}
		}
		g[L] = gL
		arrPrev = arrCurr
		basePrev = baseCurr
	}
	return g
}

func solve(n int, mod int64) int64 {
	g := computeG(n, mod)
	ans := int64(0)
	prefix := int64(1)
	for j := 1; j <= n; j++ {
		if j > 1 {
			prefix = prefix * int64(n-j+2) % mod
		}
		m := n - j + 1
		L := m - 1
		for delta := 1; delta <= m-1; delta++ {
			ans = (ans + prefix*int64(m-delta)%mod*g[L][delta]) % mod
		}
	}
	return ans % mod
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	mod := int64(rng.Intn(1_000_000_000-1) + 2)
	input := fmt.Sprintf("%d %d\n", n, mod)
	return testCase{input: input, expected: solve(n, mod)}
}

func runCase(bin string, tc testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{
		{input: "1 2\n", expected: solve(1, 2)},
		{input: "3 1000000007\n", expected: solve(3, 1000000007)},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
