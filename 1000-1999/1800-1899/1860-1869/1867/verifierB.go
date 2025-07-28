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
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, s string) string {
	mism := 0
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			mism++
		}
	}
	res := make([]byte, n+1)
	if n%2 == 1 {
		for i := mism; i <= n-mism; i++ {
			res[i] = '1'
		}
	} else {
		for i := mism; i <= n-mism; i += 2 {
			res[i] = '1'
		}
	}
	for i := 0; i <= n; i++ {
		if res[i] == 0 {
			res[i] = '0'
		}
	}
	return string(res)
}

func genCase(r *rand.Rand) (string, string) {
	n := r.Intn(20) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if r.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	s := string(b)
	input := fmt.Sprintf("1\n%d\n%s\n", n, s)
	return input, expected(n, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		in, exp := genCase(r)
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: %v\n", i, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
