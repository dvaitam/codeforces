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

func expected(input string) string {
	// compute expected output for input of problem A
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var t int
	fmt.Sscan(lines[0], &t)
	idx := 1
	out := make([]string, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		var n int
		fmt.Sscan(lines[idx], &n)
		idx++
		arrStr := strings.Fields(lines[idx])
		idx++
		mx := int64(0)
		for j := 0; j < n; j++ {
			var v int64
			fmt.Sscan(arrStr[j], &v)
			if j%2 == 0 && v > mx {
				mx = v
			}
		}
		out[caseNum] = fmt.Sprintf("%d", mx)
	}
	return strings.Join(out, "\n") + "\n"
}

func genCase(rng *rand.Rand) (string, string) {
	t := 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	n := rng.Intn(49)*2 + 1 // odd up to 99
	sb.WriteString(fmt.Sprintf("%d\n", n))
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Int63n(100) + 1
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	input := sb.String()
	return input, expected(input)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	// Predefined edge cases
	cases := []string{
		"1\n1\n6\n",
		"1\n3\n1 3 2\n",
		"1\n5\n4 7 4 2 9\n",
	}
	exps := []string{
		expected(cases[0]),
		expected(cases[1]),
		expected(cases[2]),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 100 {
		in, exp := genCase(rng)
		cases = append(cases, in)
		exps = append(exps, exp)
	}

	for i := range cases {
		if err := runCase(bin, cases[i], exps[i]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, cases[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
