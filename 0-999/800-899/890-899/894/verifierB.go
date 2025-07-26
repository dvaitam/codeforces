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

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
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

const mod int64 = 1000000007
const modExp int64 = mod - 1

func powMod(base, exp int64) int64 {
	res := int64(1)
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func solve(n, m, k int64) string {
	if k == -1 && (n%2 != m%2) {
		return "0"
	}
	e := ((n - 1) % modExp * ((m - 1) % modExp)) % modExp
	return fmt.Sprintf("%d", powMod(2, e))
}

func genCase(rng *rand.Rand) string {
	n := rng.Int63n(1_000_000_000_000) + 1
	m := rng.Int63n(1_000_000_000_000) + 1
	k := int64(1)
	if rng.Intn(2) == 0 {
		k = -1
	}
	return fmt.Sprintf("%d %d %d\n", n, m, k)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := genCase(rng)
		var n, m, k int64
		fmt.Sscan(strings.TrimSpace(input), &n, &m, &k)
		expected := solve(n, m, k)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
