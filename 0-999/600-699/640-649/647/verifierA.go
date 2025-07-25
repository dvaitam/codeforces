package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func solveA(grades []int) int {
	gifts := 0
	run := 0
	for _, g := range grades {
		if g >= 4 {
			run++
		} else {
			gifts += run / 3
			run = 0
		}
	}
	gifts += run / 3
	return gifts
}

func buildCaseA(n int, grades []int) (string, string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, g := range grades {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(g))
	}
	sb.WriteByte('\n')
	expect := fmt.Sprintf("%d", solveA(grades))
	return sb.String(), expect
}

func randomCaseA(rng *rand.Rand) (string, string) {
	n := rng.Intn(998) + 3 // 3..1000
	grades := make([]int, n)
	for i := 0; i < n; i++ {
		grades[i] = rng.Intn(5) + 1
	}
	return buildCaseA(n, grades)
}

func runCase(bin string, input string, expect string) error {
	out, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if out != expect {
		return fmt.Errorf("expected %q got %q", expect, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var tests [][2]string
	// deterministic edge cases
	in, exp := buildCaseA(3, []int{4, 4, 4})
	tests = append(tests, [2]string{in, exp})
	in, exp = buildCaseA(3, []int{1, 2, 3})
	tests = append(tests, [2]string{in, exp})
	in, exp = buildCaseA(6, []int{5, 5, 5, 4, 4, 4})
	tests = append(tests, [2]string{in, exp})
	in, exp = buildCaseA(5, []int{4, 4, 5, 1, 5})
	tests = append(tests, [2]string{in, exp})
	in, exp = buildCaseA(7, []int{4, 4, 5, 5, 5, 5, 5})
	tests = append(tests, [2]string{in, exp})

	for len(tests) < 100 {
		in, exp := randomCaseA(rng)
		tests = append(tests, [2]string{in, exp})
	}

	for i, tc := range tests {
		if err := runCase(bin, tc[0], tc[1]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
