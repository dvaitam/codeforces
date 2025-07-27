package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct{ s string }

func genCases() []Case {
	rng := rand.New(rand.NewSource(1306))
	cases := make([]Case, 100)
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := range cases {
		l := rng.Intn(100) + 1
		b := make([]rune, l)
		for j := range b {
			b[j] = letters[rng.Intn(len(letters))]
		}
		cases[i] = Case{string(b)}
	}
	return cases
}

func expected(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func runCase(bin string, c Case) error {
	input := c.s + "\n"
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	want := expected(c.s)
	if strings.TrimSpace(out) != want {
		return fmt.Errorf("expected %s got %s", want, out)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
