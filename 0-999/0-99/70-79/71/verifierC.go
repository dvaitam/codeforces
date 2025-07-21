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

func compute(n int, a []int) string {
	for step := 1; step <= n; step++ {
		if n%step != 0 {
			continue
		}
		if n/step < 3 {
			continue
		}
		for offset := 0; offset < step; offset++ {
			ok := true
			for i := offset; i < n; i += step {
				if a[i] == 0 {
					ok = false
					break
				}
			}
			if ok {
				return "YES\n"
			}
		}
	}
	return "NO\n"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 3
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(2)
	}
	var in strings.Builder
	fmt.Fprintf(&in, "%d\n", n)
	for i, v := range a {
		if i > 0 {
			in.WriteByte(' ')
		}
		in.WriteString(fmt.Sprintf("%d", v))
	}
	in.WriteByte('\n')
	return in.String(), compute(n, a)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, buf.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
