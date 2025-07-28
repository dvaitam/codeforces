package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func expected(n int, s, t string) string {
	if s[0] != t[0] || s[n-1] != t[n-1] {
		return "-1"
	}
	a := make([]int, 0)
	b := make([]int, 0)
	for i := 0; i < n-1; i++ {
		if s[i] != s[i+1] {
			a = append(a, i+1)
		}
		if t[i] != t[i+1] {
			b = append(b, i+1)
		}
	}
	if len(a) != len(b) {
		return "-1"
	}
	res := 0
	for i := range a {
		res += abs(a[i] - b[i])
	}
	return fmt.Sprint(res)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 3
	s := make([]byte, n)
	t := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			s[i] = '0'
		} else {
			s[i] = '1'
		}
		if rng.Intn(2) == 0 {
			t[i] = '0'
		} else {
			t[i] = '1'
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", n)
	sb.WriteString(string(s))
	sb.WriteByte('\n')
	sb.WriteString(string(t))
	sb.WriteByte('\n')
	exp := expected(n, string(s), string(t))
	return sb.String(), exp
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
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1705))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
