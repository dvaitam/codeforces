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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func minOps(s string) int {
	n := len(s)
	prefixOps := make([]int, n+1)
	prefixPar := make([]int, n+1)
	parity := 0
	op := 0
	for i := 1; i <= n; i++ {
		bit := int(s[i-1]-'0') ^ parity
		if bit == 1 {
			op++
			parity ^= 1
		}
		prefixOps[i] = op
		prefixPar[i] = parity
	}
	suffix0 := make([]int, n+2)
	suffix1 := make([]int, n+2)
	op = 0
	parity = 0
	for i := n; i >= 1; i-- {
		bit := int(s[i-1]-'0') ^ parity
		if bit == 0 {
			op++
			parity ^= 1
		}
		suffix0[i] = op
	}
	op = 0
	parity = 1
	for i := n; i >= 1; i-- {
		bit := int(s[i-1]-'0') ^ parity
		if bit == 0 {
			op++
			parity ^= 1
		}
		suffix1[i] = op
	}
	best := n + 5
	for k := 0; k <= n; k++ {
		p := prefixPar[k]
		ops := prefixOps[k]
		if p == 0 {
			ops += suffix0[k+1]
		} else {
			ops += suffix1[k+1]
		}
		if ops < best {
			best = ops
		}
	}
	return best
}

func solve(reader *bufio.Reader) string {
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return ""
	}
	var sb strings.Builder
	for ; T > 0; T-- {
		var n int
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)
		fmt.Fprintf(&sb, "%d\n", minOps(s))
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) string {
	T := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", T)
	for ; T > 0; T-- {
		n := rng.Intn(30) + 1
		fmt.Fprintf(&sb, "%d\n", n)
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				b[i] = '0'
			} else {
				b[i] = '1'
			}
		}
		sb.Write(b)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solve(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
