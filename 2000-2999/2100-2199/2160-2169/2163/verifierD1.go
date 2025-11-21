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
	"time"
)

const (
	randomTests = 200
	stressN     = 10000
	stressQ     = 300000
)

type caseData struct {
	n      int
	perm   []int
	ranges [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		return
	}

	candidatePath, candidateCleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("failed to prepare contestant binary:", err)
		return
	}
	if candidateCleanup != nil {
		defer candidateCleanup()
	}

	dir := sourceDir()
	oracleSrc := filepath.Join(dir, "2163D1.go")
	oraclePath, oracleCleanup, err := prepareOracle(oracleSrc)
	if err != nil {
		fmt.Println("failed to prepare reference solution:", err)
		return
	}
	defer oracleCleanup()

	deterministic := deterministicInputs()
	total := 0
	for idx, input := range deterministic {
		if err := runCase(candidatePath, oraclePath, input); err != nil {
			fmt.Printf("deterministic case %d failed: %v\ninput:\n%s", idx+1, err, input)
			return
		}
		total++
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTests; i++ {
		input := generateRandomInput(rng)
		if err := runCase(candidatePath, oraclePath, input); err != nil {
			fmt.Printf("random case %d failed: %v\ninput:\n%s", i+1, err, input)
			return
		}
		total++
	}

	fmt.Printf("All %d tests passed.\n", total)
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("contestant_%d", time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", tmp, abs)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, nil, nil
}

func prepareOracle(src string) (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracleD1_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func runCase(bin, oracle, input string) error {
	expect, err := runBinary(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	got, err := runBinary(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func deterministicInputs() []string {
	smallA := caseData{
		n:    4,
		perm: []int{0, 3, 1, 2},
		ranges: [][2]int{
			{1, 2},
			{2, 4},
			{1, 3},
		},
	}
	smallB := caseData{
		n:    6,
		perm: []int{2, 5, 1, 3, 0, 4},
		ranges: [][2]int{
			{1, 6},
			{2, 5},
			{1, 3},
			{4, 6},
			{3, 4},
		},
	}
	smallC := caseData{
		n:    5,
		perm: []int{4, 0, 2, 1, 3},
		ranges: [][2]int{
			{1, 1},
			{2, 2},
			{3, 3},
			{4, 4},
			{5, 5},
		},
	}
	fullRange := allRangesCase(8)
	stress := sequentialCase(stressN, stressQ)

	return []string{
		buildInput([]caseData{smallA}),
		buildInput([]caseData{smallA, smallB, smallC}),
		buildInput([]caseData{fullRange}),
		buildInput([]caseData{stress}),
	}
}

func allRangesCase(n int) caseData {
	perm := rotatedPermutation(n, 3)
	total := n * (n + 1) / 2
	ranges := make([][2]int, 0, total)
	for l := 1; l <= n; l++ {
		for r := l; r <= n; r++ {
			ranges = append(ranges, [2]int{l, r})
		}
	}
	return caseData{n: n, perm: perm, ranges: ranges}
}

func sequentialCase(n, q int) caseData {
	if q > n*(n+1)/2 {
		q = n * (n + 1) / 2
	}
	perm := rotatedPermutation(n, n/2+3)
	ranges := make([][2]int, 0, q)
	for l := 1; l <= n && len(ranges) < q; l++ {
		for r := l; r <= n && len(ranges) < q; r++ {
			ranges = append(ranges, [2]int{l, r})
		}
	}
	return caseData{n: n, perm: perm, ranges: ranges}
}

func rotatedPermutation(n, shift int) []int {
	if n == 0 {
		return nil
	}
	shift %= n
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = (i + shift) % n
	}
	return perm
}

func generateRandomInput(rng *rand.Rand) string {
	t := rng.Intn(4) + 1
	cases := make([]caseData, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(17) + 4 // 4..20
		maxRanges := n * (n + 1) / 2
		limit := 60
		if limit > maxRanges {
			limit = maxRanges
		}
		q := rng.Intn(limit) + 1
		cases[i] = randomCase(rng, n, q)
	}
	return buildInput(cases)
}

func randomCase(rng *rand.Rand, n, q int) caseData {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i
	}
	rng.Shuffle(n, func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	ranges := randomRanges(rng, n, q)
	return caseData{n: n, perm: perm, ranges: ranges}
}

func randomRanges(rng *rand.Rand, n, q int) [][2]int {
	used := make(map[int]struct{}, q*2)
	res := make([][2]int, 0, q)
	for len(res) < q {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		key := l*(n+1) + r
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		res = append(res, [2]int{l, r})
	}
	return res
}

func buildInput(cases []caseData) string {
	var sb strings.Builder
	sb.Grow(64)
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, c := range cases {
		if len(c.perm) != c.n {
			panic("permutation length mismatch")
		}
		fmt.Fprintf(&sb, "%d %d\n", c.n, len(c.ranges))
		for i, val := range c.perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", val)
		}
		sb.WriteByte('\n')
		for _, rg := range c.ranges {
			fmt.Fprintf(&sb, "%d %d\n", rg[0], rg[1])
		}
	}
	return sb.String()
}

func sourceDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}
