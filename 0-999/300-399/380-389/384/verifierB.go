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

func expected(n, m, k int) string {
	total := m * (m - 1) / 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", total))
	if k == 0 {
		for i := 1; i <= m; i++ {
			for j := i + 1; j <= m; j++ {
				sb.WriteString(fmt.Sprintf("%d %d\n", i, j))
			}
		}
	} else {
		for i := 1; i <= m; i++ {
			for j := i + 1; j <= m; j++ {
				sb.WriteString(fmt.Sprintf("%d %d\n", j, i))
			}
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(8) + 1
	k := rng.Intn(2)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rng.Intn(1000000) + 1
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	input := sb.String()
	return input, expected(n, m, k)
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// deterministic edge cases
	cases := []struct{ n, m, k int }{
		{1, 1, 0}, {1, 2, 0}, {1, 2, 1}, {3, 3, 0}, {2, 5, 1},
	}
	for i, tc := range cases {
		input := fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k)
		for l := 0; l < tc.n; l++ {
			for j := 0; j < tc.m; j++ {
				input += "1"
				if j+1 < tc.m {
					input += " "
				}
			}
			if l+1 < tc.n {
				input += "\n"
			}
		}
		if err := runCase(bin, input, expected(tc.n, tc.m, tc.k)); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	for i := 0; i < 100-len(cases); i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
