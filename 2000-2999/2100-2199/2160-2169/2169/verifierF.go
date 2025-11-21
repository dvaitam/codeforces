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

type testInput struct {
	text string
}

func buildReference() (string, error) {
	refDir := filepath.Join("2000-2999", "2100-2199", "2160-2169", "2169")
	tmp, err := os.CreateTemp("", "ref2169F")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "2169F.go")
	cmd.Dir = refDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return tmpPath, nil
}

func commandForPath(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := commandForPath(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func normalizeOutput(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}

func fixedTests() []testInput {
	return []testInput{
		{"1\n4 5 3\n1 1 2\n4\n1\n4 3\n"},
		{"1\n3 5 2\n1 1\n5\n2\n"},
		{"1\n2 5 2\n1 1\n1\n2\n"},
		{"1\n5 6 2\n5 5\n1 2 3 4 5\n2 3 4 5 6\n"},
	}
}

func randomInt64(rng *rand.Rand, lo, hi int64) int64 {
	return lo + rng.Int63n(hi-lo+1)
}

func generateArrays(rng *rand.Rand, k int, m int64, lens []int) [][]int64 {
	arrs := make([][]int64, k)
	for i := 0; i < k; i++ {
		cur := make([]int64, 0, lens[i])
		used := make(map[int64]struct{}, lens[i])
		for len(cur) < lens[i] {
			val := randomInt64(rng, 1, m)
			if _, ok := used[val]; ok {
				continue
			}
			used[val] = struct{}{}
			cur = append(cur, val)
		}
		arrs[i] = cur
	}
	return arrs
}

func buildTest(n int, m int64, k int, lens []int, arrays [][]int64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i, v := range lens {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i := 0; i < k; i++ {
		for j, val := range arrays[i] {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomTests() []testInput {
	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		n := int(randomInt64(rng, 2, 200000))
		m := randomInt64(rng, 5, 100000000)
		k := int(randomInt64(rng, 2, int64(n)))
		if k > 60 {
			k = 60
		}
		lens := make([]int, k)
		remaining := n
		for i := 0; i < k; i++ {
			maxAllowed := 5
			minRemain := k - i - 1
			if remaining-minRemain < maxAllowed {
				maxAllowed = remaining - minRemain
			}
			if maxAllowed < 1 {
				maxAllowed = 1
			}
			lenVal := 1
			if maxAllowed > 1 {
				lenVal += rng.Intn(maxAllowed)
			}
			lens[i] = lenVal
			remaining -= lenVal
		}
		arrays := generateArrays(rng, k, m, lens)
		tests = append(tests, testInput{text: buildTest(n, m, k, lens, arrays)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := randomTests()
	for idx, input := range tests {
		expect, err := runBinary(ref, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input.text)
			os.Exit(1)
		}
		got, err := runBinary(bin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input.text)
			os.Exit(1)
		}
		if normalizeOutput(expect) != normalizeOutput(got) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input.text, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
