package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	in  string
	out string
}

func solve(n int, f []int) string {
	for i := 1; i <= n; i++ {
		j := f[i]
		if j >= 1 && j <= n {
			k := f[j]
			if k >= 1 && k <= n && f[k] == i {
				return "YES"
			}
		}
	}
	return "NO"
}

func generateTests() []test {
	rand.Seed(1)
	tests := make([]test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 2 // 2..9
		f := make([]int, n+1)
		for j := 1; j <= n; j++ {
			x := rand.Intn(n) + 1
			if x == j {
				x = (x % n) + 1
			}
			f[j] = x
		}
		if i%2 == 0 && n >= 3 {
			f[1] = 2
			f[2] = 3
			f[3] = 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 1; j <= n; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", f[j])
		}
		sb.WriteByte('\n')
		tests = append(tests, test{in: sb.String(), out: solve(n, f)})
	}
	return tests
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := run(bin, t.in)
		if err != nil {
			fmt.Printf("Test %d failed to run: %v\n", i+1, err)
			fmt.Print(out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != t.out {
			fmt.Printf("Test %d failed. Expected %q, got %q. Input:\n%s", i+1, t.out, out, t.in)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(tests))
}
