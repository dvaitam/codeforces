package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const maxQueries = 925

type hiddenTest struct {
	n      int
	arr    []int
	unique int
}

type testCase struct {
	name string
	data hiddenTest
}

var verifierDir string

func init() {
	if _, file, _, ok := runtime.Caller(0); ok {
		verifierDir = filepath.Dir(file)
	} else {
		verifierDir = "."
	}
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solver-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("failed to build %s: %v\n%s", path, err, string(out))
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func randomHiddenTest(rng *rand.Rand, n int) hiddenTest {
	unique := rng.Intn(n) + 1
	arr := make([]int, 0, 2*n-1)
	for v := 1; v <= n; v++ {
		count := 2
		if v == unique {
			count = 1
		}
		for i := 0; i < count; i++ {
			arr = append(arr, v)
		}
	}
	rng.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return hiddenTest{n: n, arr: arr, unique: unique}
}

func deterministicTests() []testCase {
	seeds := []int64{1, 2, 3, 4, 5}
	tests := make([]testCase, 0, len(seeds))
	for idx, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		n := 10 + idx*5
		tests = append(tests, testCase{
			name: fmt.Sprintf("deterministic_%d", idx+1),
			data: randomHiddenTest(rng, n),
		})
	}
	return tests
}

func generateTests() []testCase {
	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 30 {
		n := rng.Intn(40) + 10
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", len(tests)+1),
			data: randomHiddenTest(rng, n),
		})
	}
	return tests
}

func runInteraction(bin string, tests []testCase) error {
	cmd := exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	writer := bufio.NewWriter(stdin)
	fmt.Fprintf(writer, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(writer, "%d\n", tc.data.n)
	}
	if err := writer.Flush(); err != nil {
		return err
	}

	reader := bufio.NewScanner(stdout)
	reader.Buffer(make([]byte, 0, 1024), 1<<20)

	current := 0
	queriesUsed := 0
	for current < len(tests) {
		if !reader.Scan() {
			if err := reader.Err(); err != nil {
				return fmt.Errorf("failed to read output: %v", err)
			}
			return fmt.Errorf("unexpected EOF while waiting for test %d response", current+1)
		}
		line := strings.TrimSpace(reader.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		switch fields[0] {
		case "?":
			if len(fields) < 4 {
				return fmt.Errorf("malformed query on test %d: %s", current+1, line)
			}
			x, err := strconv.Atoi(fields[1])
			if err != nil {
				return fmt.Errorf("invalid x in query: %v", err)
			}
			k, err := strconv.Atoi(fields[2])
			if err != nil {
				return fmt.Errorf("invalid subset size: %v", err)
			}
			if k < 0 || len(fields) != 3+k {
				return fmt.Errorf("query subset size mismatch on test %d", current+1)
			}
			subset := make([]int, k)
			seen := make(map[int]bool)
			for i := 0; i < k; i++ {
				val, err := strconv.Atoi(fields[3+i])
				if err != nil {
					return fmt.Errorf("invalid index in query: %v", err)
				}
				if val < 1 || val > len(tests[current].data.arr) {
					return fmt.Errorf("query index %d out of range on test %d", val, current+1)
				}
				if seen[val] {
					return fmt.Errorf("duplicate index %d in query on test %d", val, current+1)
				}
				seen[val] = true
				subset[i] = val
			}
			if x < 1 || x > tests[current].data.n {
				return fmt.Errorf("query value %d out of range on test %d", x, current+1)
			}
			answer := 0
			for _, idx := range subset {
				if tests[current].data.arr[idx-1] == x {
					answer = 1
					break
				}
			}
			fmt.Fprintf(writer, "%d\n", answer)
			if err := writer.Flush(); err != nil {
				return err
			}
			queriesUsed++
			if queriesUsed > maxQueries {
				return fmt.Errorf("too many queries on test %d (>%d)", current+1, maxQueries)
			}
		case "!":
			if len(fields) != 2 {
				return fmt.Errorf("malformed answer on test %d: %s", current+1, line)
			}
			val, err := strconv.Atoi(fields[1])
			if err != nil {
				return fmt.Errorf("invalid answer value: %v", err)
			}
			if val != tests[current].data.unique {
				return fmt.Errorf("wrong answer on test %d: expected %d got %d", current+1, tests[current].data.unique, val)
			}
			current++
			queriesUsed = 0
		default:
			return fmt.Errorf("unexpected output: %s", line)
		}
	}

	writer.Flush()
	stdin.Close()
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("program exited with error: %v", err)
	}
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin, cleanup, err := buildBinary(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	if err := runInteraction(bin, tests); err != nil {
		fmt.Fprintf(os.Stderr, "verification failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
