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
	input string
}

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "845F.go")
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(8456))
	cases := make([]Case, 100)
	for i := range cases {
		var sb strings.Builder
		var n, m int
		for {
			n = rng.Intn(15) + 1
			m = rng.Intn(15) + 1
			if n*m <= 250 {
				break
			}
		}
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for r := 0; r < n; r++ {
			for c := 0; c < m; c++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('.')
				} else {
					sb.WriteByte('x')
				}
			}
			sb.WriteByte('\n')
		}
		cases[i] = Case{sb.String()}
	}
	return cases
}

func runCase(bin, ref string, c Case) error {
	expected, err := runBinary(ref, c.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(expected) != strings.TrimSpace(got) {
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
