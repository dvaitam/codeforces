package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "981E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func formatCase(n int, ops [][3]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(ops)))
	for _, op := range ops {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", op[0], op[1], op[2]))
	}
	return sb.String()
}

func addHandcraftedTests(cases *[]string) {
	*cases = append(*cases,
		formatCase(1, [][3]int{{1, 1, 1}}),
		formatCase(5, [][3]int{{1, 5, 3}}),
		formatCase(5, [][3]int{{1, 5, 2}, {2, 5, 2}, {1, 3, 1}}),
		formatCase(6, [][3]int{{1, 2, 1}, {2, 3, 1}, {3, 4, 1}, {4, 5, 1}, {5, 6, 1}}),
	)
}

func addRandomTests(cases *[]string) {
	rng := rand.New(rand.NewSource(981))
	for t := 0; t < 300; t++ {
		n := rng.Intn(20) + 1
		q := rng.Intn(35) + 1
		ops := make([][3]int, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			v := rng.Intn(n) + 1
			ops[i] = [3]int{l, r, v}
		}
		*cases = append(*cases, formatCase(n, ops))
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	cases := make([]string, 0, 304)
	addHandcraftedTests(&cases)
	addRandomTests(&cases)

	for idx, input := range cases {
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", idx+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d generated tests passed\n", len(cases))
}
