package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod = 998244353

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func solveTwo(a, b int) int64 {
	S := a + b
	if S <= 1 {
		return 0
	}
	n := S - 1
	A := make([][]*big.Rat, n)
	B := make([]*big.Rat, n)
	for i := 0; i < n; i++ {
		A[i] = make([]*big.Rat, n)
		B[i] = new(big.Rat).SetInt64(1)
	}
	for k := 1; k < S; k++ {
		row := k - 1
		A[row][row] = new(big.Rat).SetInt64(1)
		if k-1 >= 1 {
			A[row][row-1] = new(big.Rat).SetFrac(big.NewInt(int64(-k)), big.NewInt(int64(S)))
		}
		if k+1 <= S-1 {
			A[row][row+1] = new(big.Rat).SetFrac(big.NewInt(int64(k-S)), big.NewInt(int64(S)))
		}
	}
	// Gaussian elimination
	for i := 0; i < n; i++ {
		pivot := i
		for pivot < n && (A[pivot][i] == nil || A[pivot][i].Sign() == 0) {
			pivot++
		}
		if pivot == n {
			continue
		}
		if pivot != i {
			A[i], A[pivot] = A[pivot], A[i]
			B[i], B[pivot] = B[pivot], B[i]
		}
		pv := new(big.Rat).Set(A[i][i])
		for j := i; j < n; j++ {
			if A[i][j] == nil {
				A[i][j] = new(big.Rat)
			}
			A[i][j].Quo(A[i][j], pv)
		}
		B[i].Quo(B[i], pv)
		for r := 0; r < n; r++ {
			if r == i {
				continue
			}
			if A[r][i] == nil || A[r][i].Sign() == 0 {
				continue
			}
			factor := new(big.Rat).Set(A[r][i])
			for j := i; j < n; j++ {
				if A[r][j] == nil {
					A[r][j] = new(big.Rat)
				}
				tmp := new(big.Rat).Mul(factor, A[i][j])
				A[r][j].Sub(A[r][j], tmp)
			}
			tmp := new(big.Rat).Mul(factor, B[i])
			B[r].Sub(B[r], tmp)
		}
	}
	res := B[a-1]
	num := new(big.Int).Mod(res.Num(), big.NewInt(mod))
	den := new(big.Int).Mod(res.Denom(), big.NewInt(mod))
	inv := modPow(den.Int64(), mod-2)
	return num.Int64() * inv % mod
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
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

func genCase(rng *rand.Rand) string {
	a := rng.Intn(4) + 1
	b := rng.Intn(4) + 1
	return fmt.Sprintf("2\n%d %d\n", a, b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		input := genCase(rng)
		var a, b int
		fmt.Sscanf(input, "2\n%d %d\n", &a, &b)
		exp := solveTwo(a, b)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(gotStr, &got)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected: %d\nGot: %d\ninput:\n%s", tcase+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
