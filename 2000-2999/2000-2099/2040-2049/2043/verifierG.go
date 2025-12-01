package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refSource = "./2043G.go"

type testCase struct {
	name     string
	input    string
	expected []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		expect := tc.expected
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOut, len(expect))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		if !equalSlices(refAns, expect) {
			fmt.Fprintf(os.Stderr, "reference answers mismatch on test %d (%s)\nexpected: %v\nreference: %v\ninput:\n%s", idx+1, tc.name, expect, refAns, tc.input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseAnswers(candOut, len(expect))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if !equalSlices(candAns, expect) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed\nexpected: %v\ncandidate: %v\ninput:\n%s", idx+1, tc.name, expect, candAns, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2043G-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2043G.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
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
	err := cmd.Run()
	return out.String(), err
}

func parseAnswers(output string, expected int) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, token := range fields {
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		res[i] = val
	}
	return res, nil
}

func equalSlices(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func buildTests() []testCase {
	tests := []testCase{
		newTestCase("sample1", "3\n1 2 3\n5\n2 0 2\n1 0 2\n2 0 2\n1 2 0\n2 1 0\n"),
		newTestCase("sample2", "7\n1 3 4 4 7 1 3\n3\n2 1 6\n2 1 0\n2 5 6\n"),
		newTestCase("single_element", "1\n1\n3\n2 0 0\n1 0 0\n2 0 0\n"),
		newTestCase("all_updates", "4\n1 1 1 1\n4\n1 0 1\n1 1 2\n1 2 3\n2 0 3\n"),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestCase(rng, i+1))
	}
	return tests
}

func randomTestCase(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(40) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n) + 1
	}
	q := rng.Intn(150) + 1
	type query struct {
		typ int
		a   int
		b   int
	}
	queries := make([]query, q)
	hasType2 := false
	for i := 0; i < q; i++ {
		typ := 1 + rng.Intn(2)
		if i == q-1 && !hasType2 {
			typ = 2
		}
		queries[i] = query{
			typ: typ,
			a:   rng.Intn(n),
			b:   rng.Intn(n),
		}
		if typ == 2 {
			hasType2 = true
		}
	}

	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(q))
	sb.WriteByte('\n')
	for _, qu := range queries {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", qu.typ, qu.a, qu.b))
	}
	name := fmt.Sprintf("random_%d", idx)
	return newTestCase(name, sb.String())
}

func newTestCase(name, input string) testCase {
	expect, err := solveDeterministic(input)
	if err != nil {
		panic(fmt.Sprintf("failed to solve %s: %v", name, err))
	}
	return testCase{name: name, input: input, expected: expect}
}

func solveDeterministic(input string) ([]int64, error) {
	reader := strings.NewReader(input)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return nil, fmt.Errorf("failed to read n: %v", err)
	}
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if _, err := fmt.Fscan(reader, &arr[i]); err != nil {
			return nil, fmt.Errorf("failed to read a[%d]: %v", i, err)
		}
	}
	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return nil, fmt.Errorf("failed to read q: %v", err)
	}
	if q == 0 {
		return nil, fmt.Errorf("no queries provided")
	}
	results := make([]int64, 0, q)
	last := int64(0)
	modN := int64(n)
	for i := 0; i < q; i++ {
		var typ int
		if _, err := fmt.Fscan(reader, &typ); err != nil {
			return nil, fmt.Errorf("failed to read query type %d: %v", i+1, err)
		}
		if typ == 1 {
			var pPrime, xPrime int64
			if _, err := fmt.Fscan(reader, &pPrime, &xPrime); err != nil {
				return nil, fmt.Errorf("failed to read type 1 params %d: %v", i+1, err)
			}
			p := int((pPrime + last) % modN)
			x := int((xPrime + last) % modN)
			p++
			x++
			arr[p] = x
		} else if typ == 2 {
			var lPrime, rPrime int64
			if _, err := fmt.Fscan(reader, &lPrime, &rPrime); err != nil {
				return nil, fmt.Errorf("failed to read type 2 params %d: %v", i+1, err)
			}
			l := int((lPrime + last) % modN)
			r := int((rPrime + last) % modN)
			l++
			r++
			if l > r {
				l, r = r, l
			}
			length := r - l + 1
			if length <= 1 {
				results = append(results, 0)
				last = 0
				continue
			}
			freq := make(map[int]int)
			for pos := l; pos <= r; pos++ {
				freq[arr[pos]]++
			}
			var equalPairs int64
			for _, cnt := range freq {
				equalPairs += int64(cnt*(cnt-1)) / 2
			}
			len64 := int64(length)
			totalPairs := len64 * (len64 - 1) / 2
			ans := totalPairs - equalPairs
			results = append(results, ans)
			last = ans
		} else {
			return nil, fmt.Errorf("invalid query type %d at position %d", typ, i+1)
		}
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("input must contain at least one type 2 query")
	}
	return results, nil
}
