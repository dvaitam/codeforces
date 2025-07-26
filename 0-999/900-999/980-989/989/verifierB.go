package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	in  string
	out string
}

func expected(n, p int, s string) string {
	a := []byte(s)
	possible := false
	idx := -1
	for i := 0; i < n-p; i++ {
		x, y := a[i], a[i+p]
		if x == '.' || y == '.' || x != y {
			possible = true
			idx = i
			break
		}
	}
	if !possible {
		return "No"
	}
	x, y := a[idx], a[idx+p]
	if x == '.' && y == '.' {
		a[idx] = '0'
		a[idx+p] = '1'
	} else if x == '.' {
		if y == '0' {
			a[idx] = '1'
		} else {
			a[idx] = '0'
		}
	} else if y == '.' {
		if x == '0' {
			a[idx+p] = '1'
		} else {
			a[idx+p] = '0'
		}
	}
	for i := 0; i < n; i++ {
		if a[i] == '.' {
			a[i] = '0'
		}
	}
	return string(a)
}

func generate() []testCase {
	const T = 100
	rand.Seed(2)
	cases := make([]testCase, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(20) + 1
		p := rand.Intn(n) + 1
		var s string
		for {
			var sb strings.Builder
			for j := 0; j < n; j++ {
				sb.WriteByte("01."[rand.Intn(3)])
			}
			s = sb.String()
			if strings.ContainsRune(s, '.') {
				break
			}
		}
		cases[i] = testCase{
			in:  fmt.Sprintf("%d %d\n%s\n", n, p, s),
			out: expected(n, p, s),
		}
	}
	return cases
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
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generate()
	for idx, tc := range cases {
		got, err := run(bin, tc.in)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(tc.out) {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, tc.out, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
