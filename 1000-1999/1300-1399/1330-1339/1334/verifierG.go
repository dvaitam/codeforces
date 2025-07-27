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

func expected(p []int, s, t string) string {
	n := len(t)
	m := len(s)
	bitsets := make([]*big.Int, 26)
	for i := 0; i < 26; i++ {
		bitsets[i] = new(big.Int)
	}
	for i := 0; i < n; i++ {
		c := int(t[i] - 'a')
		bitsets[c].SetBit(bitsets[c], i, 1)
	}
	idx0 := int(s[0] - 'a')
	res := new(big.Int).Or(new(big.Int).Set(bitsets[idx0]), bitsets[p[idx0]])
	for k := 1; k < m; k++ {
		idx := int(s[k] - 'a')
		tmp := new(big.Int).Or(new(big.Int).Set(bitsets[idx]), bitsets[p[idx]])
		tmp.Rsh(tmp, uint(k))
		res.And(res, tmp)
	}
	limit := n - m + 1
	var sb strings.Builder
	for i := 0; i < limit; i++ {
		if res.Bit(i) == 1 {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
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

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		// permutation p
		perm := rng.Perm(26)
		sLen := rng.Intn(6) + 1
		tLen := rng.Intn(10) + sLen
		s := randString(rng, sLen)
		t := randString(rng, tLen)
		var sb strings.Builder
		for i, v := range perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v+1))
		}
		sb.WriteByte('\n')
		sb.WriteString(s)
		sb.WriteByte('\n')
		sb.WriteString(t)
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(perm, s, t)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
