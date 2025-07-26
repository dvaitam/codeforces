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

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveCase(a, b []byte) string {
	n := len(a)
	ans := 0
	for i := 1; i < n; i++ {
		if a[i-1] == b[i] && a[i] == b[i-1] && a[i-1] != a[i] {
			ans++
			a[i-1] = b[i-1]
			a[i] = b[i]
		}
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			ans++
		}
	}
	return fmt.Sprintf("%d", ans)
}

func genCase(r *rand.Rand) (string, string) {
	n := r.Intn(100) + 1
	a := make([]byte, n)
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if r.Intn(2) == 0 {
			a[i] = '0'
		} else {
			a[i] = '1'
		}
		if r.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	input := fmt.Sprintf("%d\n%s\n%s\n", n, a, b)
	expect := solveCase(append([]byte(nil), a...), append([]byte(nil), b...))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
