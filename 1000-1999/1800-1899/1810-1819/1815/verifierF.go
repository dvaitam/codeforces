package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Triangle struct{ a, b, c int }

type Case struct {
	n       int
	m       int
	weights []int
	tris    []Triangle
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(1815))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(5) + 3 // 3..7
		m := rng.Intn(4) + 1 // 1..4
		weights := make([]int, n)
		for j := range weights {
			weights[j] = rng.Intn(10)
		}
		tris := make([]Triangle, m)
		for j := range tris {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			for b == a {
				b = rng.Intn(n) + 1
			}
			c := rng.Intn(n) + 1
			for c == a || c == b {
				c = rng.Intn(n) + 1
			}
			if a > b {
				a, b = b, a
			}
			if b > c {
				b, c = c, b
			}
			if a > b {
				a, b = b, a
			}
			tris[j] = Triangle{a, b, c}
		}
		cases[i] = Case{n, m, weights, tris}
	}
	return cases
}

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1815F.go")
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
	for i, w := range c.weights {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", w)
	}
	input += "\n"
	for _, t := range c.tris {
		input += fmt.Sprintf("%d %d %d\n", t.a, t.b, t.c)
	}
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
			for idx, w := range c.weights {
				if idx > 0 {
					input += " "
				}
				input += fmt.Sprintf("%d", w)
			}
			input += "\n"
			for _, t := range c.tris {
				input += fmt.Sprintf("%d %d %d\n", t.a, t.b, t.c)
			}
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
