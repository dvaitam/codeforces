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
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierF.go /path/to/candidate")
	}
	candidate := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	numCases := 30
	for tc := 0; tc < numCases; tc++ {
		inputData, answerCount := genCase(rng)

		refOut, err := runProgram(exec.Command(refBin), inputData)
		if err != nil {
			fail("reference execution failed on case %d: %v", tc+1, err)
		}
		expected, err := parseInts(refOut, answerCount)
		if err != nil {
			fail("invalid reference output on case %d: %v", tc+1, err)
		}

		candOut, err := runProgram(commandFor(candidate), inputData)
		if err != nil {
			fail("candidate execution failed on case %d: %v", tc+1, err)
		}
		got, err := parseInts(candOut, answerCount)
		if err != nil {
			fail("invalid candidate output on case %d: %v", tc+1, err)
		}

		for i := 0; i < answerCount; i++ {
			if got[i] != expected[i] {
				fail("wrong answer on case %d query %d: expected %d, got %d", tc+1, i+1, expected[i], got[i])
			}
		}
	}

	fmt.Println("OK")
}

func genCase(rng *rand.Rand) ([]byte, int) {
	n := rng.Intn(8) + 2 // 2..9
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(201)) - 100 // -100..100
		b[i] = int64(rng.Intn(201)) - 100
	}
	q := rng.Intn(10) + 1
	type query struct {
		t    int
		p    int
		x    int64
		l, r int
	}
	queries := make([]query, q)
	answerCount := 0
	for i := 0; i < q; i++ {
		t := rng.Intn(3) + 1
		if t == 1 || t == 2 {
			p := rng.Intn(n) + 1
			x := int64(rng.Intn(201)) - 100
			queries[i] = query{t: t, p: p, x: x}
		} else {
			l := rng.Intn(n-1) + 1
			r := l + rng.Intn(n-l) + 1
			queries[i] = query{t: 3, l: l, r: r}
			answerCount++
		}
	}
	// Ensure at least one type-3 query
	if answerCount == 0 {
		l := rng.Intn(n-1) + 1
		r := l + rng.Intn(n-l) + 1
		queries = append(queries, query{t: 3, l: l, r: r})
		q++
		answerCount++
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", b[i])
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", q)
	for _, qr := range queries {
		if qr.t == 1 || qr.t == 2 {
			fmt.Fprintf(&sb, "%d %d %d\n", qr.t, qr.p, qr.x)
		} else {
			fmt.Fprintf(&sb, "%d %d %d\n", qr.t, qr.l, qr.r)
		}
	}
	return []byte(sb.String()), answerCount
}

func buildReference() (string, error) {
	refPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if refPath == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}

	// Read reference source to detect language
	srcBytes, err := os.ReadFile(refPath)
	if err != nil {
		return "", fmt.Errorf("cannot read reference source: %v", err)
	}
	srcContent := string(srcBytes)

	tmp, err := os.CreateTemp("", "2042F-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	if strings.Contains(srcContent, "#include") {
		// C++ source: copy to .cpp file, compile with g++
		cppPath := tmp.Name() + ".cpp"
		if err := os.WriteFile(cppPath, srcBytes, 0644); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("failed to write cpp source: %v", err)
		}
		defer os.Remove(cppPath)
		cmd := exec.Command("g++", "-O2", "-o", tmp.Name(), cppPath)
		var combined bytes.Buffer
		cmd.Stdout = &combined
		cmd.Stderr = &combined
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("%v\n%s", err, combined.String())
		}
	} else {
		// Go source
		cmd := exec.Command("go", "build", "-o", tmp.Name(), refPath)
		var combined bytes.Buffer
		cmd.Stdout = &combined
		cmd.Stderr = &combined
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("%v\n%s", err, combined.String())
		}
	}
	return tmp.Name(), nil
}

func runProgram(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func parseInts(out string, expectedCount int) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([]int64, 0, expectedCount)
	for len(ans) < expectedCount {
		var v int64
		if _, err := fmt.Fscan(reader, &v); err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("expected %d answers, got %d", expectedCount, len(ans))
			}
			return nil, err
		}
		ans = append(ans, v)
	}
	return ans, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
