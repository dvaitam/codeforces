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

type query struct {
	n int
	a []int
}

type testCase struct {
	input   string
	queries []query
}

type answer struct {
	possible bool
	l, r     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
	for idx, tc := range tests {
		wantOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		wantAns, err := parseAnswers(wantOut, len(tc.queries))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		gotAns, err := parseAnswers(gotOut, len(tc.queries))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}
		for i, q := range tc.queries {
			ref := wantAns[i]
			cand := gotAns[i]
			if !ref.possible {
				if cand.possible {
					fmt.Fprintf(os.Stderr, "test %d query %d: candidate claims YES but reference says NO\nInput query length %d\n", idx+1, i+1, q.n)
					os.Exit(1)
				}
				continue
			}
			if !cand.possible {
				fmt.Fprintf(os.Stderr, "test %d query %d: candidate output NO but reference has solution\n", idx+1, i+1)
				os.Exit(1)
			}
			if cand.l < 1 || cand.r < cand.l || cand.r > q.n {
				fmt.Fprintf(os.Stderr, "test %d query %d: invalid indices %d %d for n=%d\n", idx+1, i+1, cand.l, cand.r, q.n)
				os.Exit(1)
			}
			if !isInteresting(q.a[cand.l-1 : cand.r]) {
				fmt.Fprintf(os.Stderr, "test %d query %d: subarray [%d,%d] is not interesting\nArray: %v\n", idx+1, i+1, cand.l, cand.r, q.a)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1270B.go",
		filepath.Join("1000-1999", "1200-1299", "1270-1279", "1270", "1270B.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not locate 1270B.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1270B_%d.bin", time.Now().UnixNano()))
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

func parseAnswers(out string, t int) ([]answer, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	res := make([]answer, t)
	for i := 0; i < t; i++ {
		var word string
		if _, err := fmt.Fscan(reader, &word); err != nil {
			return nil, fmt.Errorf("failed to read YES/NO for query %d: %v", i+1, err)
		}
		word = strings.ToUpper(word)
		if word == "NO" {
			res[i] = answer{}
			continue
		}
		if word != "YES" {
			return nil, fmt.Errorf("unexpected token %q for query %d", word, i+1)
		}
		var l, r int
		if _, err := fmt.Fscan(reader, &l, &r); err != nil {
			return nil, fmt.Errorf("failed to read indices for query %d: %v", i+1, err)
		}
		res[i] = answer{possible: true, l: l, r: r}
	}
	// ensure no extra tokens with non-whitespace characters
	rest, _ := io.ReadAll(reader)
	if strings.TrimSpace(string(rest)) != "" {
		return nil, fmt.Errorf("extra output detected: %q", strings.TrimSpace(string(rest)))
	}
	return res, nil
}

func isInteresting(arr []int) bool {
	if len(arr) == 0 {
		return false
	}
	minVal, maxVal := arr[0], arr[0]
	for _, v := range arr {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal-minVal >= len(arr)
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildCase([][]int{{1, 1}}),
		buildCase([][]int{{1, 2}}),
		buildCase([][]int{{1, 1, 3}}),
		buildCase([][]int{{5, 5, 5, 5}}),
		buildCase([][]int{{0, 3, 0, 3}}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		tests = append(tests, randomCase(rng, rng.Intn(4)+1, rng.Intn(20)+2))
	}
	tests = append(tests, randomCase(rng, 1, 500))
	return tests
}

func buildCase(arrays [][]int) testCase {
	var queries []query
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arrays)))
	for _, arr := range arrays {
		queries = append(queries, query{n: len(arr), a: append([]int(nil), arr...)})
		sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String(), queries: queries}
}

func randomCase(rng *rand.Rand, t, maxN int) testCase {
	if t < 1 {
		t = 1
	}
	var arrays [][]int
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN-1) + 2
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(20)
		}
		arrays = append(arrays, arr)
	}
	return buildCase(arrays)
}
