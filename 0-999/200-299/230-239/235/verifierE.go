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

const MOD = 1 << 30

func solveCase(a, b, c int) int64 {
	ab := a * b
	cnt := make([]int32, ab+1)
	for i := 1; i <= a; i++ {
		for j := 1; j <= b; j++ {
			cnt[i*j]++
		}
	}
	D := make([]int32, ab+1)
	for u := 1; u <= ab; u++ {
		var sum int32
		for v := u; v <= ab; v += u {
			sum += cnt[v]
		}
		D[u] = sum
	}
	mu := make([]int8, c+1)
	isPrime := make([]bool, c+1)
	primes := make([]int, 0, c/10)
	if c >= 1 {
		mu[1] = 1
	}
	for i := 2; i <= c; i++ {
		isPrime[i] = true
	}
	for i := 2; i <= c; i++ {
		if isPrime[i] {
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			v := i * p
			if v > c {
				break
			}
			isPrime[v] = false
			if i%p == 0 {
				mu[v] = 0
				break
			}
			mu[v] = -mu[i]
		}
	}
	H := make([]int32, c+1)
	for m := 1; m <= c; m++ {
		var s int32
		for k := 1; k <= m; k++ {
			s += int32(m / k)
		}
		H[m] = s
	}
	T := make([]int64, ab+1)
	for d := 1; d <= c; d++ {
		md := mu[d]
		if md == 0 {
			continue
		}
		h := int64(H[c/d]) * int64(md)
		for u := d; u <= ab; u += d {
			T[u] += h
		}
	}
	var ans int64
	for u := 1; u <= ab; u++ {
		if D[u] != 0 && T[u] != 0 {
			ans += int64(D[u]) * T[u]
			if ans > (1<<61) || ans < -(1<<61) {
				ans %= MOD
			}
		}
	}
	ans %= MOD
	if ans < 0 {
		ans += MOD
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	a := rng.Intn(8) + 1
	b := rng.Intn(8) + 1
	c := rng.Intn(8) + 1
	input := fmt.Sprintf("%d %d %d\n", a, b, c)
	expect := fmt.Sprintf("%d", solveCase(a, b, c))
	return input, expect
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
