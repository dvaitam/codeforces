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
	n uint64
	m int64
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(1815))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Uint64() % 1000000000 // up to 1e9
		m := int64(rng.Intn(10) + 1)
		cases[i] = Case{n, m}
	}
	return cases
}

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1815D.go")
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
	input := fmt.Sprintf("1\n%d %d\n", c.n, c.m)
	expected, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
			input := fmt.Sprintf("1\n%d %d\n", c.n, c.m)
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
