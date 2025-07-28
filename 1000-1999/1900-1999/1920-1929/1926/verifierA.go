package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(s string) string {
	cntA := 0
	for i := 0; i < len(s); i++ {
		if s[i] == 'A' {
			cntA++
		}
	}
	if cntA > len(s)-cntA {
		return "A"
	}
	return "B"
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	var all []string
	for i := 0; i < 32; i++ {
		var sb strings.Builder
		for j := 4; j >= 0; j-- {
			if i&(1<<j) != 0 {
				sb.WriteByte('A')
			} else {
				sb.WriteByte('B')
			}
		}
		all = append(all, sb.String())
	}
	var tests []string
	for len(tests) < 100 {
		tests = append(tests, all...)
	}
	tests = tests[:100]

	for idx, s := range tests {
		in := fmt.Sprintf("1\n%s\n", s)
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp := expected(s)
		if got != exp {
			fmt.Printf("test %d failed: input=%s expected=%s got=%s\n", idx+1, s, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
