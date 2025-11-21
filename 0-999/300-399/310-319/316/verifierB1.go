package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const testCount = 200

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "316B1.go")
	tmp, err := os.CreateTemp("", "oracle316B1")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	return path, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(100) + 1
	x := r.Intn(n) + 1
	perm := r.Perm(n)
	a := make([]int, n+1)
	zeros := 0
	for idx, p := range perm {
		id := p + 1
		if idx == 0 || r.Intn(100) < 35 || (zeros < 1 && idx == len(perm)-1) {
			a[id] = 0
			zeros++
			continue
		}
		prevID := perm[idx-1] + 1
		a[id] = prevID
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, x)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string) ([]int, error) {
	out = strings.TrimSpace(out)
	if len(out) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	reader := strings.NewReader(out)
	var res []int
	for {
		var v int
		_, err := fmt.Fscan(reader, &v)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no numbers in output")
	}
	for i := 1; i < len(res); i++ {
		if res[i] <= res[i-1] {
			return nil, fmt.Errorf("output not strictly increasing")
		}
	}
	return res, nil
}

func compareOutputs(expectStr, gotStr string) error {
	expect, err := parseOutput(expectStr)
	if err != nil {
		return fmt.Errorf("oracle output invalid: %v", err)
	}
	got, err := parseOutput(gotStr)
	if err != nil {
		return fmt.Errorf("contestant output invalid: %v", err)
	}
	if len(expect) != len(got) {
		return fmt.Errorf("different number of answers: expect %d got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			return fmt.Errorf("answer mismatch at index %d: expect %d got %d", i, expect[i], got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(1))
	for i := 0; i < testCount; i++ {
		input := genCase(r)
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := compareOutputs(expect, got); err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", i+1, input, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
