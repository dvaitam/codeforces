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
	n    int
	m    int
	grid []string
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(1194))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		grid := make([]string, n)
		for r := 0; r < n; r++ {
			var sb strings.Builder
			for c := 0; c < m; c++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('.')
				} else {
					sb.WriteByte('*')
				}
			}
			grid[r] = sb.String()
		}
		cases[i] = Case{n, m, grid}
	}
	return cases
}

func buildRef() (string, error) {
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1194B.go")
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
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", c.n, c.m)
	for _, row := range c.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	input := sb.String()
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
			var sb strings.Builder
			fmt.Fprintf(&sb, "1\n%d %d\n", c.n, c.m)
			for _, row := range c.grid {
				sb.WriteString(row)
				sb.WriteByte('\n')
			}
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
