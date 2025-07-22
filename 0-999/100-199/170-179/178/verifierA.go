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
	n   int
	arr []int
}

func genTests() []test {
	rand.Seed(42)
	var tests []test
	for i := 0; i < 100; i++ {
		n := rand.Intn(9) + 2
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(20)
		}
		tests = append(tests, test{n, arr})
	}
	return tests
}

func expected(t test) string {
	sum := 0
	var sb strings.Builder
	for i, v := range t.arr {
		sum += v
		if i < t.n-1 {
			sb.WriteString(fmt.Sprintf("%d\n", sum))
		}
	}
	return strings.TrimSpace(sb.String())
}

func inputString(t test) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.n))
	for i, v := range t.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func run(bin string, in string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		in := inputString(t)
		exp := expected(t)
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\nactual:\n%s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
