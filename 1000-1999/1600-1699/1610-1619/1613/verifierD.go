package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const mod int64 = 998244353

func run(bin, input string) (string, error) {
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

func solve(a []int) int64 {
	n := len(a)
	dp0 := make([]int64, n+5)
	dp1 := make([]int64, n+5)
	dp0[0] = 1
	for _, x := range a {
		dp0[x+1] = (dp0[x+1] * 2) % mod
		dp1[x+1] = (dp1[x+1] * 2) % mod
		tmp := dp0[x]
		dp0[x+1] = (dp0[x+1] + tmp) % mod
		if x >= 1 {
			tmp0 := dp0[x-1]
			tmp1 := dp1[x-1]
			dp1[x-1] = (tmp1*2 + tmp0) % mod
		}
	}
	var ans int64
	for i := 0; i < len(dp0); i++ {
		ans += dp0[i] + dp1[i]
	}
	ans = (ans - 1) % mod
	if ans < 0 {
		ans += mod
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(n + 1)
	}
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	input := fmt.Sprintf("1\n%d\n%s\n", n, sb.String())
	expect := fmt.Sprintf("%d", solve(a))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
