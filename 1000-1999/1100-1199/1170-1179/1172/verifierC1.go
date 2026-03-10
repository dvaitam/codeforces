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

func modpow(a, e int64) int64 {
	a %= MOD
	if a < 0 {
		a += MOD
	}
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modinv(a int64) int64 { return modpow(a, MOD-2) }

// refSolve computes expected weights via DP over number of liked-visits.
// dp[j] = probability (mod) that exactly j of the first 'step' visits were to liked pictures.
// At state (step, j): total weight = A+B + 2j - step, liked weight = A+j, disliked weight = B-(step-j).
func refSolve(n, m int, a []int, w []int64) []int64 {
	var A, B int64
	for i := 0; i < n; i++ {
		if a[i] == 1 {
			A += w[i]
		} else {
			B += w[i]
		}
	}

	dp := make([]int64, m+2)
	dp[0] = 1
	for step := 0; step < m; step++ {
		ndp := make([]int64, m+2)
		for j := 0; j <= step; j++ {
			if dp[j] == 0 {
				continue
			}
			total := A + B + int64(2*j-step)
			inv := modinv(total % MOD)
			likedW := A + int64(j)
			disW := B - int64(step-j)
			// liked transition: j -> j+1
			ndp[j+1] = (ndp[j+1] + dp[j]*((likedW%MOD)*inv%MOD)) % MOD
			// disliked transition: j -> j, only if disliked weight > 0
			if disW > 0 {
				ndp[j] = (ndp[j] + dp[j]*((disW%MOD)*inv%MOD)) % MOD
			}
		}
		dp = ndp
	}

	// E[A_final] and E[B_final]
	var EA, EB int64
	for j := 0; j <= m; j++ {
		if dp[j] == 0 {
			continue
		}
		EA = (EA + dp[j]*((A+int64(j))%MOD)) % MOD
		bFinal := B - int64(m-j)
		if bFinal >= 0 {
			EB = (EB + dp[j]*(bFinal%MOD)) % MOD
		}
	}

	invA := modinv(A % MOD)
	var invB int64
	if B > 0 {
		invB = modinv(B % MOD)
	}

	res := make([]int64, n)
	for i := 0; i < n; i++ {
		if a[i] == 1 {
			res[i] = w[i] % MOD * invA % MOD * EA % MOD
		} else {
			res[i] = w[i] % MOD * invB % MOD * EB % MOD
		}
	}
	return res
}

func formatExpected(n, m int, a []int, w []int64) string {
	res := refSolve(n, m, a, w)
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	return sb.String()
}

type testCase struct {
	n, m int
	a    []int
	w    []int64
}

func (tc testCase) Input() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range tc.w {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	a := make([]int, n)
	hasLiked := false
	for i := range a {
		a[i] = rng.Intn(2)
		if a[i] == 1 {
			hasLiked = true
		}
	}
	if !hasLiked {
		a[rng.Intn(n)] = 1
	}
	w := make([]int64, n)
	for i := range w {
		w[i] = int64(rng.Intn(5) + 1)
	}
	return testCase{n: n, m: m, a: a, w: w}
}

func runProg(bin, input string) (string, error) {
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
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tc := generateCase(rng)
		expect := formatExpected(tc.n, tc.m, tc.a, tc.w)
		got, err := runProg(bin, tc.Input())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\n     got: %s\ninput:\n%s", i+1, expect, got, tc.Input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
