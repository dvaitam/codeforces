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

const mod = 1000000007

func modExp(base, exp int64) int64 {
	res := int64(1)
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		exp >>= 1
	}
	return res
}

func solveCase(n int, x int64, a []int64) int64 {
	var sumA int64
	for _, v := range a {
		sumA += v
	}
	maxA := a[n-1]
	cnt := make(map[int64]int64, n)
	for _, v := range a {
		d := maxA - v
		cnt[d]++
	}
	var carry int64
	var vT int64
	for j := int64(0); ; j++ {
		total := carry
		if c, ok := cnt[j]; ok {
			total += c
		}
		if total%x != 0 {
			vT = j
			break
		}
		carry = total / x
	}
	exp := (sumA - maxA) + vT
	if exp > sumA {
		exp = sumA
	}
	return modExp(x, exp)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	x := int64(rng.Intn(10) + 2)
	a := make([]int64, n)
	base := int64(rng.Intn(5))
	a[0] = base
	for i := 1; i < n; i++ {
		base += int64(rng.Intn(3))
		a[i] = base
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	exp := fmt.Sprintf("%d", solveCase(n, x, a))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
