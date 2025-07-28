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

type Case struct{ a, b string }

func genSize(rng *rand.Rand) string {
	typ := rng.Intn(3)
	if typ == 1 { // M
		return "M"
	}
	cnt := rng.Intn(6)
	s := strings.Repeat("X", cnt)
	if typ == 0 {
		return s + "S"
	}
	return s + "L"
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]Case, 100)
	for i := range cases {
		cases[i] = Case{genSize(rng), genSize(rng)}
	}
	return cases
}

func expected(a, b string) string {
	rank := func(c byte) int {
		if c == 'S' {
			return 0
		}
		if c == 'M' {
			return 1
		}
		return 2
	}
	la, lb := len(a), len(b)
	ca, cb := a[la-1], b[lb-1]
	ra, rb := rank(ca), rank(cb)
	if ra != rb {
		if ra > rb {
			return ">"
		}
		return "<"
	}
	if ca == 'M' {
		return "="
	}
	if ca == 'S' {
		if la == lb {
			return "="
		}
		if la < lb {
			return ">"
		}
		return "<"
	}
	if la == lb {
		return "="
	}
	if la > lb {
		return ">"
	}
	return "<"
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, c Case) error {
	exp := expected(c.a, c.b)
	input := fmt.Sprintf("1\n%s %s\n", c.a, c.b)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != exp {
		return fmt.Errorf("expected %s got %s", exp, out)
	}
	return nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
