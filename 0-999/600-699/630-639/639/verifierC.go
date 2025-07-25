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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(n int, k int64, a []int64) int {
	p := big.NewInt(0)
	for i := n; i >= 0; i-- {
		p.Lsh(p, 1)
		p.Add(p, big.NewInt(a[i]))
	}
	y := new(big.Int).Set(p)
	divisible := true
	ans := 0
	for j := 0; j <= n; j++ {
		if divisible && y.BitLen() <= 62 {
			val := y.Int64()
			diff := a[j] - val
			if diff >= -k && diff <= k {
				if j != n || diff != 0 {
					ans++
				}
			}
		}
		divisible = divisible && y.Bit(0) == 0
		y.Rsh(y, 1)
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	k := int64(rng.Intn(5) + 1)
	a := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		a[i] = int64(rng.Intn(11) - 5)
	}
	input := fmt.Sprintf("%d %d\n", n, k)
	for i := 0; i <= n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", a[i])
	}
	input += "\n"
	exp := fmt.Sprintf("%d", solve(n, k, a))
	return input, exp
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
