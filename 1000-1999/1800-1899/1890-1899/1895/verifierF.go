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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type state struct {
	pos  int
	val  int64
	seen bool
}

const MOD int64 = 1_000_000_007

func solveF(n int64, x, k int64) int64 {
	maxVal := x + k - 1 + k*(n-1)
	memo := make(map[state]int64)
	var dfs func(pos int, val int64, seen bool) int64
	dfs = func(pos int, val int64, seen bool) int64 {
		if val < 0 || val > maxVal {
			return 0
		}
		if pos == int(n) {
			if seen {
				return 1
			}
			return 0
		}
		st := state{pos, val, seen}
		if v, ok := memo[st]; ok {
			return v
		}
		var res int64
		for d := -k; d <= k; d++ {
			nv := val + d
			ns := seen || (nv >= x && nv <= x+k-1)
			res = (res + dfs(pos+1, nv, ns)) % MOD
		}
		memo[st] = res
		return res
	}
	var total int64
	for start := int64(0); start <= maxVal; start++ {
		total = (total + dfs(1, start, start >= x && start <= x+k-1)) % MOD
	}
	return total
}

func genCase(rng *rand.Rand) (int64, int64, int64) {
	n := int64(rng.Intn(4) + 1)
	x := int64(rng.Intn(5))
	k := int64(rng.Intn(3) + 1)
	return n, x, k
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, x, k := genCase(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, x, k))
		expect := fmt.Sprint(solveF(n, x, k))
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
