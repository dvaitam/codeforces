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

func solve(arr []int) int {
	res := 0
	for _, x := range arr {
		for _, y := range arr {
			if x >= y {
				if m := x % y; m > res {
					res = m
				}
			}
		}
	}
	return res
}

func generateTests() []test {
	rnd := rand.New(rand.NewSource(2))
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		n := rnd.Intn(50) + 1
		arr := make([]int, n)
		var in bytes.Buffer
		fmt.Fprintln(&in, n)
		for i := 0; i < n; i++ {
			v := rnd.Intn(1000) + 1
			arr[i] = v
			if i > 0 {
				fmt.Fprint(&in, " ")
			}
			fmt.Fprint(&in, v)
		}
		in.WriteByte('\n')
		out := fmt.Sprintf("%d\n", solve(arr))
		tests = append(tests, test{in: in.String(), out: out})
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		g := strings.TrimSpace(got)
		e := strings.TrimSpace(t.out)
		if g != e {
			fmt.Fprintf(os.Stderr, "Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, t.in, e, g)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
