package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const harnessSource = `package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
)

type simulator struct {
	num   int
	state []complex128
}

var sim = &simulator{}

type Qubit struct {
	index int
}

func ensureStateSize(num int) {
	sim.num = num
	size := 1 << num
	if sim.state == nil || len(sim.state) != size {
		sim.state = make([]complex128, size)
	} else {
		for i := range sim.state {
			sim.state[i] = 0
		}
	}
}

func prepareBasis(numQubits int, basis int) {
	if basis < 0 || basis >= (1<<numQubits) {
		panic("basis index out of range")
	}
	ensureStateSize(numQubits)
	sim.state[basis] = 1
}

func validateQubit(q Qubit) {
	if q.index < 0 || q.index >= sim.num {
		panic("invalid qubit index")
	}
}

func (q Qubit) X() {
	validateQubit(q)
	applyX(q.index)
}

func X(q Qubit) {
	q.X()
}

func CNOT(control, target Qubit) {
	validateQubit(control)
	validateQubit(target)
	applyControlledX([]int{control.index}, target.index, false)
}

func CCNOT(a, b, target Qubit) {
	validateQubit(a)
	validateQubit(b)
	validateQubit(target)
	applyControlledX([]int{a.index, b.index}, target.index, false)
}

func MultiControlledX(ctrls []Qubit, target Qubit) {
	validateQubit(target)
	indices := qubitsToIndices(ctrls)
	applyControlledX(indices, target.index, false)
}

func MultiControlledXInvert(ctrls []Qubit, target Qubit) {
	validateQubit(target)
	indices := qubitsToIndices(ctrls)
	applyControlledX(indices, target.index, true)
}

func qubitsToIndices(qs []Qubit) []int {
	if len(qs) == 0 {
		return nil
	}
	out := make([]int, len(qs))
	for i, q := range qs {
		validateQubit(q)
		out[i] = q.index
	}
	return out
}

func applyX(bit int) {
	if bit < 0 || bit >= sim.num {
		panic("invalid bit index")
	}
	mask := 1 << bit
	size := len(sim.state)
	for i := 0; i < size; i++ {
		if i&mask == 0 {
			j := i | mask
			sim.state[i], sim.state[j] = sim.state[j], sim.state[i]
		}
	}
}

func applyControlledX(ctrls []int, target int, invert bool) {
	if target < 0 || target >= sim.num {
		panic("invalid target index")
	}
	if len(ctrls) == 0 {
		applyX(target)
		return
	}
	maskTarget := 1 << target
	size := len(sim.state)
	for _, c := range ctrls {
		if c < 0 || c >= sim.num {
			panic("invalid control index")
		}
		if c == target {
			panic("control cannot equal target")
		}
	}
	for basis := 0; basis < size; basis++ {
		if basis&maskTarget != 0 {
			continue
		}
		if !controlsMatch(basis, ctrls, invert) {
			continue
		}
		other := basis | maskTarget
		sim.state[basis], sim.state[other] = sim.state[other], sim.state[basis]
	}
}

func controlsMatch(val int, ctrls []int, invert bool) bool {
	for _, c := range ctrls {
		bit := (val >> c) & 1
		if invert {
			if bit != 0 {
				return false
			}
		} else {
			if bit != 1 {
				return false
			}
		}
	}
	return true
}

func expectedIndex(base int, n int) int {
	xMask := base & ((1 << n) - 1)
	yBit := (base >> n) & 1
	if xMask == (1<<n)-1 {
		yBit ^= 1
	}
	return xMask | (yBit << n)
}

func verifyBasis(expected int) error {
	const eps = 1e-9
	for i, amp := range sim.state {
		mag := cmplx.Abs(amp)
		if i == expected {
			if math.Abs(mag-1) > eps {
				return fmt.Errorf("expected amplitude 1 at %d, got %v", i, amp)
			}
		} else {
			if mag > eps {
				return fmt.Errorf("unexpected amplitude at %d: %v", i, amp)
			}
		}
	}
	return nil
}

func checkCase(n int) error {
	totalQubits := n + 1
	xs := make([]Qubit, n)
	for i := range xs {
		xs[i] = Qubit{index: i}
	}
	y := Qubit{index: n}
	limit := 1 << totalQubits
	for basis := 0; basis < limit; basis++ {
		prepareBasis(totalQubits, basis)
		AndOracle(xs, y)
		expected := expectedIndex(basis, n)
		if err := verifyBasis(expected); err != nil {
			xMask := basis & ((1 << n) - 1)
			yBit := (basis >> n) & 1
			return fmt.Errorf("input x=%0*b y=%d: %w", n, xMask, yBit, err)
		}
	}
	return nil
}

func main() {
	for n := 1; n <= 8; n++ {
		if err := checkCase(n); err != nil {
			fmt.Println("FAIL")
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	fmt.Println("OK")
}
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG1.go /path/to/solution.go")
		os.Exit(1)
	}
	candidatePath := os.Args[1]
	if filepath.Ext(candidatePath) != ".go" {
		fmt.Fprintln(os.Stderr, "candidate path must point to a Go source file (.go)")
		os.Exit(1)
	}

	refPath, err := locateReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to locate reference solution: %v\n", err)
		os.Exit(1)
	}

	if err := runHarness(refPath); err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\n", err)
		os.Exit(1)
	}

	if err := runHarness(candidatePath); err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}

func locateReference() (string, error) {
	_, verifierFile, _, ok := runtime.Caller(0)
	if ok {
		dir := filepath.Dir(verifierFile)
		path := filepath.Join(dir, "1115G1.go")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	if _, err := os.Stat("1000-1999/1100-1199/1110-1119/1115/1115G1.go"); err == nil {
		return "1000-1999/1100-1199/1110-1119/1115/1115G1.go", nil
	}
	return "", fmt.Errorf("1115G1.go not found near verifier")
}

func runHarness(solutionPath string) error {
	content, err := os.ReadFile(solutionPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", solutionPath, err)
	}

	tmpDir, err := os.MkdirTemp("", "verifier-g1-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	harnessFile := filepath.Join(tmpDir, "harness.go")
	if err := os.WriteFile(harnessFile, []byte(harnessSource), 0644); err != nil {
		return fmt.Errorf("failed to write harness: %v", err)
	}

	solutionCopy := filepath.Join(tmpDir, "solution.go")
	if err := os.WriteFile(solutionCopy, content, 0644); err != nil {
		return fmt.Errorf("failed to copy solution: %v", err)
	}

	cmd := exec.Command("go", "run", harnessFile, solutionCopy)
	cmd.Dir = tmpDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go run failed: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	out := strings.TrimSpace(stdout.String())
	if out != "OK" {
		return fmt.Errorf("harness reported failure:\n%s", stdout.String())
	}
	return nil
}
