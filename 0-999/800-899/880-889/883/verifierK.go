package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input string
	n     int
	s     []int
	g     []int
}

type solution struct {
	possible bool
	total    int64
	roads    []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Printf("Reference runtime error on test %d: %v\nInput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		refSol, err := parseSolution(refOut, tc.n)
		if err != nil {
			fmt.Printf("Failed to parse reference output on test %d: %v\nOutput:\n%s\n", i+1, err, refOut)
			os.Exit(1)
		}
		if refSol.possible {
			if err := validateFeasible(tc, refSol); err != nil {
				fmt.Printf("Reference produced invalid solution on test %d: %v\n", i+1, err)
				os.Exit(1)
			}
		}

		out, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\nInput:\n%sOutput:\n%s\n", i+1, err, tc.input, out)
			os.Exit(1)
		}
		cand, err := parseSolution(out, tc.n)
		if err != nil {
			fmt.Printf("Failed to parse output on test %d: %v\nInput:\n%sOutput:\n%s\n", i+1, err, tc.input, out)
			os.Exit(1)
		}
		if err := compareWithReference(tc, refSol, cand); err != nil {
			fmt.Printf("Wrong answer on test %d: %v\nInput:\n%sOutput:\n%s\n", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	path := "./ref883K.bin"
	cmd := exec.Command("go", "build", "-o", path, "883K.go")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	return path, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func parseSolution(out string, n int) (solution, error) {
	var sol solution
	reader := bufio.NewReader(strings.NewReader(out))
	var firstStr string
	if _, err := fmt.Fscan(reader, &firstStr); err != nil {
		return sol, fmt.Errorf("failed to read first token: %v", err)
	}
	if firstStr == "-1" {
		rest, _ := io.ReadAll(reader)
		if strings.TrimSpace(string(rest)) != "" {
			return sol, fmt.Errorf("extra output after -1")
		}
		return sol, nil
	}
	var total int64
	if _, err := fmt.Sscan(firstStr, &total); err != nil {
		return sol, fmt.Errorf("first token is neither -1 nor integer total")
	}
	sol.possible = true
	sol.total = total
	sol.roads = make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &sol.roads[i]); err != nil {
			return sol, fmt.Errorf("failed to read road width %d: %v", i+1, err)
		}
	}
	if extra := strings.TrimSpace(readAllRemaining(reader)); extra != "" {
		return sol, fmt.Errorf("unexpected extra tokens: %q", extra)
	}
	return sol, nil
}

func readAllRemaining(r *bufio.Reader) string {
	var sb strings.Builder
	for {
		chunk, err := r.ReadString('\n')
		sb.WriteString(chunk)
		if err != nil {
			break
		}
	}
	return sb.String()
}

func compareWithReference(tc testCase, ref, cand solution) error {
	if !ref.possible {
		if cand.possible {
			return fmt.Errorf("should output -1 but provided a solution")
		}
		return nil
	}
	if !cand.possible {
		return fmt.Errorf("reported -1 but solution exists")
	}
	if err := validateFeasible(tc, cand); err != nil {
		return err
	}
	if cand.total != ref.total {
		return fmt.Errorf("total demolished %d differs from optimal %d", cand.total, ref.total)
	}
	return nil
}

func validateFeasible(tc testCase, sol solution) error {
	if len(sol.roads) != tc.n {
		return fmt.Errorf("expected %d widths, got %d", tc.n, len(sol.roads))
	}
	var destroyed int64
	for i, w := range sol.roads {
		lo := tc.s[i]
		hi := tc.s[i] + tc.g[i]
		if w < lo || w > hi {
			return fmt.Errorf("road width %d at position %d out of range [%d,%d]", w, i+1, lo, hi)
		}
		if i > 0 && abs(w-sol.roads[i-1]) > 1 {
			return fmt.Errorf("adjacent widths %d and %d differ by more than 1 at positions %d,%d", sol.roads[i-1], w, i, i+1)
		}
		destroyed += int64(w - tc.s[i])
	}
	if destroyed != sol.total {
		return fmt.Errorf("reported total %d but actual demolished %d", sol.total, destroyed)
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	addCase := func(sVals, gVals []int) {
		cpS := append([]int(nil), sVals...)
		cpG := append([]int(nil), gVals...)
		tests = append(tests, testCase{
			input: formatInput(cpS, cpG),
			n:     len(cpS),
			s:     cpS,
			g:     cpG,
		})
	}

	addCase([]int{5}, []int{0})
	addCase([]int{3}, []int{10})
	addCase([]int{1, 100}, []int{0, 0})
	addCase([]int{4, 4, 4}, []int{3, 0, 3})
	addCase([]int{10, 1, 10, 1}, []int{0, 0, 0, 0})
	addCase([]int{2, 2, 2, 2, 2}, []int{5, 5, 5, 5, 5})
	addCase([]int{1, 2, 3, 4, 5}, []int{0, 0, 0, 0, 0})

	rng := rand.New(rand.NewSource(123456789))
	for len(tests) < 120 {
		n := rng.Intn(8) + 2
		s := make([]int, n)
		g := make([]int, n)
		base := rng.Intn(5)
		for i := 0; i < n; i++ {
			s[i] = rng.Intn(15) + base
			if rng.Intn(4) == 0 {
				g[i] = rng.Intn(20)
			} else {
				g[i] = rng.Intn(6)
			}
		}
		addCase(s, g)
	}

	largeN := 500
	sLarge := make([]int, largeN)
	gLarge := make([]int, largeN)
	for i := 0; i < largeN; i++ {
		sLarge[i] = (i % 5) * 2
		gLarge[i] = 1000
	}
	addCase(sLarge, gLarge)

	return tests
}

func formatInput(s, g []int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(s))
	for i := range s {
		fmt.Fprintf(&sb, "%d %d\n", s[i], g[i])
	}
	return sb.String()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
