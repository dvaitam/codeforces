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

func enumerateSeq(m, length int, prev int, mask int, cnt []int) {
	if length == 0 {
		cnt[mask]++
		return
	}
	for c := 0; c < m; c++ {
		if c == prev {
			continue
		}
		enumerateSeq(m, length-1, c, mask|(1<<c), cnt)
	}
}

func layerWays(m int, l int) []int {
	size := 1 << m
	cnt := make([]int, size)
	enumerateSeq(m, l, -1, 0, cnt)
	return cnt
}

func solveE(n, m int, p int64, ls []int) string {
	ways := make([][]int, n)
	for i := 0; i < n; i++ {
		ways[i] = layerWays(m, ls[i])
	}
	size := 1 << m
	dp := make([]int64, size)
	for mask, v := range ways[0] {
		dp[mask] = int64(v) % p
	}
	for i := 1; i < n; i++ {
		ndp := make([]int64, size)
		for pm, pv := range dp {
			if pv == 0 {
				continue
			}
			for cm, cv := range ways[i] {
				if cm != pm {
					ndp[cm] = (ndp[cm] + pv*int64(cv)) % p
				}
			}
		}
		dp = ndp
	}
	var total int64
	for _, v := range dp {
		total = (total + v) % p
	}
	return fmt.Sprint(total)
}

func generateCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 2
	ls := make([]int, n)
	for i := 0; i < n; i++ {
		ls[i] = rng.Intn(3) + 1
	}
	p := int64(1000000007)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, p))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(ls[i]))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expect := solveE(n, m, p, ls)
	return input, expect
}

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

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCaseE(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
