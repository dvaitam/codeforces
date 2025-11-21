package main

import (
	"bytes"
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceA2 = "1000-1999/1000-1099/1000-1009/1002/1002A2.go"
	refBinaryA2 = "ref1002A2.bin"
	totalTests  = 80
)

type testCase struct {
	n    int
	bits string
}

type operation struct {
	kind string
	a    int
	b    int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, cleanup, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		if err := verifyOutput(tc, refOut); err != nil {
			fmt.Printf("reference output invalid on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		if err := verifyOutput(tc, candOut); err != nil {
			fmt.Printf("candidate output invalid on test %d: %v\n", idx+1, err)
			printInput(input)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref-1002A2-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1002A2.bin")
	cmd := exec.Command("go", "build", "-o", bin, refSourceA2)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("build error: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return bin, cleanup, nil
}

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func verifyOutput(tc testCase, out string) error {
	ops, err := parseOperations(out)
	if err != nil {
		return err
	}
	state, err := simulate(tc.n, ops)
	if err != nil {
		return err
	}
	exp := expectedState(tc)
	if err := compareStates(state, exp); err != nil {
		return err
	}
	return nil
}

func parseOperations(out string) ([]operation, error) {
	lines := make([]string, 0)
	for _, raw := range strings.Split(out, "\n") {
		line := strings.TrimSpace(raw)
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	start := 0
	expectCount := -1
	if len(lines) > 0 {
		if cnt, err := strconv.Atoi(lines[0]); err == nil {
			expectCount = cnt
			start = 1
		}
	}
	ops := make([]operation, 0, len(lines)-start)
	for i := start; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		if len(fields) == 0 {
			continue
		}
		cmd := strings.ToUpper(fields[0])
		switch cmd {
		case "H":
			if len(fields) != 2 {
				return nil, fmt.Errorf("invalid H command format")
			}
			q, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, fmt.Errorf("invalid qubit index in H command")
			}
			ops = append(ops, operation{kind: "H", a: q})
		case "CNOT":
			if len(fields) != 3 {
				return nil, fmt.Errorf("invalid CNOT command format")
			}
			c, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, fmt.Errorf("invalid control qubit in CNOT")
			}
			t, err := strconv.Atoi(fields[2])
			if err != nil {
				return nil, fmt.Errorf("invalid target qubit in CNOT")
			}
			ops = append(ops, operation{kind: "CNOT", a: c, b: t})
		default:
			return nil, fmt.Errorf("unsupported command %q", cmd)
		}
	}
	if expectCount >= 0 && expectCount != len(ops) {
		return nil, fmt.Errorf("operation count mismatch: declared %d actual %d", expectCount, len(ops))
	}
	return ops, nil
}

func simulate(n int, ops []operation) ([]complex128, error) {
	if n <= 0 || n > 8 {
		return nil, fmt.Errorf("invalid number of qubits: %d", n)
	}
	size := 1 << n
	state := make([]complex128, size)
	state[0] = 1 + 0i
	for _, op := range ops {
		switch op.kind {
		case "H":
			if op.a < 1 || op.a > n {
				return nil, fmt.Errorf("H qubit %d out of range", op.a)
			}
			applyH(state, n-op.a)
		case "CNOT":
			if op.a < 1 || op.a > n || op.b < 1 || op.b > n {
				return nil, fmt.Errorf("CNOT qubit index out of range")
			}
			if op.a == op.b {
				return nil, fmt.Errorf("CNOT control and target cannot be the same")
			}
			applyCNOT(state, n-op.a, n-op.b)
		default:
			return nil, fmt.Errorf("unknown operation %q", op.kind)
		}
	}
	return state, nil
}

func applyH(state []complex128, pos int) {
	mask := 1 << pos
	inv := 1 / math.Sqrt2
	for i := 0; i < len(state); i++ {
		if i&mask == 0 {
			j := i | mask
			a := state[i]
			b := state[j]
			state[i] = (a + b) * complex(inv, 0)
			state[j] = (a - b) * complex(inv, 0)
		}
	}
}

func applyCNOT(state []complex128, controlPos, targetPos int) {
	maskC := 1 << controlPos
	maskT := 1 << targetPos
	for i := 0; i < len(state); i++ {
		if i&maskC != 0 && i&maskT == 0 {
			j := i | maskT
			state[i], state[j] = state[j], state[i]
		}
	}
}

func expectedState(tc testCase) []complex128 {
	size := 1 << tc.n
	state := make([]complex128, size)
	alpha := complex(1/math.Sqrt2, 0)
	state[0] = alpha
	state[psiIndex(tc.bits)] = alpha
	return state
}

func psiIndex(bits string) int {
	val := 0
	for i := 0; i < len(bits); i++ {
		val <<= 1
		if bits[i] == '1' {
			val |= 1
		}
	}
	return val
}

func compareStates(actual, expected []complex128) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("state size mismatch")
	}
	eps := 1e-8
	var phase complex128 = 1
	phaseSet := false
	for i := range actual {
		exp := expected[i]
		act := actual[i]
		if cmplx.Abs(exp) > eps {
			if !phaseSet {
				phase = act / exp
				if math.Abs(cmplx.Abs(phase)-1) > 1e-6 {
					return fmt.Errorf("state not normalized correctly")
				}
				phaseSet = true
			}
			if cmplx.Abs(act-phase*exp) > 1e-6 {
				return fmt.Errorf("state mismatch at basis index %d", i)
			}
		} else if cmplx.Abs(act) > 1e-6 {
			return fmt.Errorf("unexpected amplitude at basis index %d", i)
		}
	}
	return nil
}

func formatInput(tc testCase) []byte {
	return []byte(fmt.Sprintf("%d\n%s\n", tc.n, tc.bits))
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 1, bits: "1"},
		{n: 2, bits: "10"},
		{n: 2, bits: "11"},
		{n: 3, bits: "101"},
		{n: 3, bits: "111"},
		{n: 4, bits: "1001"},
		{n: 4, bits: "1111"},
		{n: 5, bits: "10000"},
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests {
		n := rnd.Intn(8) + 1
		b := make([]byte, n)
		b[0] = '1'
		for i := 1; i < n; i++ {
			if rnd.Intn(2) == 1 {
				b[i] = '1'
			} else {
				b[i] = '0'
			}
		}
		tests = append(tests, testCase{n: n, bits: string(b)})
	}
	return tests
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
