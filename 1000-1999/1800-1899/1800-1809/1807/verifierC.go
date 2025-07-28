package main

import (
	"bytes"
	"fmt"
	"math/rand"
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

type Case struct{ s string }

func genCases() []Case {
	rng := rand.New(rand.NewSource(1807))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(10) + 1
		b := make([]byte, n)
		for j := range b {
			b[j] = byte(rng.Intn(26) + 'a')
		}
		cases[i] = Case{s: string(b)}
	}
	return cases
}

func expected(s string) string {
	pos := make([]int, 26)
	for i := range pos {
		pos[i] = -1
	}
	for i, ch := range s {
		idx := ch - 'a'
		parity := i % 2
		if pos[idx] == -1 {
			pos[idx] = parity
		} else if pos[idx] != parity {
			return "NO"
		}
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		input := fmt.Sprintf("1\n%d\n%s\n", len(c.s), c.s)
		exp := expected(c.s)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(strings.ToUpper(got)) != exp {
			fmt.Printf("case %d failed: expected %s got %s (string %s)\n", i+1, exp, got, c.s)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
