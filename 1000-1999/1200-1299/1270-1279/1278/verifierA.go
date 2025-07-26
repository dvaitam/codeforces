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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func solveA(p, h string) string {
	if len(p) > len(h) {
		return "NO"
	}
	var fp [26]int
	for i := 0; i < len(p); i++ {
		fp[p[i]-'a']++
	}
	var fh [26]int
	l := len(p)
	for i := 0; i < l; i++ {
		fh[h[i]-'a']++
	}
	equal := func() bool {
		for i := 0; i < 26; i++ {
			if fp[i] != fh[i] {
				return false
			}
		}
		return true
	}
	if equal() {
		return "YES"
	}
	for i := l; i < len(h); i++ {
		fh[h[i]-'a']++
		fh[h[i-l]-'a']--
		if equal() {
			return "YES"
		}
	}
	return "NO"
}

func generateCase(r *rand.Rand) (string, string) {
	pLen := r.Intn(100) + 1
	hLen := r.Intn(100) + 1
	b := make([]byte, pLen)
	for i := 0; i < pLen; i++ {
		b[i] = byte('a' + r.Intn(26))
	}
	p := string(b)
	b = make([]byte, hLen)
	for i := 0; i < hLen; i++ {
		b[i] = byte('a' + r.Intn(26))
	}
	h := string(b)
	expect := solveA(p, h)
	input := fmt.Sprintf("1\n%s\n%s\n", p, h)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
