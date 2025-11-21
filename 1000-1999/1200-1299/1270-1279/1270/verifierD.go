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
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	referenceSource = "1000-1999/1200-1299/1270-1279/1270/1270D.go"
	totalTests      = 50
)

type testCase struct {
	name string
	n    int
	k    int
	m    int
	arr  []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary-or-source")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()

	for i, tc := range tests {
		if err := runInteractive(refBin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\n", i+1, tc.name, err)
			os.Exit(1)
		}
	}

	for i, tc := range tests {
		if err := runInteractive(candidate, tc); err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\n", i+1, tc.name, err)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	tmp, err := os.CreateTemp("", "ref1270D-*")
	if err != nil {
		return "", nil, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmpPath, referenceSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpPath)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	cleanup := func() { os.Remove(tmpPath) }
	return tmpPath, cleanup, nil
}

func runInteractive(target string, tc testCase) error {
	cmd, err := prepareCommand(target)
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start program: %w", err)
	}
	finished := false
	defer func() {
		if !finished && cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	writer := bufio.NewWriter(stdin)
	reader := bufio.NewReader(stdout)

	if _, err := fmt.Fprintf(writer, "%d %d\n", tc.n, tc.k); err != nil {
		return fmt.Errorf("failed to write initial input: %v", err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush initial input: %v", err)
	}

	queries := 0
	answered := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if answered {
					break
				}
				return fmt.Errorf("unexpected EOF before final answer (queries=%d)\nstderr:\n%s", queries, stderr.String())
			}
			return fmt.Errorf("failed to read output: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch line[0] {
		case '?':
			args := strings.Fields(line[1:])
			if len(args) != tc.k {
				return fmt.Errorf("expected %d indices in query, got %d", tc.k, len(args))
			}
			positions := make([]int, tc.k)
			seen := make([]bool, tc.n+1)
			for i, tok := range args {
				idx, err := strconv.Atoi(tok)
				if err != nil {
					return fmt.Errorf("failed to parse index %d: %v", i+1, err)
				}
				if idx < 1 || idx > tc.n {
					return fmt.Errorf("query index %d out of range", idx)
				}
				if seen[idx] {
					return fmt.Errorf("query contains duplicate index %d", idx)
				}
				seen[idx] = true
				positions[i] = idx
			}
			queries++
			if queries > tc.n {
				return fmt.Errorf("query limit exceeded (%d > %d)", queries, tc.n)
			}
			posAns, valAns := simulateDevice(tc, positions)
			if _, err := fmt.Fprintf(writer, "%d %d\n", posAns, valAns); err != nil {
				return fmt.Errorf("failed to write answer: %v", err)
			}
			if err := writer.Flush(); err != nil {
				return fmt.Errorf("failed to flush answer: %v", err)
			}
		case '!':
			args := strings.Fields(line[1:])
			if len(args) == 0 {
				return fmt.Errorf("expected answer after '!'")
			}
			ans, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid answer token %q: %v", args[0], err)
			}
			if ans != tc.m {
				return fmt.Errorf("wrong answer: expected %d got %d", tc.m, ans)
			}
			answered = true
			writer.Flush()
			goto finish
		default:
			return fmt.Errorf("unexpected output line: %q", line)
		}
	}

finish:
	if !answered {
		return fmt.Errorf("program terminated without reporting an answer")
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	finished = true
	return nil
}

func prepareCommand(target string) (*exec.Cmd, error) {
	if strings.HasSuffix(target, ".go") {
		abs, err := filepath.Abs(target)
		if err != nil {
			return nil, fmt.Errorf("resolve candidate path: %w", err)
		}
		return exec.Command("go", "run", abs), nil
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return nil, fmt.Errorf("resolve candidate path: %w", err)
	}
	return exec.Command(abs), nil
}

type deviceItem struct {
	pos int
	val int
}

func simulateDevice(tc testCase, positions []int) (int, int) {
	items := make([]deviceItem, len(positions))
	for i, pos := range positions {
		items[i] = deviceItem{pos: pos, val: tc.arr[pos-1]}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].val == items[j].val {
			return items[i].pos < items[j].pos
		}
		return items[i].val < items[j].val
	})
	chosen := items[tc.m-1]
	return chosen.pos, chosen.val
}

func generateTests() []testCase {
	tests := []testCase{
		{
			name: "sample",
			n:    4,
			k:    3,
			m:    3,
			arr:  []int{2, 0, 1, 9},
		},
		{
			name: "minimal",
			n:    2,
			k:    1,
			m:    1,
			arr:  []int{5, 1},
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, testCase{
		name: "max_case",
		n:    500,
		k:    499,
		m:    250,
		arr:  randomArray(500, rng),
	})

	for len(tests) < totalTests {
		n := rng.Intn(499) + 2
		k := rng.Intn(n-1) + 1
		m := rng.Intn(k) + 1
		arr := randomArray(n, rng)
		tests = append(tests, testCase{
			name: fmt.Sprintf("rand_%d", len(tests)),
			n:    n,
			k:    k,
			m:    m,
			arr:  arr,
		})
	}
	return tests
}

func randomArray(n int, rng *rand.Rand) []int {
	arr := make([]int, n)
	base := rng.Intn(1_000_000_000 - n)
	for i := 0; i < n; i++ {
		arr[i] = base + i
	}
	rng.Shuffle(n, func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr
}
