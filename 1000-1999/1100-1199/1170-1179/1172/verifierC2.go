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

const MOD = 998244353

func modPow(a, e int) int {
	res := 1
	base := a % MOD
	if base < 0 {
		base += MOD
	}
	for e > 0 {
		if e&1 == 1 {
			res = int((int64(res) * int64(base)) % int64(MOD))
		}
		base = int((int64(base) * int64(base)) % int64(MOD))
		e >>= 1
	}
	return res
}

func modInv(a int) int {
	return modPow((a%MOD+MOD)%MOD, MOD-2)
}

func solveC2(n, m int, a []int, w []int) []int {
	var S_plus, S_minus int
	for i := 0; i < n; i++ {
		if a[i] == 1 {
			S_plus += w[i]
		} else {
			S_minus += w[i]
		}
	}
	S := S_plus + S_minus
	sMin := S - m
	if sMin < S_plus {
		sMin = S_plus
	}
	sMax := S + m
	invSize := sMax - sMin + 1
	inv := make([]int, invSize)
	for i := 0; i < invSize; i++ {
		tot := sMin + i
		inv[i] = modInv(tot)
	}
	dp := make([]int, m+1)
	ndp := make([]int, m+1)
	dp[0] = 1
	for k := 0; k < m; k++ {
		for i := 0; i <= m; i++ {
			ndp[i] = 0
		}
		for l := 0; l <= k; l++ {
			v := dp[l]
			if v == 0 {
				continue
			}
			S_total := S + 2*l - k
			idx := S_total - sMin
			invTot := inv[idx]
			pPlus := int((int64(S_plus+l) * int64(invTot)) % int64(MOD))
			ndp[l+1] = (ndp[l+1] + int((int64(v)*int64(pPlus))%int64(MOD))) % MOD
			rem := S_minus - (k - l)
			if rem > 0 {
				pMinus := int((int64(rem) * int64(invTot)) % int64(MOD))
				ndp[l] = (ndp[l] + int((int64(v)*int64(pMinus))%int64(MOD))) % MOD
			}
		}
		dp, ndp = ndp, dp
	}
	EX := 0
	for l := 0; l <= m; l++ {
		EX = (EX + int((int64(dp[l])*int64(l))%int64(MOD))) % MOD
	}
	invSp := modInv(S_plus)
	tPlus := int((int64(S_plus+EX) * int64(invSp)) % int64(MOD))
	mMinusEX := (m - EX) % MOD
	if mMinusEX < 0 {
		mMinusEX += MOD
	}
	invSm := 1
	if S_minus > 0 {
		invSm = modInv(S_minus)
	}
	tMinus := 0
	if S_minus > 0 {
		tMinus = int((int64((S_minus%MOD-mMinusEX+MOD)%MOD) * int64(invSm)) % int64(MOD))
	}
	results := make([]int, n)
	for i := 0; i < n; i++ {
		if a[i] == 1 {
			results[i] = int((int64(w[i]) * int64(tPlus)) % int64(MOD))
		} else {
			results[i] = int((int64(w[i]) * int64(tMinus)) % int64(MOD))
		}
	}
	return results
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(2)
	}
	// Ensure at least one liked picture
	hasLiked := false
	for _, v := range a {
		if v == 1 {
			hasLiked = true
			break
		}
	}
	if !hasLiked {
		a[rng.Intn(n)] = 1
	}
	for i, v := range a {
		sb.WriteString(fmt.Sprintf("%d", v))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	w := make([]int, n)
	for i := 0; i < n; i++ {
		w[i] = rng.Intn(5) + 1
		sb.WriteString(fmt.Sprintf("%d", w[i]))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func oracleFromInput(input string) string {
	// parse input
	var n, m int
	r := strings.NewReader(input)
	fmt.Fscan(r, &n, &m)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
	}
	w := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &w[i])
	}
	results := solveC2(n, m, a, w)
	parts := make([]string, n)
	for i, v := range results {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, " ")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		expect := oracleFromInput(input)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\n got: %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
