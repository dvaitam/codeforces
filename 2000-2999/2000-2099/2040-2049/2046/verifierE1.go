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

const refSource = "./2046E1.go"

type participant struct {
	a, b int64
	s    int64
	city int
}

type testCase struct {
	name         string
	m            int
	participants []participant
}

type problem struct {
	d, t int64
}

type caseResult struct {
	possible bool
	problems []problem
}

func main() {
	candidatePath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}

	tests := buildTests()
	input := buildInput(tests)

	refBin, refCleanup, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := buildBinary(candidatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	refOut, err := runBinary(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\n", err)
		os.Exit(1)
	}
	refResults, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\nreference output:\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runBinary(candBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ncandidate output:\n%s", err, candOut)
		os.Exit(1)
	}
	candResults, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		refRes := refResults[i]
		candRes := candResults[i]
		if refRes.possible {
			if err := verifyArrangement(tc, refRes.problems); err != nil {
				fmt.Fprintf(os.Stderr, "internal error: reference produced invalid solution on test %d (%s): %v\n", i+1, tc.name, err)
				os.Exit(1)
			}
			if !candRes.possible {
				fmt.Fprintf(os.Stderr, "test %d (%s): expected a valid construction, but candidate output -1\n", i+1, tc.name)
				os.Exit(1)
			}
			if err := verifyArrangement(tc, candRes.problems); err != nil {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed validation: %v\nInput:\n%sCandidate output:\n%s", i+1, tc.name, err, formatSingleInput(tc), formatSingleOutput(candRes))
				os.Exit(1)
			}
		} else {
			if candRes.possible {
				fmt.Fprintf(os.Stderr, "test %d (%s): candidate produced a solution but reference determined it is impossible\n", i+1, tc.name)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier-2046E1-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(path))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(binPath, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(output string, tests []testCase) ([]caseResult, error) {
	reader := bufio.NewReader(strings.NewReader(output))
	results := make([]caseResult, len(tests))
	for i := range tests {
		tok, err := nextToken(reader)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("test %d: unexpected end of output", i+1)
			}
			return nil, fmt.Errorf("test %d: failed to read output: %v", i+1, err)
		}
		if tok == "-1" {
			results[i] = caseResult{possible: false}
			continue
		}
		p, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("test %d: invalid integer %q: %v", i+1, tok, err)
		}
		problems := make([]problem, p)
		for j := 0; j < p; j++ {
			dTok, err := nextToken(reader)
			if err != nil {
				return nil, fmt.Errorf("test %d: unable to read difficulty for problem %d: %v", i+1, j+1, err)
			}
			tTok, err := nextToken(reader)
			if err != nil {
				return nil, fmt.Errorf("test %d: unable to read topic for problem %d: %v", i+1, j+1, err)
			}
			dVal, err := strconv.ParseInt(dTok, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid difficulty %q: %v", i+1, dTok, err)
			}
			tVal, err := strconv.ParseInt(tTok, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid topic %q: %v", i+1, tTok, err)
			}
			problems[j] = problem{d: dVal, t: tVal}
		}
		results[i] = caseResult{possible: true, problems: problems}
	}
	if tok, err := nextToken(reader); err == nil && tok != "" {
		return nil, fmt.Errorf("extra output detected: %q", tok)
	}
	return results, nil
}

func nextToken(r *bufio.Reader) (string, error) {
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if !isSpace(ch) {
			var sb strings.Builder
			sb.WriteByte(ch)
			for {
				ch, err := r.ReadByte()
				if err != nil {
					if err == io.EOF {
						return sb.String(), nil
					}
					return "", err
				}
				if isSpace(ch) {
					return sb.String(), nil
				}
				sb.WriteByte(ch)
			}
		}
	}
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\t'
}

func verifyArrangement(tc testCase, problems []problem) error {
	n := len(tc.participants)
	if len(problems) == 0 {
		return fmt.Errorf("must output at least one problem")
	}
	if len(problems) > 5*n {
		return fmt.Errorf("too many problems: %d > %d", len(problems), 5*n)
	}
	seenTopics := make(map[int64]struct{}, len(problems))
	solved := make([]int, n)
	for idx, pr := range problems {
		if pr.d < 0 || pr.t < 0 || pr.d > 1_000_000_000 || pr.t > 1_000_000_000 {
			return fmt.Errorf("problem %d has out-of-range values (%d,%d)", idx+1, pr.d, pr.t)
		}
		if _, ok := seenTopics[pr.t]; ok {
			return fmt.Errorf("topic %d appears multiple times", pr.t)
		}
		seenTopics[pr.t] = struct{}{}
		for i, part := range tc.participants {
			if part.a >= pr.d || (part.s == pr.t && part.b >= pr.d) {
				solved[i]++
			}
		}
	}
	minSolved := make([]int, tc.m+1)
	maxSolved := make([]int, tc.m+1)
	const inf = int(1 << 30)
	for city := 1; city <= tc.m; city++ {
		minSolved[city] = inf
	}
	for i, part := range tc.participants {
		city := part.city
		val := solved[i]
		if val < minSolved[city] {
			minSolved[city] = val
		}
		if val > maxSolved[city] {
			maxSolved[city] = val
		}
	}
	for city := 1; city <= tc.m; city++ {
		if minSolved[city] == inf {
			return fmt.Errorf("city %d has no participants", city)
		}
	}
	for i := 1; i <= tc.m; i++ {
		for j := i + 1; j <= tc.m; j++ {
			if minSolved[i] <= maxSolved[j] {
				return fmt.Errorf("city %d min solved %d is not greater than city %d max solved %d", i, minSolved[i], j, maxSolved[j])
			}
		}
	}
	return nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		n := len(tc.participants)
		fmt.Fprintf(&sb, "%d %d\n", n, tc.m)
		for _, part := range tc.participants {
			fmt.Fprintf(&sb, "%d %d %d\n", part.a, part.b, part.s)
		}
		cityMembers := make([][]int, tc.m+1)
		for idx, part := range tc.participants {
			if part.city < 1 || part.city > tc.m {
				panic("invalid city assignment")
			}
			cityMembers[part.city] = append(cityMembers[part.city], idx+1)
		}
		for city := 1; city <= tc.m; city++ {
			members := cityMembers[city]
			if len(members) == 0 {
				panic("city without participants in generated test")
			}
			fmt.Fprintf(&sb, "%d", len(members))
			for _, idx := range members {
				fmt.Fprintf(&sb, " %d", idx)
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func formatSingleInput(tc testCase) string {
	return buildInput([]testCase{tc})
}

func formatSingleOutput(res caseResult) string {
	if !res.possible {
		return "-1\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(res.problems))
	for _, pr := range res.problems {
		fmt.Fprintf(&sb, "%d %d\n", pr.d, pr.t)
	}
	return sb.String()
}

func buildTests() []testCase {
	var tests []testCase
	tests = append(tests, simpleFeasibleTest(), simpleImpossibleTest())
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		tests = append(tests, randomFeasibleTest(rng, i))
	}
	for i := 0; i < 3; i++ {
		base := randomFeasibleTest(rng, 100+i)
		tests = append(tests, base)
		tests = append(tests, impossibleFromFeasible(base, rng, i))
	}
	for i := 0; i < 5; i++ {
		tests = append(tests, randomGeneralTest(rng, i))
	}
	return tests
}

func simpleFeasibleTest() testCase {
	parts := []participant{
		{a: 5, b: 20, s: 1, city: 1},
		{a: 2, b: 15, s: 2, city: 1},
		{a: 4, b: 4, s: 3, city: 2},
		{a: 3, b: 6, s: 2, city: 2},
	}
	return testCase{name: "simple_feasible", m: 2, participants: parts}
}

func simpleImpossibleTest() testCase {
	parts := []participant{
		{a: 5, b: 6, s: 1, city: 1},
		{a: 1, b: 5, s: 2, city: 1},
		{a: 6, b: 7, s: 3, city: 2},
		{a: 4, b: 5, s: 2, city: 2},
	}
	return testCase{name: "simple_impossible", m: 2, participants: parts}
}

func randomFeasibleTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(40) + 2
	city1Count := rng.Intn(n-1) + 1
	city2Count := n - city1Count
	if city2Count == 0 {
		city2Count = 1
		city1Count = n - 1
	}
	specRange := 5*n + 5
	city2Parts := make([]participant, 0, city2Count)
	var maxStrength int64
	bMax := make(map[int64]int64)
	for i := 0; i < city2Count; i++ {
		a := rng.Int63n(40)
		b := a + rng.Int63n(5)
		s := int64(rng.Intn(specRange))
		if b < a {
			b = a
		}
		city2Parts = append(city2Parts, participant{a: a, b: b, s: s, city: 2})
		if a > maxStrength {
			maxStrength = a
		}
		if val, ok := bMax[s]; !ok || b > val {
			bMax[s] = b
		}
	}
	city1Parts := make([]participant, 0, city1Count)
	for i := 0; i < city1Count; i++ {
		s := int64(rng.Intn(specRange))
		limit := maxStrength
		if val, ok := bMax[s]; ok && val > limit {
			limit = val
		}
		b := limit + 1 + rng.Int63n(20)
		if b < 0 {
			b = limit + 5
		}
		a := int64(0)
		if b > 0 {
			maxA := b
			if maxA > limit+5 {
				maxA = limit + 5
			}
			if maxA < 0 {
				maxA = 0
			}
			a = rng.Int63n(maxA + 1)
		}
		if a > b {
			a = b
		}
		city1Parts = append(city1Parts, participant{a: a, b: b, s: s, city: 1})
	}
	parts := append(city1Parts, city2Parts...)
	rng.Shuffle(len(parts), func(i, j int) { parts[i], parts[j] = parts[j], parts[i] })
	return testCase{
		name:         fmt.Sprintf("feasible_%d", idx),
		m:            2,
		participants: parts,
	}
}

func randomGeneralTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(60) + 2
	city1Count := rng.Intn(n-1) + 1
	city2Count := n - city1Count
	if city2Count == 0 {
		city2Count = 1
		city1Count = n - 1
	}
	specRange := 5*n + 5
	parts := make([]participant, 0, n)
	for i := 0; i < city1Count; i++ {
		a := rng.Int63n(60)
		b := a + rng.Int63n(60)
		s := int64(rng.Intn(specRange))
		if b < a {
			b = a
		}
		parts = append(parts, participant{a: a, b: b, s: s, city: 1})
	}
	for i := 0; i < city2Count; i++ {
		a := rng.Int63n(60)
		b := a + rng.Int63n(60)
		if b < a {
			b = a
		}
		s := int64(rng.Intn(specRange))
		parts = append(parts, participant{a: a, b: b, s: s, city: 2})
	}
	rng.Shuffle(len(parts), func(i, j int) { parts[i], parts[j] = parts[j], parts[i] })
	return testCase{
		name:         fmt.Sprintf("random_%d", idx),
		m:            2,
		participants: parts,
	}
}

func impossibleFromFeasible(base testCase, rng *rand.Rand, idx int) testCase {
	parts := cloneParticipants(base.participants)
	tc := testCase{
		name:         fmt.Sprintf("forced_impossible_%d", idx),
		m:            base.m,
		participants: parts,
	}
	maxStrength, specMap := city2Stats(tc)
	var city1Indices []int
	for i, part := range tc.participants {
		if part.city == 1 {
			city1Indices = append(city1Indices, i)
		}
	}
	if len(city1Indices) == 0 {
		return tc
	}
	target := city1Indices[rng.Intn(len(city1Indices))]
	spec := tc.participants[target].s
	limit := maxStrength
	if val, ok := specMap[spec]; ok && val > limit {
		limit = val
	}
	if limit < 0 {
		limit = 0
	}
	tc.participants[target].b = limit
	if tc.participants[target].a > limit {
		tc.participants[target].a = limit
	}
	return tc
}

func city2Stats(tc testCase) (int64, map[int64]int64) {
	var maxStrength int64
	specMap := make(map[int64]int64)
	for _, part := range tc.participants {
		if part.city != 2 {
			continue
		}
		if part.a > maxStrength {
			maxStrength = part.a
		}
		if val, ok := specMap[part.s]; !ok || part.b > val {
			specMap[part.s] = part.b
		}
	}
	return maxStrength, specMap
}

func cloneParticipants(parts []participant) []participant {
	cp := make([]participant, len(parts))
	copy(cp, parts)
	return cp
}
