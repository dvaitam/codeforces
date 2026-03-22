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

const mod = 1000000007

func solveAdd(a, b int) int {
	res := a + b
	if res >= mod {
		res -= mod
	}
	return res
}

func solveSub(a, b int) int {
	res := a - b
	if res < 0 {
		res += mod
	}
	return res
}

func solveMul(a, b int) int {
	return int((int64(a) * int64(b)) % mod)
}

func solvePower(base, exp int) int {
	res := 1
	base %= mod
	if base < 0 {
		base += mod
	}
	for exp > 0 {
		if exp%2 == 1 {
			res = solveMul(res, base)
		}
		base = solveMul(base, base)
		exp /= 2
	}
	return res
}

func solveInv(n int) int {
	return solvePower(n, mod-2)
}

func solve(input string) string {
	var n, k, m int
	fmt.Sscan(input, &n, &k, &m)

	if k == 1 || k == n {
		return fmt.Sprintf("%d", n)
	}

	maxA, maxB, vecSize := m+1, m+1, m
	pool := make([]int, maxA*maxB*vecSize)
	E := make([][][]int, maxA)
	idx := 0
	for i := 0; i < maxA; i++ {
		E[i] = make([][]int, maxB)
		for j := 0; j < maxB; j++ {
			E[i][j] = pool[idx : idx+vecSize]
			idx += vecSize
		}
	}

	sumA := make([][]int, maxB)
	for i := 0; i < maxB; i++ {
		sumA[i] = make([]int, vecSize)
	}

	sumB := make([][]int, maxB)
	for i := 0; i < maxB; i++ {
		sumB[i] = make([]int, vecSize)
	}

	Eq := make([][]int, maxA)
	for i := 0; i < maxA; i++ {
		Eq[i] = make([]int, vecSize)
	}

	for b := 0; b <= m-1; b++ {
		E[0][b][0] = b + 1
		copy(sumA[b], E[0][b])
	}
	for a := 1; a <= m-1; a++ {
		E[a][0][0] = a + 1
	}
	for b := 1; b <= m-2; b++ {
		E[1][b][b] = 1
	}

	for a := 1; a <= m-2; a++ {
		copy(sumB[1], E[a][0])
		for b := 1; b <= m-a-1; b++ {
			for i := 0; i < vecSize; i++ {
				sumB[b+1][i] = solveAdd(sumB[b][i], E[a][b][i])
			}
		}

		bEq := m - 1 - a
		invM1 := solveInv(m - 1)
		for i := 0; i < vecSize; i++ {
			sumAB := solveAdd(sumA[bEq][i], sumB[bEq][i])
			term := solveMul(sumAB, invM1)
			Eq[a][i] = solveSub(E[a][bEq][i], term)
		}

		if a <= m-3 {
			for b := 1; b <= m-a-2; b++ {
				l := a + b + 1
				PC := solveMul(m-l, solveInv(m))
				PB := solveMul(l, solveInv(m))

				coeff_Eab1 := solveMul(solveSub(0, PC), solveMul(b+1, solveInv(l+1)))
				coeff_sum := solveMul(solveSub(0, PB), solveInv(l-1))
				multiplier := solveMul(l+1, solveInv(solveMul(PC, a+1)))

				for i := 0; i < vecSize; i++ {
					val := E[a][b][i]
					val = solveAdd(val, solveMul(E[a][b+1][i], coeff_Eab1))
					sumAB := solveAdd(sumA[b][i], sumB[b][i])
					val = solveAdd(val, solveMul(sumAB, coeff_sum))
					E[a+1][b][i] = solveMul(val, multiplier)
				}
			}
		}

		for b := 1; b <= m-a-1; b++ {
			for i := 0; i < vecSize; i++ {
				sumA[b][i] = solveAdd(sumA[b][i], E[a][b][i])
			}
		}
	}

	nVars := m - 2
	Matrix := make([][]int, nVars)
	for i := 0; i < nVars; i++ {
		Matrix[i] = make([]int, nVars+1)
	}

	for i := 0; i < nVars; i++ {
		a := i + 1
		for j := 0; j < nVars; j++ {
			Matrix[i][j] = Eq[a][j+1]
		}
		Matrix[i][nVars] = solveSub(0, Eq[a][0])
	}

	for i := 0; i < nVars; i++ {
		pivot := i
		for j := i; j < nVars; j++ {
			if Matrix[j][i] != 0 {
				pivot = j
				break
			}
		}
		Matrix[i], Matrix[pivot] = Matrix[pivot], Matrix[i]

		invPivot := solveInv(Matrix[i][i])
		for j := i; j <= nVars; j++ {
			Matrix[i][j] = solveMul(Matrix[i][j], invPivot)
		}

		for j := 0; j < nVars; j++ {
			if i != j && Matrix[j][i] != 0 {
				factor := Matrix[j][i]
				for kk := i; kk <= nVars; kk++ {
					Matrix[j][kk] = solveSub(Matrix[j][kk], solveMul(factor, Matrix[i][kk]))
				}
			}
		}
	}

	U := make([]int, vecSize)
	U[0] = 1
	for i := 1; i <= nVars; i++ {
		U[i] = Matrix[i-1][nVars]
	}

	aAns := k - 1
	bAns := n - k
	ans := 0
	for i := 0; i < vecSize; i++ {
		ans = solveAdd(ans, solveMul(E[aAns][bAns][i], U[i]))
	}
	return fmt.Sprintf("%d", ans)
}

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	for idx, tc := range tests {
		refOut := solve(tc.input)
		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if refVal != candVal {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, refVal, candVal, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(output string) (int64, error) {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return 0, fmt.Errorf("empty output")
	}
	tokens := strings.Fields(trimmed)
	if len(tokens) != 1 {
		return 0, fmt.Errorf("expected single integer, got %q", trimmed)
	}
	val, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", tokens[0])
	}
	return val % int64(mod), nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeManualTest("minimal", 1, 1, 1),
		makeManualTest("already_end_left", 5, 1, 5),
		makeManualTest("already_end_right", 6, 6, 6),
		makeManualTest("middle_no_growth", 2, 1, 2),
		makeManualTest("middle_can_grow", 3, 2, 5),
		makeManualTest("large_m", 7, 3, 250),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func makeManualTest(name string, n, k, m int) testCase {
	return testCase{
		name:  name,
		input: fmt.Sprintf("%d %d %d\n", n, k, m),
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(250) + 1
	m := n + rng.Intn(251-n)
	k := rng.Intn(n) + 1
	name := fmt.Sprintf("random_%d_n%d_k%d_m%d", idx+1, n, k, m)
	return testCase{
		name:  name,
		input: fmt.Sprintf("%d %d %d\n", n, k, m),
	}
}
