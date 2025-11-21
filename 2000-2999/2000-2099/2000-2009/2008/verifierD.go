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

const refSource = "2000-2999/2000-2099/2000-2009/2008/2008D.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for i, input := range tests {
		refOut, err := runExecutable(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		if normalize(refOut) != normalize(candOut) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2008D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, combined.String())
	}
	return tmp.Name(), nil
}

func runExecutable(path, input string) (string, error) {
	cmd := exec.Command(path)
	return execute(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return execute(cmd, input)
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

func execute(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func normalize(out string) string {
	return strings.Join(strings.Fields(out), " ")
}

func buildTests() []string {
	tests := []string{
		"5\n1\n1\n0\n5\n1 2 4 5 3\n10101\n5\n5 4 1 3 2\n10011\n6\n3 5 6 1 2 4\n010000\n6\n1 2 3 4 5 6\n100110\n",
	}

	randomConfigs := []struct {
		t      int
		maxSum int
		seed   int64
	}{
		{5, 50, 1},
		{10, 200, 2},
		{20, 1000, 3},
		{40, 5000, 4},
		{80, 200000, 5},
		{60, 200000, time.Now().UnixNano()},
	}
	for _, cfg := range randomConfigs {
		tests = append(tests, randomTest(cfg.t, cfg.maxSum, cfg.seed))
	}
	return tests
}

func randomTest(t, maxSum int, seed int64) string {
	if t < 1 {
		t = 1
	}
	r := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	remaining := maxSum
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		needCases := t - caseIdx - 1
		minRemain := needCases
		if minRemain < 0 {
			minRemain = 0
		}
		maxN := remaining - minRemain
		if maxN < 1 {
			maxN = 1
		}
		if maxN > 200000 {
			maxN = 200000
		}
		var n int
		if remaining-minRemain <= 1 {
			n = 1
		} else {
			n = r.Intn(maxN) + 1
		}
		remaining -= n
		sb.WriteString(fmt.Sprintf("%d\n", n))
		perm := randPerm(r, n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", perm[i]))
		}
		sb.WriteByte('\n')
		var str strings.Builder
		for i := 0; i < n; i++ {
			str.WriteByte(randomBit(r, i))
		}
		sb.WriteString(str.String())
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randPerm(r *rand.Rand, n int) []int {
	arr := r.Perm(n)
	res := make([]int, n)
	for i, v := range arr {
		res[i] = v + 1
	}
	return res
}

func randomBit(r *rand.Rand, pos int) byte {
	mode := r.Intn(4)
	switch mode {
	case 0:
		return '0'
	case 1:
		return '1'
	case 2:
		if pos%2 == 0 {
			return '0'
		}
		return '1'
	default:
		if r.Intn(2) == 0 {
			return '0'
		}
		return '1'
	}
}
