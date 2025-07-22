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

func expected(n, k int) string {
	if k > n || (k == 1 && n > 1) {
		return "-1"
	}
	if n == 1 {
		return "a"
	}
	res := make([]byte, n)
	altLen := n - (k - 2)
	for i := 0; i < altLen; i++ {
		if i%2 == 0 {
			res[i] = 'a'
		} else {
			res[i] = 'b'
		}
	}
	for j := 0; j < k-2; j++ {
		res[altLen+j] = byte('c' + j)
	}
	return string(res)
}

func runCase(bin string, n, k int) error {
	input := fmt.Sprintf("%d %d\n", n, k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expected(n, k)
	if got != exp {
		return fmt.Errorf("expected %q got %q (n=%d k=%d)", exp, got, n, k)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := []struct{ n, k int }{
		{1, 1},
		{2, 1},
		{5, 3},
		{3, 2},
		{4, 4},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 100 {
		n := rng.Intn(50) + 1
		k := rng.Intn(26) + 1
		cases = append(cases, struct{ n, k int }{n, k})
	}

	for i, c := range cases {
		if err := runCase(bin, c.n, c.k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
