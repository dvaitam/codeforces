package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "2000-2999/2000-2099/2010-2019/2013/2013F2.go"

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierF2.go /path/to/candidate")
	}
	candidate := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	expectedCount, err := countAnswers(inputData)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(exec.Command(refBin), inputData)
	if err != nil {
		fail("reference solution failed: %v", err)
	}
	expected, err := parseWinners(refOut, expectedCount)
	if err != nil {
		fail("invalid reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), inputData)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	got, err := parseWinners(candOut, expectedCount)
	if err != nil {
		fail("invalid candidate output: %v", err)
	}

	for i := 0; i < expectedCount; i++ {
		if got[i] != expected[i] {
			fail("mismatch at answer %d: expected %s, got %s", i+1, expected[i], got[i])
		}
	}

	fmt.Println("OK")
}

func countAnswers(data []byte) (int, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, err
	}
	total := 0
	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return 0, err
		}
		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var a, b int
			if _, err := fmt.Fscan(reader, &a, &b); err != nil {
				return 0, err
			}
			a--
			b--
			adj[a] = append(adj[a], b)
			adj[b] = append(adj[b], a)
		}
		var u, v int
		if _, err := fmt.Fscan(reader, &u, &v); err != nil {
			return 0, err
		}
		u--
		v--
		dist := distance(adj, u, v)
		if dist < 0 {
			return 0, fmt.Errorf("vertices %d and %d not connected", u+1, v+1)
		}
		total += dist + 1
	}
	return total, nil
}

func distance(adj [][]int, start, target int) int {
	if start == target {
		return 0
	}
	n := len(adj)
	q := make([]int, n)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	head, tail := 0, 0
	q[tail] = start
	tail++
	dist[start] = 0
	for head < tail {
		u := q[head]
		head++
		nd := dist[u] + 1
		for _, v := range adj[u] {
			if dist[v] != -1 {
				continue
			}
			dist[v] = nd
			if v == target {
				return nd
			}
			q[tail] = v
			tail++
		}
	}
	return -1
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2013F2-ref-*")
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

func runProgram(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
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

func parseWinners(out string, expectedCount int) ([]string, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([]string, 0, expectedCount)
	for len(ans) < expectedCount {
		token, err := readToken(reader)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("expected %d answers, got %d", expectedCount, len(ans))
			}
			return nil, err
		}
		low := strings.ToLower(token)
		if low != "alice" && low != "bob" {
			return nil, fmt.Errorf("invalid answer token %q", token)
		}
		ans = append(ans, low)
	}
	if extra, err := readToken(reader); err == nil {
		return nil, fmt.Errorf("unexpected extra output token %q", extra)
	} else if err != io.EOF {
		return nil, err
	}
	return ans, nil
}

func readToken(r *bufio.Reader) (string, error) {
	var b strings.Builder
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if ch > ' ' {
			b.WriteByte(ch)
			break
		}
	}
	for {
		ch, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return b.String(), nil
			}
			return "", err
		}
		if ch <= ' ' {
			return b.String(), nil
		}
		b.WriteByte(ch)
	}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
