package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const mod int64 = 1000000007

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func fib(n int) []int64 {
	f := make([]int64, n+2)
	f[0] = 1
	f[1] = 1
	for i := 2; i <= n; i++ {
		f[i] = (f[i-1] + f[i-2]) % mod
	}
	return f
}

func solve(s string) int64 {
	for _, c := range s {
		if c == 'w' || c == 'm' {
			return 0
		}
	}
	n := len(s)
	f := fib(n)
	res := int64(1)
	i := 0
	for i < n {
		if s[i] == 'u' || s[i] == 'n' {
			j := i
			for j < n && s[j] == s[i] {
				j++
			}
			length := j - i
			res = (res * f[length]) % mod
			i = j
		} else {
			i++
		}
	}
	return res
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := generateCase(rng)
		input := s + "\n"
		want := solve(s)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\noutput:%s", i+1, err, out)
			os.Exit(1)
		}
		if got%mod != want%mod {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
