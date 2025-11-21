package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type test struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, t := range tests {
		want, err := runProgram(refBin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, t.input, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, t.input, err)
			os.Exit(1)
		}
		if normalize(got) != normalize(want) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nInput:\n%sExpected:\n%sGot:\n%s\n", idx+1, t.input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"720F.go",
		filepath.Join("0-999", "700-799", "720-729", "720", "720F.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 720F.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref720F_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func generateTests() []test {
	var tests []test
	tests = append(tests,
		newTest([]int{5}, 1),
		newTest([]int{1, 2}, 3),
		newTest([]int{-1, -2, -3}, 1),
		newTest([]int{100000, -100000, 100000}, 6),
	)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 150 {
		var n int
		switch {
		case len(tests) < 50:
			n = rng.Intn(8) + 1
		case len(tests) < 90:
			n = rng.Intn(50) + 10
		case len(tests) < 120:
			n = rng.Intn(500) + 100
		default:
			n = rng.Intn(5000) + 5000
		}
		if n > 100000 {
			n = 100000
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(100001) - 50000
		}
		maxK := int64(n) * int64(n+1) / 2
		var k int64
		if len(tests)%10 == 0 {
			k = maxK
		} else {
			k = rng.Int63n(maxK) + 1
		}
		if k > maxK {
			k = maxK
		}
		tests = append(tests, newTest(arr, k))
	}
	return tests
}

func newTest(arr []int, k int64) test {
	n := len(arr)
	if k < 1 {
		k = 1
	}
	maxK := int64(n) * int64(n+1) / 2
	if k > maxK {
		k = maxK
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmt.Sprintf("%d", v))
	}
	b.WriteByte('\n')
	return test{input: b.String()}
}
