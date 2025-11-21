package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

const harnessTemplate = `package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"
)

type simulator struct {
	state    [4]complex128
	script   []int
	scriptID int
}

type simQubit struct {
	idx int
	sim *simulator
}

var (
	matX       = [4]complex128{0, 1, 1, 0}
	matY       = [4]complex128{0, complex(0, -1), complex(0, 1), 0}
	matZ       = [4]complex128{1, 0, 0, -1}
	matHFactor = 1 / math.Sqrt(2)
	matH       = [4]complex128{
		complex(matHFactor, 0), complex(matHFactor, 0),
		complex(matHFactor, 0), complex(-matHFactor, 0),
	}
)

func newSimulator(script []int) *simulator {
	s := &simulator{
		script: append([]int(nil), script...),
	}
	s.state[0] = 1
	return s
}

func maskFor(idx int) int {
	if idx == 0 {
		return 2
	}
	return 1
}

func (s *simulator) applySingle(idx int, mat [4]complex128) {
	mask := maskFor(idx)
	for base := 0; base < len(s.state); base++ {
		if base&mask != 0 {
			continue
		}
		i0 := base
		i1 := base | mask
		a0 := s.state[i0]
		a1 := s.state[i1]
		s.state[i0] = mat[0]*a0 + mat[1]*a1
		s.state[i1] = mat[2]*a0 + mat[3]*a1
	}
}

func (s *simulator) applyCNOT(control, target int) {
	if control == target {
		panic("control and target must differ")
	}
	cMask := maskFor(control)
	tMask := maskFor(target)
	for base := 0; base < len(s.state); base++ {
		if base&cMask == 0 || base&tMask != 0 {
			continue
		}
		zero := base
		one := base | tMask
		s.state[zero], s.state[one] = s.state[one], s.state[zero]
	}
}

func (s *simulator) measure(idx int) int {
	mask := maskFor(idx)
	var prob1 float64
	for i := 0; i < len(s.state); i++ {
		amp := s.state[i]
		if i&mask != 0 {
			prob1 += cmplx.Abs(amp) * cmplx.Abs(amp)
		}
	}
	const eps = 1e-12
	outcome := 0
	forced := -1
	if s.scriptID < len(s.script) {
		forced = s.script[s.scriptID]
		s.scriptID++
	}
	switch forced {
	case 1:
		if prob1 > eps {
			outcome = 1
		}
	case 0:
		if 1-prob1 <= eps && prob1 > eps {
			outcome = 1
		}
	default:
		if prob1 > 0.5 {
			outcome = 1
		}
	}
	if outcome == 1 && prob1 <= eps {
		outcome = 0
	}
	s.collapse(mask, outcome)
	return outcome
}

func (s *simulator) collapse(mask int, outcome int) {
	var norm float64
	for i := 0; i < len(s.state); i++ {
		if ((i & mask) != 0) == (outcome == 1) {
			amp := s.state[i]
			norm += cmplx.Abs(amp) * cmplx.Abs(amp)
			continue
		}
		s.state[i] = 0
	}
	if norm == 0 {
		return
	}
	scale := 1 / math.Sqrt(norm)
	for i := 0; i < len(s.state); i++ {
		s.state[i] *= complex(scale, 0)
	}
}

func (s *simulator) checkTarget() error {
	const tol = 1e-6
	amp := complex(1/math.Sqrt(3), 0)
	target := [4]complex128{0, amp, amp, amp}
	var phase complex128
	for i := 0; i < len(target); i++ {
		if cmplx.Abs(target[i]) > tol {
			if cmplx.Abs(s.state[i]) < tol {
				return fmt.Errorf("basis %d amplitude missing", i)
			}
			phase = s.state[i] / target[i]
			break
		}
	}
	if phase == 0 {
		return fmt.Errorf("state lacks overlap with target")
	}
	for i := 0; i < len(target); i++ {
		if cmplx.Abs(target[i]) < tol {
			if cmplx.Abs(s.state[i]) > tol {
				return fmt.Errorf("basis %d should be zero, |amp|=%.3g", i, cmplx.Abs(s.state[i]))
			}
			continue
		}
		expected := target[i] * phase
		if cmplx.Abs(s.state[i]-expected) > tol {
			return fmt.Errorf("basis %d mismatch: got %v want %v", i, s.state[i], expected)
		}
	}
	sum := 0.0
	for _, amp := range s.state {
		sum += cmplx.Abs(amp) * cmplx.Abs(amp)
	}
	if math.Abs(sum-1) > 1e-6 {
		return fmt.Errorf("state not normalized (norm %.6f)", sum)
	}
	return nil
}

func (q *simQubit) X() {
	q.sim.applySingle(q.idx, matX)
}

func (q *simQubit) Y() {
	q.sim.applySingle(q.idx, matY)
}

func (q *simQubit) Z() {
	q.sim.applySingle(q.idx, matZ)
}

func (q *simQubit) H() {
	q.sim.applySingle(q.idx, matH)
}

func (q *simQubit) CNOT(target Qubit) {
	t, ok := target.(*simQubit)
	if !ok {
		panic("CNOT target is not a simulator qubit")
	}
	q.sim.applyCNOT(q.idx, t.idx)
}

func (q *simQubit) Measure() int {
	return q.sim.measure(q.idx)
}

func runTrial(script []int) error {
	sim := newSimulator(script)
	qubits := []Qubit{
		&simQubit{idx: 0, sim: sim},
		&simQubit{idx: 1, sim: sim},
	}
	var execErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				execErr = fmt.Errorf("panic: %v", r)
			}
		}()
		PrepareTripleState(qubits)
	}()
	if execErr != nil {
		return execErr
	}
	return sim.checkTarget()
}

func TestPrepareTripleState(t *testing.T) {
	scripts := [][]int{
		nil,
		{0},
		{1},
		{0, 1},
		{1, 0},
		{0, 0, 1},
		{1, 1, 0},
	}
	for idx, script := range scripts {
		if err := runTrial(script); err != nil {
			t.Fatalf("trial %d with script %v failed: %v", idx+1, script, err)
		}
	}
}
`

func main() {
	args := os.Args[1:]
	if len(args) >= 1 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC.go /path/to/solution.go")
		os.Exit(1)
	}
	target := args[0]
	if err := runHarness(target); err != nil {
		fmt.Fprintf(os.Stderr, "verification failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}

func runHarness(target string) error {
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("failed to resolve target path: %v", err)
	}
	info, err := os.Stat(absTarget)
	if err != nil {
		return fmt.Errorf("cannot access target: %v", err)
	}
	if info.IsDir() {
		return fmt.Errorf("target must be a Go source file, got directory: %s", absTarget)
	}

	tmpDir, err := os.MkdirTemp("", "verify-1356C-")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	destPath := filepath.Join(tmpDir, filepath.Base(absTarget))
	if err := copyFile(absTarget, destPath); err != nil {
		return fmt.Errorf("failed to copy candidate: %v", err)
	}
	hPath := filepath.Join(tmpDir, "prepare_state_test.go")
	if err := os.WriteFile(hPath, []byte(harnessTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write harness: %v", err)
	}

	cmd := exec.Command("go", "test", "-run", "TestPrepareTripleState", "-count", "1")
	cmd.Dir = tmpDir
	cmd.Env = append(os.Environ(), "GO111MODULE=off")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go test failed: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		out.Close()
	}()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	if err := out.Sync(); err != nil {
		return err
	}
	return nil
}
