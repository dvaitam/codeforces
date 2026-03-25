package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveF(r *bufio.Reader) string {
	var n int
	var m int64
	if _, err := fmt.Fscan(r, &n); err != nil {
		return ""
	}
	if _, err := fmt.Fscan(r, &m); err != nil {
		return ""
	}

	if n == 1 {
		return fmt.Sprintf("%d", int64(1)%m)
	}

	MOD := m

	R := make([]int64, n+1)
	SumM := make([]int64, 2*n+2)
	PrefSum := make([]int64, 2*n+2)

	pow2 := int64(2) % MOD

	for length := 1; length <= n; length++ {
		if length == 1 {
			R[1] = 1 % MOD
			SumM[2] = (SumM[2] + R[1]) % MOD
		} else if length == 2 {
			R[2] = 1 % MOD
			SumM[4] = (SumM[4] + R[2]) % MOD
		} else {
			R[length] = pow2
			pow2 = (pow2 * 2) % MOD

			for L := 1; L < length; L++ {
				K := length - 2*L - 1
				curS := int64(0)
				if K >= 1 {
					curS = (R[L] * PrefSum[K]) % MOD
				}
				R[length] = (R[length] - curS + MOD) % MOD
				SumM[length+L] = (SumM[length+L] + curS) % MOD
			}
			SumM[2*length] = (SumM[2*length] + R[length]) % MOD
		}

		PrefSum[length+1] = (PrefSum[length] + SumM[length+1]) % MOD
	}

	return fmt.Sprintf("%d", R[n]%MOD)
}

func generateCaseF(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	return fmt.Sprintf("%d %d\n", n, m)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseF(rng)
		expect := solveF(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
