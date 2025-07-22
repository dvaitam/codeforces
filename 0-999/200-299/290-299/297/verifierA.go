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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func canTransform(a, b string) bool {
	n, m := len(a), len(b)
	if a == b {
		return true
	}
	const B uint64 = 91138233
	pow := make([]uint64, m)
	if m > 0 {
		pow[0] = 1
		for i := 1; i < m; i++ {
			pow[i] = pow[i-1] * B
		}
	}
	var hb uint64
	for i := 0; i < m; i++ {
		if b[i] == '1' {
			hb += pow[m-1-i]
		}
	}
	for k := 0; k < n; k++ {
		s := make([]byte, m)
		if n-k >= m {
			copy(s, a[k:k+m])
		} else {
			copy(s, a[k:])
			s = s[:n-k]
			parity := byte(0)
			for _, c := range s {
				parity ^= c - '0'
			}
			for len(s) < m {
				s = append(s, '0'+parity)
				parity = 0
			}
		}
		var hs uint64
		var par byte
		for i := 0; i < m; i++ {
			if s[i] == '1' {
				hs += pow[m-1-i]
				par ^= 1
			}
		}
		for t := 0; t <= m; t++ {
			if hs == hb && string(s) == b {
				return true
			}
			if t == m {
				break
			}
			p := par
			front := s[0] - '0'
			hs = (hs - uint64(front)*pow[m-1]) * B
			if p == 1 {
				hs += 1
			}
			copy(s[0:m-1], s[1:m])
			s[m-1] = '0' + p
			par = front
		}
	}
	return false
}

func randBits(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		if rng.Intn(2) == 1 {
			b[i] = '1'
		} else {
			b[i] = '0'
		}
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		a := randBits(rng, n)
		b := randBits(rng, m)
		input := fmt.Sprintf("%s\n%s\n", a, b)
		expect := "NO"
		if canTransform(a, b) {
			expect = "YES"
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		res := strings.ToUpper(strings.TrimSpace(out))
		if res != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
