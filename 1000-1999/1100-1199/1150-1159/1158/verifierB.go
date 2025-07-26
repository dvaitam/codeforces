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

func minUniqueLength(s string) int {
	n := len(s)
	for L := 1; L <= n; L++ {
		occ := make(map[string]int)
		for i := 0; i+L <= n; i++ {
			sub := s[i : i+L]
			occ[sub]++
		}
		for _, c := range occ {
			if c == 1 {
				return L
			}
		}
	}
	return n
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(10) + 1
		k := rng.Intn(n) + 1
		if n%2 != k%2 {
			if k%2 == 0 {
				n++
			} else {
				k++
			}
		}
		input := fmt.Sprintf("%d %d\n", n, k)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		s := strings.TrimSpace(out)
		if len(s) != n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected length %d got %d\ninput:%s", t+1, n, len(s), input)
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			if s[i] != '0' && s[i] != '1' {
				fmt.Fprintf(os.Stderr, "case %d failed: invalid character\ninput:%s", t+1, input)
				os.Exit(1)
			}
		}
		mu := minUniqueLength(s)
		if mu != k {
			fmt.Fprintf(os.Stderr, "case %d failed: expected minimal unique length %d got %d\ninput:%s output:%s", t+1, k, mu, input, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
