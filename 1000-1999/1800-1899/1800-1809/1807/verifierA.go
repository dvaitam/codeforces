package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runProg(bin, input string) (string, error) {
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
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct{ a, b, c int }

func genCases() []Case {
	cases := make([]Case, 0, 162)
	for a := 1; a <= 9; a++ {
		for b := 1; b <= 9; b++ {
			cases = append(cases, Case{a, b, a + b})
			cases = append(cases, Case{a, b, a - b})
		}
	}
	return cases
}

func expected(a, b, c int) string {
	if a+b == c {
		return "+"
	}
	return "-"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		input := fmt.Sprintf("1\n%d %d %d\n", c.a, c.b, c.c)
		exp := expected(c.a, c.b, c.c)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("case %d failed: expected %s got %s (input %d %d %d)\n", i+1, exp, got, c.a, c.b, c.c)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
