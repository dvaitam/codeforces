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

func expected(k int) string {
	if k%2 == 0 {
		return "NO"
	}
	n := 4*k - 2
	m := 2*((k-1)+(k-1)*(k-1)+((k-1)/2)) + 1
	var sb strings.Builder
	sb.WriteString("YES\n")
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	goFunc := func(start, limit, n int) {
		for i := start; i < limit; i += 2 {
			fmt.Fprintf(&sb, "%d %d\n", i, i+1)
			for j := 0; j < n-1; j++ {
				fmt.Fprintf(&sb, "%d %d\n", i, limit+j)
				fmt.Fprintf(&sb, "%d %d\n", i+1, limit+j)
			}
		}
		for i := limit; i < limit+n-1; i++ {
			fmt.Fprintf(&sb, "%d %d\n", i, limit+n-1)
		}
	}
	goFunc(1, k, k)
	goFunc(2*k, 3*k-1, k)
	fmt.Fprintf(&sb, "%d %d\n", 2*k-1, 4*k-2)
	return strings.TrimRight(sb.String(), "\n")
}

func genCase(rng *rand.Rand) (string, string) {
	k := rng.Intn(10) + 1
	input := fmt.Sprintf("%d\n", k)
	exp := expected(k)
	return input, exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected \n%s\ngot \n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
