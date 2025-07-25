package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

type Case struct {
	input string
	n     int
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(712))
	cases := make([]Case, 0, 100)
	// add some fixed cases
	cases = append(cases, Case{input: "1\n", n: 1})
	cases = append(cases, Case{input: "3\n", n: 3})
	for len(cases) < 100 {
		n := rng.Intn(25)*2 + 1 // odd between 1 and 49
		cases = append(cases, Case{input: fmt.Sprintf("%d\n", n), n: n})
	}
	return cases
}

func verify(n int, out string) error {
	tokens := strings.Fields(out)
	if len(tokens) != n*n {
		return fmt.Errorf("expected %d numbers got %d", n*n, len(tokens))
	}
	seen := make([]bool, n*n+1)
	matrix := make([][]int, n)
	idx := 0
	for i := 0; i < n; i++ {
		row := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(tokens[idx])
			if err != nil {
				return fmt.Errorf("invalid number %q", tokens[idx])
			}
			if v < 1 || v > n*n {
				return fmt.Errorf("value %d out of range", v)
			}
			if seen[v] {
				return fmt.Errorf("duplicate value %d", v)
			}
			seen[v] = true
			row[j] = v
			idx++
		}
		matrix[i] = row
	}
	// check sums
	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < n; j++ {
			sum += matrix[i][j]
		}
		if sum%2 == 0 {
			return fmt.Errorf("row %d sum not odd", i+1)
		}
	}
	for j := 0; j < n; j++ {
		sum := 0
		for i := 0; i < n; i++ {
			sum += matrix[i][j]
		}
		if sum%2 == 0 {
			return fmt.Errorf("column %d sum not odd", j+1)
		}
	}
	sum := 0
	for i := 0; i < n; i++ {
		sum += matrix[i][i]
	}
	if sum%2 == 0 {
		return fmt.Errorf("main diagonal sum not odd")
	}
	sum = 0
	for i := 0; i < n; i++ {
		sum += matrix[i][n-i-1]
	}
	if sum%2 == 0 {
		return fmt.Errorf("second diagonal sum not odd")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		out, err := runBinary(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		if err := verify(c.n, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, c.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
