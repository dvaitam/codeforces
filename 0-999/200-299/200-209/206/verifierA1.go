package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "0-999/200-299/200-209/206/206A1.go"

type genParams struct {
	k      int
	a1     int64
	x, y   int64
	modulo int64
}

type testCase struct {
	input     string
	params    []genParams
	total     int
	needOrder bool
}

type task struct {
	value int64
	id    int
}

type parsedOutput struct {
	bad   int64
	tasks []task
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA1.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for i, tc := range tests {
		refOutStr, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		refRes, err := parseOutput(refOutStr, tc.total, tc.needOrder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, refOutStr)
			os.Exit(1)
		}
		if tc.needOrder {
			if len(refRes.tasks) != tc.total {
				fmt.Fprintf(os.Stderr, "reference did not print full order on test %d\n", i+1)
				os.Exit(1)
			}
			if err := validateSchedule(refRes.tasks, tc); err != nil {
				fmt.Fprintf(os.Stderr, "reference produced invalid schedule on test %d: %v\n", i+1, err)
				os.Exit(1)
			}
			if countBadPairs(refRes.tasks) != refRes.bad {
				fmt.Fprintf(os.Stderr, "reference bad pair count mismatch on test %d\n", i+1)
				os.Exit(1)
			}
		}

		candOutStr, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, candOutStr)
			os.Exit(1)
		}
		candRes, err := parseOutput(candOutStr, tc.total, tc.needOrder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\n", i+1, err, candOutStr)
			os.Exit(1)
		}

		if candRes.bad != refRes.bad {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d bad pairs, got %d\n", i+1, refRes.bad, candRes.bad)
			fmt.Fprintf(os.Stderr, "input:\n%s", tc.input)
			os.Exit(1)
		}

		if tc.needOrder {
			if len(candRes.tasks) != tc.total {
				fmt.Fprintf(os.Stderr, "test %d requires printing full order of %d tasks, but got %d lines\n", i+1, tc.total, len(candRes.tasks))
				os.Exit(1)
			}
		} else if len(candRes.tasks) != 0 && len(candRes.tasks) != tc.total {
			fmt.Fprintf(os.Stderr, "test %d: candidate printed %d tasks, expected either 0 or %d\n", i+1, len(candRes.tasks), tc.total)
			os.Exit(1)
		}

		if len(candRes.tasks) > 0 {
			if err := validateSchedule(candRes.tasks, tc); err != nil {
				fmt.Fprintf(os.Stderr, "invalid schedule on test %d: %v\n", i+1, err)
				os.Exit(1)
			}
			if countBadPairs(candRes.tasks) != candRes.bad {
				fmt.Fprintf(os.Stderr, "bad pair count mismatch on test %d: claimed %d but actual %d\n", i+1, candRes.bad, countBadPairs(candRes.tasks))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "206A1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	return out.String(), nil
}

func parseOutput(out string, total int, mustPrint bool) (parsedOutput, error) {
	var res parsedOutput
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return res, err
	}
	i := 0
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}
	if i == len(lines) {
		return res, errors.New("output is empty")
	}
	first := strings.TrimSpace(lines[i])
	i++
	val, err := strconv.ParseInt(first, 10, 64)
	if err != nil {
		return res, fmt.Errorf("invalid first line %q: %v", first, err)
	}
	res.bad = val
	for ; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return res, fmt.Errorf("expected two integers per line, got %q", line)
		}
		v, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return res, fmt.Errorf("invalid value %q: %v", fields[0], err)
		}
		id, err := strconv.Atoi(fields[1])
		if err != nil {
			return res, fmt.Errorf("invalid scientist id %q: %v", fields[1], err)
		}
		res.tasks = append(res.tasks, task{value: v, id: id})
	}
	if mustPrint {
		if len(res.tasks) != total {
			return res, fmt.Errorf("expected %d task lines, got %d", total, len(res.tasks))
		}
	} else if len(res.tasks) != 0 && len(res.tasks) != total {
		return res, fmt.Errorf("unexpected number of task lines: got %d, expected 0 or %d", len(res.tasks), total)
	}
	return res, nil
}

type generator struct {
	params   genParams
	produced int
	curr     int64
	initDone bool
}

func (g *generator) next() (int64, bool) {
	if g.produced >= g.params.k {
		return 0, false
	}
	var val int64
	if !g.initDone {
		val = g.params.a1
		g.curr = val
		g.initDone = true
	} else {
		g.curr = (g.curr*g.params.x + g.params.y) % g.params.modulo
		val = g.curr
	}
	g.produced++
	return val, true
}

func validateSchedule(tasks []task, tc testCase) error {
	if len(tasks) != tc.total {
		return fmt.Errorf("expected %d tasks, got %d", tc.total, len(tasks))
	}
	gens := make([]generator, len(tc.params))
	for i, p := range tc.params {
		gens[i] = generator{params: p}
	}
	for idx, t := range tasks {
		id := t.id - 1
		if id < 0 || id >= len(gens) {
			return fmt.Errorf("task %d has invalid scientist id %d", idx+1, t.id)
		}
		val, ok := gens[id].next()
		if !ok {
			return fmt.Errorf("task %d repeats scientist %d more than %d times", idx+1, t.id, tc.params[id].k)
		}
		if val != t.value {
			return fmt.Errorf("task %d mismatch: expected %d from scientist %d, got %d", idx+1, val, t.id, t.value)
		}
	}
	for i, g := range gens {
		if g.produced != tc.params[i].k {
			return fmt.Errorf("scientist %d used %d tasks instead of %d", i+1, g.produced, tc.params[i].k)
		}
	}
	return nil
}

func countBadPairs(tasks []task) int64 {
	var bad int64
	for i := 1; i < len(tasks); i++ {
		if tasks[i-1].value > tasks[i].value {
			bad++
		}
	}
	return bad
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20602060))
	var tests []testCase

	tests = append(tests, buildCase([]genParams{
		{k: 2, a1: 1, x: 1, y: 0, modulo: 100},
		{k: 2, a1: 3, x: 1, y: 0, modulo: 100},
	}))

	tests = append(tests, buildCase([]genParams{
		{k: 3, a1: 5, x: 2, y: 3, modulo: 11},
		{k: 1, a1: 7, x: 4, y: 5, modulo: 13},
		{k: 0, a1: 0, x: 1, y: 1, modulo: 2},
	}))

	tests = append(tests, buildCase([]genParams{
		{k: 100000, a1: 5, x: 17, y: 23, modulo: 1_000_000_000},
		{k: 100000, a1: 7, x: 31, y: 11, modulo: 1_000_000_000},
	}))

	tests = append(tests, buildCase([]genParams{
		{k: 70000, a1: 2, x: 37, y: 19, modulo: 999_999_937},
		{k: 70000, a1: 3, x: 73, y: 97, modulo: 1_000_000_000},
		{k: 70000, a1: 5, x: 53, y: 59, modulo: 999_999_893},
	}))

	for i := 0; i < 30; i++ {
		tests = append(tests, randomCase(rng, 1+rng.Intn(6), 5000))
	}

	return tests
}

func buildCase(params []genParams) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(params))
	total := 0
	for _, p := range params {
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", p.k, p.a1, p.x, p.y, p.modulo)
		total += p.k
	}
	return testCase{
		input:     sb.String(),
		params:    append([]genParams(nil), params...),
		total:     total,
		needOrder: total <= 200000,
	}
}

func randomCase(rng *rand.Rand, n int, maxK int) testCase {
	params := make([]genParams, n)
	total := 0
	for i := 0; i < n; i++ {
		k := rng.Intn(maxK) + 1
		a1 := rng.Int63n(1_000_000_000)
		mod := rng.Int63n(1_000_000_000-1) + 2
		if a1 >= mod {
			a1 = mod - 1
		}
		params[i] = genParams{
			k:      k,
			a1:     a1,
			x:      rng.Int63n(1_000_000_000) + 1,
			y:      rng.Int63n(1_000_000_000) + 1,
			modulo: mod,
		}
		total += k
	}
	return buildCase(params)
}
