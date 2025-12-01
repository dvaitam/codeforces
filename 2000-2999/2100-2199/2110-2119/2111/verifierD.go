package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./2111D.go"
const classesPerGroup = 6

type testCase struct {
	n     int
	m     int
	rooms []int
}

type testInput struct {
	raw string
	tc  []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]
	tests := generateTests()

	for i, ti := range tests {
		refOut, err := runProgram(refBin, ti.raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, ti.raw)
			os.Exit(1)
		}
		refOK, refTotal, refErr := verify(ti, refOut)
		if !refOK {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %s\ninput:\n%s\noutput:\n%s\n", i+1, refErr, ti.raw, refOut)
			os.Exit(1)
		}

		gotOut, err := runCandidate(candidate, ti.raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, ti.raw, gotOut)
			os.Exit(1)
		}
		ok, total, reason := verify(ti, gotOut)
		if !ok {
			fmt.Fprintf(os.Stderr, "invalid output on test %d: %s\ninput:\n%s\noutput:\n%s\n", i+1, reason, ti.raw, gotOut)
			os.Exit(1)
		}
		if total != refTotal {
			fmt.Fprintf(os.Stderr, "non-optimal output on test %d\ninput:\n%s\nexpected total moves: %d\nyour total: %d\n", i+1, ti.raw, refTotal, total)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2111D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func verify(ti testInput, output string) (bool, int64, string) {
	tokens := strings.Fields(output)
	pos := 0
	totalMoves := int64(0)

	for _, tc := range ti.tc {
		roomSet := make(map[int]struct{}, tc.m)
		for _, r := range tc.rooms {
			roomSet[r] = struct{}{}
		}

		for classIdx := 0; classIdx < classesPerGroup; classIdx++ {
			used := make(map[int]struct{}, tc.n)
			_ = used
		}
		usedPerClass := make([]map[int]struct{}, classesPerGroup)
		for k := 0; k < classesPerGroup; k++ {
			usedPerClass[k] = make(map[int]struct{}, tc.n)
		}

		for g := 0; g < tc.n; g++ {
			if pos+classesPerGroup > len(tokens) {
				return false, 0, "not enough output tokens"
			}
			groupRooms := make([]int, classesPerGroup)
			for k := 0; k < classesPerGroup; k++ {
				val, err := atoi(tokens[pos+k])
				if err != nil {
					return false, 0, "non-integer output"
				}
				if _, ok := roomSet[val]; !ok {
					return false, 0, "uses unavailable classroom"
				}
				if _, exists := usedPerClass[k][val]; exists {
					return false, 0, "classroom conflict at period"
				}
				usedPerClass[k][val] = struct{}{}
				groupRooms[k] = val
			}
			pos += classesPerGroup
			totalMoves += movement(groupRooms)
		}
	}

	if pos != len(tokens) {
		return false, 0, "extra output tokens"
	}
	return true, totalMoves, ""
}

func movement(rooms []int) int64 {
	if len(rooms) <= 1 {
		return 0
	}
	var sum int64
	for i := 0; i+1 < len(rooms); i++ {
		a := rooms[i] / 100
		b := rooms[i+1] / 100
		sum += int64(abs(a - b))
	}
	return sum
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func atoi(s string) (int, error) {
	var x int
	_, err := fmt.Sscan(s, &x)
	return x, err
}

func parseInput(raw string) (testInput, error) {
	toks := strings.Fields(raw)
	if len(toks) == 0 {
		return testInput{}, fmt.Errorf("empty input")
	}
	ptr := 0
	t, err := atoi(toks[ptr])
	if err != nil {
		return testInput{}, fmt.Errorf("invalid t")
	}
	ptr++
	var cases []testCase
	for c := 0; c < t; c++ {
		if ptr+1 >= len(toks) {
			return testInput{}, fmt.Errorf("incomplete test case header")
		}
		n, err := atoi(toks[ptr])
		if err != nil {
			return testInput{}, fmt.Errorf("invalid n")
		}
		m, err := atoi(toks[ptr+1])
		if err != nil {
			return testInput{}, fmt.Errorf("invalid m")
		}
		ptr += 2
		if ptr+m > len(toks) {
			return testInput{}, fmt.Errorf("missing rooms")
		}
		rooms := make([]int, m)
		for i := 0; i < m; i++ {
			val, err := atoi(toks[ptr+i])
			if err != nil {
				return testInput{}, fmt.Errorf("bad room")
			}
			rooms[i] = val
		}
		ptr += m
		cases = append(cases, testCase{n: n, m: m, rooms: rooms})
	}
	if ptr != len(toks) {
		return testInput{}, fmt.Errorf("extra input data")
	}
	return testInput{raw: raw, tc: cases}, nil
}

func generateTests() []testInput {
	var tests []testInput
	rng := rand.New(rand.NewSource(21112111))

	tests = append(tests, mustParse(buildInput([]testCase{
		{n: 1, m: 1, rooms: []int{100}},
		{n: 2, m: 4, rooms: []int{479, 290, 478, 293}},
		{n: 3, m: 6, rooms: []int{200, 201, 300, 301, 400, 401}},
	})))

	tests = append(tests, mustParse(buildInput([]testCase{
		{n: 2, m: 2, rooms: []int{500, 600}},
		{n: 4, m: 8, rooms: []int{110, 111, 210, 211, 310, 311, 410, 411}},
	})))

	tests = append(tests, mustParse(buildInput([]testCase{
		{n: 5, m: 9, rooms: []int{100, 101, 102, 600, 601, 602, 603, 604, 605}},
	})))

	for i := 0; i < 10; i++ {
		tests = append(tests, randomTest(rng, 5, 50))
	}
	tests = append(tests, randomTest(rng, 10, 200))

	return tests
}

func mustParse(raw string, err error) testInput {
	if err != nil {
		panic(err)
	}
	ti, e := parseInput(raw)
	if e != nil {
		panic(e)
	}
	return ti
}

func buildInput(cases []testCase) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, tc := range cases {
		if tc.n > tc.m {
			return "", fmt.Errorf("invalid case: n > m")
		}
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.m)
		for i, v := range tc.rooms {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String(), nil
}

func randomTest(rng *rand.Rand, maxCases int, maxN int) testInput {
	t := rng.Intn(maxCases) + 1
	var cases []testCase
	totalM := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN-1) + 1
		m := rng.Intn(maxN-n+1) + n
		totalM += m
		rooms := randomRooms(rng, m)
		cases = append(cases, testCase{n: n, m: m, rooms: rooms})
	}
	raw, _ := buildInput(cases)
	ti, _ := parseInput(raw)
	return ti
}

func randomRooms(rng *rand.Rand, m int) []int {
	set := make(map[int]struct{}, m)
	res := make([]int, 0, m)
	for len(res) < m {
		val := randomRoomNumber(rng)
		if _, ok := set[val]; ok {
			continue
		}
		set[val] = struct{}{}
		res = append(res, val)
	}
	return res
}

func randomRoomNumber(rng *rand.Rand) int {
	floor := rng.Intn(20_000_000) + 1 // floors from 1 to 2e7 to stay within int
	room := rng.Intn(100)
	return int(math.Max(100, float64(floor*100+room)))
}
