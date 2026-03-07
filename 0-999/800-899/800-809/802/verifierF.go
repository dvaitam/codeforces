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

// oracle mirrors the C++ candidate's classification logic exactly.
func oracle(vals []int) string {
	var tag float64
	mx := 0
	for _, v := range vals {
		tag += float64(v * v)
		av := v
		if av < 0 {
			av = -av
		}
		if av > mx {
			mx = av
		}
	}
	if mx == 0 {
		// All zeros: sum of squares is 0, tag = 0 < 0.25 → poisson.
		return "poisson"
	}
	tag /= 250.0 * float64(mx) * float64(mx)
	if tag < 0.25 {
		return "poisson"
	}
	return "uniform"
}

func genTests(rng *rand.Rand) []string {
	tests := make([]string, 0, 102)

	// Generate random tests.
	for i := 0; i < 90; i++ {
		vals := make([]int, 250)
		for j := range vals {
			vals[j] = rng.Intn(5)
		}
		tests = append(tests, buildInput(vals))
	}

	// Explicit "uniform" case (all same non-zero value → tag = 1.0).
	tests = append(tests, buildInput(repeat(2, 250)))
	tests = append(tests, buildInput(repeat(3, 250)))

	// Explicit "poisson" case (sparse: mostly 0, a few non-zeros → tag << 0.25).
	sparse := make([]int, 250)
	sparse[5] = 1
	sparse[100] = 2
	tests = append(tests, buildInput(sparse))

	// Negative values (problem allows them).
	for i := 0; i < 7; i++ {
		vals := make([]int, 250)
		for j := range vals {
			vals[j] = rng.Intn(11) - 5 // [-5, 5]
		}
		tests = append(tests, buildInput(vals))
	}

	return tests
}

func repeat(v, n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = v
	}
	return s
}

func buildInput(vals []int) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	for j, v := range vals {
		if j > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runProg(bin, input string) (string, error) {
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
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseVals(input string) []int {
	r := strings.NewReader(input)
	var t int
	fmt.Fscan(r, &t)
	vals := make([]int, 250)
	for i := range vals {
		fmt.Fscan(r, &vals[i])
	}
	return vals
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := genTests(rng)
	for i, input := range tests {
		exp := oracle(parseVals(input))
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\ninput:%sexpected: %s got: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
