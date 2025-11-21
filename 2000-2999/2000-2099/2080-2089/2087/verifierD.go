package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"time"
)

const maxQueries = 10

type testCase struct {
	handles   []string
	targetIdx int
	name      string
}

func (tc testCase) target() string {
	return tc.handles[tc.targetIdx]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{
		manualCase([]string{"Alpha"}, "Alpha", true, "single uppercase"),
		manualCase([]string{"zed", "Apple"}, "zed", true, "upper-first uppercase before lowercase"),
		manualCase([]string{"apple", "Banana"}, "Banana", false, "lower-first lowercase before uppercase"),
	}
	for i := 0; i < 200; i++ {
		cases = append(cases, randomCase(rng, i))
	}

	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (%s) failed: %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runCase(bin string, tc testCase) error {
	var stderr bytes.Buffer
	cmd := exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to open stdin pipe: %v", err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to open stdout pipe: %v", err)
	}
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start candidate: %v", err)
	}

	if _, err := fmt.Fprintf(stdin, "%d %s\n", len(tc.handles), tc.target()); err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return fmt.Errorf("failed to write initial input: %v", err)
	}

	reader := bufio.NewReader(stdout)
	queryCount := 0
	for {
		token, err := readToken(reader)
		if err != nil {
			cmd.Process.Kill()
			cmd.Wait()
			return fmt.Errorf("failed to read token: %v stderr:%s", err, stderr.String())
		}
		switch token {
		case "?":
			idxToken, err := readToken(reader)
			if err != nil {
				cmd.Process.Kill()
				cmd.Wait()
				return fmt.Errorf("failed to read query index: %v stderr:%s", err, stderr.String())
			}
			pos, err := strconv.Atoi(idxToken)
			if err != nil {
				cmd.Process.Kill()
				cmd.Wait()
				return fmt.Errorf("invalid query index %q: %v", idxToken, err)
			}
			if pos < 1 || pos > len(tc.handles) {
				cmd.Process.Kill()
				cmd.Wait()
				return fmt.Errorf("query index %d out of range 1..%d", pos, len(tc.handles))
			}
			queryCount++
			if queryCount > maxQueries {
				cmd.Process.Kill()
				cmd.Wait()
				return fmt.Errorf("too many queries: %d", queryCount)
			}
			if _, err := fmt.Fprintf(stdin, "%s\n", tc.handles[pos-1]); err != nil {
				cmd.Process.Kill()
				cmd.Wait()
				return fmt.Errorf("failed to write query response: %v", err)
			}
		case "!":
			idxToken, err := readToken(reader)
			if err != nil {
				cmd.Process.Kill()
				cmd.Wait()
				return fmt.Errorf("failed to read reported index: %v stderr:%s", err, stderr.String())
			}
			answer, err := strconv.Atoi(idxToken)
			if err != nil {
				cmd.Process.Kill()
				cmd.Wait()
				return fmt.Errorf("invalid reported index %q: %v", idxToken, err)
			}
			if answer != tc.targetIdx+1 {
				cmd.Process.Kill()
				cmd.Wait()
				return fmt.Errorf("reported index %d, expected %d", answer, tc.targetIdx+1)
			}
			stdin.Close()
			if err := cmd.Wait(); err != nil {
				return fmt.Errorf("runtime error after correct answer: %v stderr:%s", err, stderr.String())
			}
			return nil
		default:
			cmd.Process.Kill()
			cmd.Wait()
			return fmt.Errorf("unexpected token %q", token)
		}
	}
}

func readToken(r *bufio.Reader) (string, error) {
	var token string
	_, err := fmt.Fscan(r, &token)
	return token, err
}

func randomCase(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(500) + 1
	upperFirst := rng.Intn(2) == 0
	set := make(map[string]struct{}, n*2)
	handles := make([]string, 0, n)
	for len(handles) < n {
		uppercase := rng.Intn(2) == 0
		h := randomHandle(rng, uppercase)
		if _, exists := set[h]; exists {
			continue
		}
		set[h] = struct{}{}
		handles = append(handles, h)
	}
	sortHandles(handles, upperFirst)
	targetIdx := rng.Intn(n)
	return testCase{
		handles:   handles,
		targetIdx: targetIdx,
		name:      fmt.Sprintf("random-%d n=%d upperFirst=%v", idx+1, n, upperFirst),
	}
}

func manualCase(raw []string, target string, upperFirst bool, name string) testCase {
	handles := append([]string(nil), raw...)
	sortHandles(handles, upperFirst)
	targetIdx := -1
	for i, h := range handles {
		if h == target {
			targetIdx = i
			break
		}
	}
	if targetIdx == -1 {
		panic(fmt.Sprintf("target %s not found in manual case %s", target, name))
	}
	return testCase{
		handles:   handles,
		targetIdx: targetIdx,
		name:      name,
	}
}

func randomHandle(rng *rand.Rand, uppercase bool) string {
	length := rng.Intn(20) + 1
	b := make([]byte, length)
	if uppercase {
		b[0] = byte('A' + rng.Intn(26))
	} else {
		b[0] = byte('a' + rng.Intn(26))
	}
	for i := 1; i < length; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func sortHandles(handles []string, upperFirst bool) {
	sort.Slice(handles, func(i, j int) bool {
		return compare(handles[i], handles[j], upperFirst) < 0
	})
}

func compare(a, b string, upperFirst bool) int {
	ua := isUpper(a[0])
	ub := isUpper(b[0])
	if ua != ub {
		if upperFirst {
			if ua {
				return -1
			}
			return 1
		}
		if ua {
			return 1
		}
		return -1
	}
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

func isUpper(ch byte) bool {
	return ch >= 'A' && ch <= 'Z'
}
