package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const maxN = 200000

var prefix [maxN + 1]int64
var sumDigits [maxN + 1]int

func init() {
	for i := 1; i <= maxN; i++ {
		sumDigits[i] = sumDigits[i/10] + i%10
		prefix[i] = prefix[i-1] + int64(sumDigits[i])
	}
}

func expected(n int) string {
	if n > maxN {
		n = maxN
	}
	return fmt.Sprint(prefix[n])
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	var tests []int
	for i := 1; i <= 99; i++ {
		n := i * 2000
		if n > maxN {
			n = maxN
		}
		tests = append(tests, n)
	}
	tests = append(tests, maxN)

	for idx, n := range tests {
		in := fmt.Sprintf("1\n%d\n", n)
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp := expected(n)
		if got != exp {
			fmt.Printf("test %d failed: n=%d expected=%s got=%s\n", idx+1, n, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
