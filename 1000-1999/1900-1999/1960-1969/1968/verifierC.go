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

func solve(n int, arr []int) []int {
	b := make([]int, n)
	b[0] = 501
	for i := 1; i < n; i++ {
		b[i] = b[i-1] + arr[i-1]
	}
	return b
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(9) + 2 // 2..10
	arr := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		arr[i] = rng.Intn(500) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprint(n))
	sb.WriteByte('\n')
	for i := 0; i < n-1; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(arr[i]))
	}
	sb.WriteByte('\n')
	b := solve(n, arr)
	exp := make([]string, n)
	for i := 0; i < n; i++ {
		exp[i] = strconv.Itoa(b[i])
	}
	return sb.String(), strings.Join(exp, " ")
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		inp, exp := genCase(rng)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", i+1, exp, got, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
