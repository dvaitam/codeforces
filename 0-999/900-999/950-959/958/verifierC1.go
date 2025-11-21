package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "0-999/900-999/950-959/958/958C1.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "958C1.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	n int
	p int
	a []int64
}

func inputString(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.p))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func solve(tc testCase) int64 {
	var total int64
	for _, v := range tc.a {
		total += v
	}
	var best int64
	var prefix int64
	mod := int64(tc.p)
	for i := 0; i < tc.n-1; i++ {
		prefix += tc.a[i]
		cur := (prefix % mod) + ((total - prefix) % mod)
		if cur > best {
			best = cur
		}
	}
	return best
}

func runProgram(path, input string) (string, error) {
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
	err := cmd.Run()
	return out.String(), err
}

func parseAnswer(out string) (int64, error) {
	reader := strings.NewReader(out)
	var val int64
	if _, err := fmt.Fscan(reader, &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer: %v\nfull output:\n%s", err, out)
	}
	return val, nil
}

func buildReferenceBinary() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "958C1-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_958C1")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func randomArray(rng *rand.Rand, n int, maxVal int64) []int64 {
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = 1 + rng.Int63n(maxVal)
	}
	return arr
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(20250309))
	tests := []testCase{
		{n: 2, p: 2, a: []int64{1, 1}},
		{n: 3, p: 10, a: []int64{3, 4, 7}},
		{n: 5, p: 12, a: []int64{16, 3, 24, 13, 9}},
	}
	for i := 0; i < 80; i++ {
		n := rng.Intn(10) + 2
		p := rng.Intn(50) + 2
		tests = append(tests, testCase{
			n: n,
			p: p,
			a: randomArray(rng, n, 100),
		})
	}
	for i := 0; i < 80; i++ {
		n := rng.Intn(200) + 2
		p := rng.Intn(1000) + 2
		tests = append(tests, testCase{
			n: n,
			p: p,
			a: randomArray(rng, n, 1000000),
		})
	}
	for i := 0; i < 30; i++ {
		n := rng.Intn(1000) + 2
		p := rng.Intn(10000) + 2
		tests = append(tests, testCase{
			n: n,
			p: p,
			a: randomArray(rng, n, 1000000),
		})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	if bin == "--" {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := genTests()
	for i, tc := range tests {
		in := inputString(tc)
		expected := solve(tc)

		refOut, err := runProgram(refBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, refOut)
			os.Exit(1)
		}
		refAns, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference parse error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, refOut)
			os.Exit(1)
		}
		if refAns != expected {
			fmt.Fprintf(os.Stderr, "reference mismatch on test %d: expected %d got %d\ninput:\n%soutput:\n%s\n", i+1, expected, refAns, in, refOut)
			os.Exit(1)
		}

		out, runErr := runProgram(bin, in)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\ninput:\n%soutput:\n%s\n", i+1, runErr, in, out)
			os.Exit(1)
		}
		ans, err := parseAnswer(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, out)
			os.Exit(1)
		}
		if ans != expected {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\ninput:\n%soutput:\n%s\n", i+1, expected, ans, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
