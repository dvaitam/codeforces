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

const (
	refSourceE       = "./1142E.go"
	randomCaseCount  = 40
	queryLimitFactor = 2
)

type testCase struct {
	name      string
	n         int
	pinkEdges [][2]int
	color     [][]uint8
	dir       [][]uint8
	valid     map[int]bool
}

type solverProcess struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr *bytes.Buffer
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceE)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomCases(rng, randomCaseCount)...)

	for idx, tc := range tests {
		if err := runSolver(refBin, tc); err != nil {
			fail("reference failed on case %d (%s): %v", idx+1, tc.name, err)
		}
		if err := runSolver(candidate, tc); err != nil {
			fail("candidate failed on case %d (%s): %v", idx+1, tc.name, err)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "1142E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(src))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func deterministicCases() []testCase {
	var cases []testCase
	cases = append(cases, singleNodeCase())
	cases = append(cases, manualCase(2, "manual-n2"))
	cases = append(cases, manualCaseWithRoot(4, 1, "manual-root1"))
	cases = append(cases, manualCaseWithRoot(5, 3, "manual-root3"))
	seedRng := rand.New(rand.NewSource(2024))
	for i := 0; i < 3; i++ {
		cases = append(cases, randomCase(seedRng, fmt.Sprintf("seeded-%d", i)))
	}
	return cases
}

func randomCases(rng *rand.Rand, count int) []testCase {
	cases := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		cases = append(cases, randomCase(rng, fmt.Sprintf("random-%d", i)))
	}
	return cases
}

func singleNodeCase() testCase {
	n := 1
	color := newMatrix(n)
	dir := newMatrix(n)
	return buildTestCase(n, color, dir, "single-node")
}

func manualCase(n int, name string) testCase {
	color := newMatrix(n)
	dir := newMatrix(n)
	root := 1
	for u := 1; u <= n; u++ {
		for v := u + 1; v <= n; v++ {
			col := uint8((u + v) % 2)
			color[u][v] = col
			color[v][u] = col
			orientation := uint8(0)
			if u == root {
				orientation = 1
			} else if v == root {
				orientation = 0
			} else if (u+v)%3 == 0 {
				orientation = 1
			}
			dir[u][v] = orientation
			dir[v][u] = 1 - orientation
		}
	}
	for v := 1; v <= n; v++ {
		if v == root {
			continue
		}
		dir[root][v] = 1
		dir[v][root] = 0
	}
	return buildTestCase(n, color, dir, name)
}

func manualCaseWithRoot(n, root int, name string) testCase {
	color := newMatrix(n)
	dir := newMatrix(n)
	for u := 1; u <= n; u++ {
		for v := u + 1; v <= n; v++ {
			col := uint8((u*v + n) % 2)
			color[u][v] = col
			color[v][u] = col
			dir[u][v] = uint8((u + v) % 2)
			dir[v][u] = 1 - dir[u][v]
		}
	}
	for v := 1; v <= n; v++ {
		if v == root {
			continue
		}
		dir[root][v] = 1
		dir[v][root] = 0
		if v%2 == 0 {
			color[root][v] = 1
			color[v][root] = 1
		}
	}
	return buildTestCase(n, color, dir, name)
}

func randomCase(rng *rand.Rand, name string) testCase {
	n := rng.Intn(39) + 2
	color := newMatrix(n)
	dir := newMatrix(n)
	for u := 1; u <= n; u++ {
		for v := u + 1; v <= n; v++ {
			col := uint8(rng.Intn(2))
			color[u][v] = col
			color[v][u] = col
			orientation := uint8(rng.Intn(2))
			dir[u][v] = orientation
			dir[v][u] = 1 - orientation
		}
	}
	root := rng.Intn(n) + 1
	for v := 1; v <= n; v++ {
		if v == root {
			continue
		}
		dir[root][v] = 1
		dir[v][root] = 0
		if rng.Intn(3) == 0 {
			color[root][v] = 1
			color[v][root] = 1
		}
	}
	return buildTestCase(n, color, dir, name)
}

func newMatrix(n int) [][]uint8 {
	mat := make([][]uint8, n+1)
	for i := range mat {
		mat[i] = make([]uint8, n+1)
	}
	return mat
}

func buildTestCase(n int, color, dir [][]uint8, name string) testCase {
	pink := make([][2]int, 0)
	for u := 1; u <= n; u++ {
		for v := 1; v <= n; v++ {
			if u == v {
				continue
			}
			if color[u][v] == 1 && dir[u][v] == 1 {
				pink = append(pink, [2]int{u, v})
			}
		}
	}
	valid := computeValid(n, color, dir)
	if len(valid) == 0 {
		panic("generated case without valid nodes")
	}
	return testCase{
		name:      name,
		n:         n,
		pinkEdges: pink,
		color:     color,
		dir:       dir,
		valid:     valid,
	}
}

func computeValid(n int, color, dir [][]uint8) map[int]bool {
	valid := make(map[int]bool)
	for s := 1; s <= n; s++ {
		pinkReach := reach(n, color, dir, s, 1)
		greenReach := reach(n, color, dir, s, 0)
		ok := true
		for v := 1; v <= n; v++ {
			if v == s {
				continue
			}
			if !pinkReach[v] && !greenReach[v] {
				ok = false
				break
			}
		}
		if ok {
			valid[s] = true
		}
	}
	return valid
}

func reach(n int, color, dir [][]uint8, start int, targetColor uint8) []bool {
	visited := make([]bool, n+1)
	stack := []int{start}
	visited[start] = true
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for v := 1; v <= n; v++ {
			if dir[u][v] == 1 && color[u][v] == targetColor && !visited[v] {
				visited[v] = true
				stack = append(stack, v)
			}
		}
	}
	return visited
}

func runSolver(target string, tc testCase) error {
	proc, err := startProcess(target)
	if err != nil {
		return err
	}
	defer proc.cmd.Process.Kill()

	writer := bufio.NewWriter(proc.stdin)
	fmt.Fprintf(writer, "%d %d\n", tc.n, len(tc.pinkEdges))
	for _, e := range tc.pinkEdges {
		fmt.Fprintf(writer, "%d %d\n", e[0], e[1])
	}
	writer.Flush()

	reader := bufio.NewReader(proc.stdout)
	queryLimit := queryLimitFactor * tc.n
	queries := 0

	for {
		token, err := readToken(reader)
		if err != nil {
			proc.cmd.Process.Kill()
			proc.cmd.Wait()
			return fmt.Errorf("unexpected EOF: %v (stderr: %s)", err, proc.stderr.String())
		}
		switch token {
		case "?":
			u, err := readIntToken(reader)
			if err != nil {
				proc.cmd.Process.Kill()
				proc.cmd.Wait()
				return fmt.Errorf("failed to read first node in query: %v", err)
			}
			v, err := readIntToken(reader)
			if err != nil {
				proc.cmd.Process.Kill()
				proc.cmd.Wait()
				return fmt.Errorf("failed to read second node in query: %v", err)
			}
			if u < 1 || u > tc.n || v < 1 || v > tc.n || u == v {
				proc.cmd.Process.Kill()
				proc.cmd.Wait()
				return fmt.Errorf("invalid query nodes (%d,%d)", u, v)
			}
			queries++
			if queries > queryLimit {
				proc.cmd.Process.Kill()
				proc.cmd.Wait()
				return fmt.Errorf("query limit exceeded (%d > %d)", queries, queryLimit)
			}
			ans := 0
			if tc.dir[u][v] == 1 {
				ans = 1
			}
			fmt.Fprintf(writer, "%d\n", ans)
			writer.Flush()
		case "!":
			x, err := readIntToken(reader)
			if err != nil {
				proc.cmd.Process.Kill()
				proc.cmd.Wait()
				return fmt.Errorf("failed to read answer: %v", err)
			}
			writer.Flush()
			proc.stdin.Close()
			if waitErr := proc.cmd.Wait(); waitErr != nil {
				return fmt.Errorf("process exit error: %v (stderr: %s)", waitErr, proc.stderr.String())
			}
			if !tc.valid[x] {
				return fmt.Errorf("reported node %d is invalid", x)
			}
			return nil
		default:
			proc.cmd.Process.Kill()
			proc.cmd.Wait()
			return fmt.Errorf("unexpected token %q", token)
		}
	}
}

func startProcess(target string) (*solverProcess, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &solverProcess{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: &stderr,
	}, nil
}

func readToken(r *bufio.Reader) (string, error) {
	var sb strings.Builder
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if ch <= ' ' {
			continue
		}
		sb.WriteByte(ch)
		break
	}
	for {
		ch, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return sb.String(), nil
			}
			return "", err
		}
		if ch <= ' ' {
			break
		}
		sb.WriteByte(ch)
	}
	return sb.String(), nil
}

func readIntToken(r *bufio.Reader) (int, error) {
	tok, err := readToken(r)
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(tok)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", tok)
	}
	return val, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
