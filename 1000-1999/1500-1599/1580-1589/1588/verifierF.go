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

const (
	refSource  = "./1588F.go"
	totalTests = 80
)

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if strings.TrimSpace(refOut) != strings.TrimSpace(candOut) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n", idx+1, tc.name, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref1588F-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1588F.bin")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCase {
	tests := []testCase{
		deterministicSample(),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-2 {
		tests = append(tests, randomSmall(rng, len(tests)+1))
	}
	tests = append(tests,
		randomLarge(rand.New(rand.NewSource(1)), "large1", 50000, 50000),
		randomLarge(rand.New(rand.NewSource(2)), "large2", 200000, 200000),
	)
	return tests
}

func deterministicSample() testCase {
	input := strings.TrimSpace(`
5
6 9 -5 3 0
2 3 1 5 4
6
1 1 5
2 1 1
1 1 5
3 1 5
2 2 -1
1 1 5
`)
	return testCase{name: "sample_like", input: input + "\n"}
}

func randomSmall(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(6) + 1
	q := rng.Intn(10) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(21) - 10
	}
	p := rng.Perm(n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(p[i] + 1))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	ensureType1 := false
	for t := 0; t < q; t++ {
		queryType := rng.Intn(3) + 1
		if !ensureType1 && t == q-1 {
			queryType = 1
		}
		switch queryType {
		case 1:
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			sb.WriteString(fmt.Sprintf("1 %d %d\n", l, r))
			ensureType1 = true
		case 2:
			v := rng.Intn(n) + 1
			x := rng.Intn(21) - 10
			sb.WriteString(fmt.Sprintf("2 %d %d\n", v, x))
		case 3:
			i := rng.Intn(n) + 1
			j := rng.Intn(n) + 1
			sb.WriteString(fmt.Sprintf("3 %d %d\n", i, j))
		}
	}
	return testCase{name: fmt.Sprintf("rand_small_%d", idx), input: sb.String()}
}

func randomLarge(rng *rand.Rand, name string, n, q int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(200001) - 100000))
	}
	sb.WriteByte('\n')
	perm := rng.Perm(n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(perm[i] + 1))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		t := rng.Intn(3) + 1
		switch t {
		case 1:
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			sb.WriteString(fmt.Sprintf("1 %d %d\n", l, r))
		case 2:
			v := rng.Intn(n) + 1
			x := rng.Intn(200001) - 100000
			sb.WriteString(fmt.Sprintf("2 %d %d\n", v, x))
		case 3:
			i := rng.Intn(n) + 1
			j := rng.Intn(n) + 1
			sb.WriteString(fmt.Sprintf("3 %d %d\n", i, j))
		}
	}
	return testCase{name: name, input: sb.String()}
}
