package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Case struct {
	s string
	t string
	p string
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(1194))
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	cases := make([]Case, 100)
	for i := range cases {
		ls := rng.Intn(5) + 1
		lt := ls + rng.Intn(5)
		lp := rng.Intn(5) + 1
		b := make([]rune, ls)
		for j := range b {
			b[j] = letters[rng.Intn(len(letters))]
		}
		s := string(b)
		b = make([]rune, lt)
		for j := range b {
			b[j] = letters[rng.Intn(len(letters))]
		}
		t := string(b)
		b = make([]rune, lp)
		for j := range b {
			b[j] = letters[rng.Intn(len(letters))]
		}
		p := string(b)
		cases[i] = Case{s, t, p}
	}
	return cases
}

func buildRef() (string, error) {
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1194C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin, ref string, c Case) error {
	input := fmt.Sprintf("1\n%s\n%s\n%s\n", c.s, c.t, c.p)
	expected, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if expected != strings.TrimSpace(got) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, ref, c); err != nil {
			input := fmt.Sprintf("1\n%s\n%s\n%s\n", c.s, c.t, c.p)
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
