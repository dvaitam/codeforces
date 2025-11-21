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
)

const (
	maxTotalN = 10000
	maxTotalQ = 300000
	workDir   = "."
	refSource = "2163D2.go"
	refBinary = "refD2.bin"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(ref)

	input, cases := generateTests()

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}

	expVals, err := parseOutput(refOut, cases)
	if err != nil {
		fmt.Printf("reference output parse error: %v\n", err)
		return
	}
	gotVals, err := parseOutput(candOut, cases)
	if err != nil {
		fmt.Printf("candidate output parse error: %v\n", err)
		return
	}

	for i := 0; i < cases; i++ {
		if expVals[i] != gotVals[i] {
			fmt.Printf("Mismatch on test %d: expected %d, got %d\n", i+1, expVals[i], gotVals[i])
			fmt.Println("Input (truncated):")
			fmt.Println(previewInput(input, 1200))
			return
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary, refSource)
	cmd.Dir = workDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return filepath.Join(workDir, refBinary), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, cases int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != cases {
		return nil, fmt.Errorf("expected %d numbers, got %d", cases, len(fields))
	}
	res := make([]int, cases)
	for i, tok := range fields {
		v, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("token %d is not an integer: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func generateTests() ([]byte, int) {
	rnd := rand.New(rand.NewSource(2163))
	var cases []string
	totalN, totalQ := 0, 0

	for len(cases) < 35 && totalN+4 <= maxTotalN && totalQ < maxTotalQ {
		remainingN := maxTotalN - totalN
		if remainingN < 4 {
			break
		}
		maxNForCase := remainingN
		if maxNForCase > 800 {
			maxNForCase = 800
		}
		n := 4 + rnd.Intn(maxNForCase-3)
		if totalN+n > maxTotalN {
			n = remainingN
		}

		remainingQ := maxTotalQ - totalQ
		if remainingQ <= 0 {
			break
		}
		maxPossible := remainingQ
		limit := n * (n + 1) / 2
		if maxPossible > limit {
			maxPossible = limit
		}
		if maxPossible == 0 {
			break
		}
		q := 1 + rnd.Intn(maxPossible)
		cases = append(cases, formatCase(n, q, rnd))
		totalN += n
		totalQ += q
	}

	remainingN := maxTotalN - totalN
	remainingQ := maxTotalQ - totalQ
	if remainingN >= 4 && remainingQ > 0 {
		limit := remainingN * (remainingN + 1) / 2
		if remainingQ > limit {
			remainingQ = limit
		}
		if remainingQ == 0 {
			remainingQ = 1
		}
		cases = append(cases, formatCase(remainingN, remainingQ, rnd))
	}

	if len(cases) == 0 {
		cases = append(cases, formatCase(4, 1, rnd))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(cs)
	}
	return []byte(sb.String()), len(cases)
}

func formatCase(n, q int, rnd *rand.Rand) string {
	perm := rnd.Perm(n)
	ranges := uniqueRanges(n, q, rnd)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i, v := range perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, rg := range ranges {
		fmt.Fprintf(&sb, "%d %d\n", rg[0], rg[1])
	}
	return sb.String()
}

func uniqueRanges(n, q int, rnd *rand.Rand) [][2]int {
	limit := n * (n + 1) / 2
	if q > limit {
		q = limit
	}
	if limit <= 200000 {
		all := make([][2]int, 0, limit)
		for l := 1; l <= n; l++ {
			for r := l; r <= n; r++ {
				all = append(all, [2]int{l, r})
			}
		}
		rnd.Shuffle(len(all), func(i, j int) {
			all[i], all[j] = all[j], all[i]
		})
		return all[:q]
	}
	ranges := make([][2]int, 0, q)
	seen := make(map[int]struct{}, q*2)
	for len(ranges) < q {
		l := rnd.Intn(n) + 1
		r := l + rnd.Intn(n-l+1)
		key := l*(n+1) + r
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		ranges = append(ranges, [2]int{l, r})
	}
	return ranges
}

func previewInput(data []byte, max int) string {
	if len(data) <= max {
		return string(data)
	}
	return string(data[:max]) + "\n..."
}
