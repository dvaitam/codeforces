package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func compileRef() (string, error) {
	out := filepath.Join(os.TempDir(), fmt.Sprintf("refB_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", out, "1641B.go")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func runCandidate(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(42))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		var b strings.Builder
		fmt.Fprintf(&b, "1\n%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", rng.Intn(20)+1)
		}
		b.WriteByte('\n')
		tests[i] = b.String()
	}
	return tests
}

func parseInput(input string) ([]int, error) {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return nil, fmt.Errorf("bad input")
	}
	// skip T=1
	if !sc.Scan() {
		return nil, fmt.Errorf("missing n")
	}
	n, _ := strconv.Atoi(sc.Text())
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("bad array")
		}
		v, _ := strconv.Atoi(sc.Text())
		arr[i] = v
	}
	return arr, nil
}

func applyOperations(arr []int, ops [][2]int) ([]int, error) {
	res := append([]int{}, arr...)
	for _, op := range ops {
		p, c := op[0], op[1]
		if p < 0 || p > len(res) {
			return nil, fmt.Errorf("invalid position")
		}
		tmp := append([]int{}, res[:p]...)
		tmp = append(tmp, c, c)
		tmp = append(tmp, res[p:]...)
		res = tmp
	}
	return res, nil
}

func verifySegments(arr []int, segs []int) error {
	sum := 0
	for _, l := range segs {
		if l%2 != 0 {
			return fmt.Errorf("segment length not even")
		}
		if sum+l > len(arr) {
			return fmt.Errorf("segment overflow")
		}
		for i := 0; i < l/2; i++ {
			if arr[sum+i] != arr[sum+l/2+i] {
				return fmt.Errorf("segment not tandem")
			}
		}
		sum += l
	}
	if sum != len(arr) {
		return fmt.Errorf("segments do not cover array")
	}
	return nil
}

func verifyOutput(input, output string, expectImpossible bool) error {
	sc := bufio.NewScanner(strings.NewReader(output))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return fmt.Errorf("no output")
	}
	first := sc.Text()
	if first == "-1" {
		if expectImpossible {
			if sc.Scan() {
				return fmt.Errorf("extra output")
			}
			return nil
		}
		return fmt.Errorf("expected solution but got -1")
	}
	if expectImpossible {
		return fmt.Errorf("expected -1 but got solution")
	}
	q, err := strconv.Atoi(first)
	if err != nil {
		return fmt.Errorf("bad q")
	}
	ops := make([][2]int, q)
	for i := 0; i < q; i++ {
		if !sc.Scan() {
			return fmt.Errorf("missing op p")
		}
		p, _ := strconv.Atoi(sc.Text())
		if !sc.Scan() {
			return fmt.Errorf("missing op c")
		}
		c, _ := strconv.Atoi(sc.Text())
		ops[i] = [2]int{p, c}
	}
	if !sc.Scan() {
		return fmt.Errorf("missing d")
	}
	d, _ := strconv.Atoi(sc.Text())
	segs := make([]int, d)
	for i := 0; i < d; i++ {
		if !sc.Scan() {
			return fmt.Errorf("missing segment")
		}
		segs[i], _ = strconv.Atoi(sc.Text())
	}
	if sc.Scan() {
		return fmt.Errorf("extra output")
	}

	arr, err := parseInput(input)
	if err != nil {
		return err
	}
	arr2, err := applyOperations(arr, ops)
	if err != nil {
		return err
	}
	if len(arr2) != len(arr)+2*q {
		return fmt.Errorf("final length mismatch")
	}
	if err := verifySegments(arr2, segs); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, t := range tests {
		expectOut, err := runBinary(ref, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expectImpossible := strings.HasPrefix(strings.TrimSpace(expectOut), "-1")

		got, err := runCandidate(candidate, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, t)
			os.Exit(1)
		}
		if err := verifyOutput(t, got, expectImpossible); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, t, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
