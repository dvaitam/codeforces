package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TestCase struct {
	n      int
	parent []int
	events []string
}

func parseCases(filename string) ([]TestCase, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var cases []TestCase
	for {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		n, _ := strconv.Atoi(line)
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing parents line")
		}
		parentParts := strings.Fields(scanner.Text())
		parent := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parentParts[i])
			parent[i] = v
		}
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing m line")
		}
		m, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		events := make([]string, 0, m)
		for i := 0; i < m; i++ {
			if !scanner.Scan() {
				return nil, fmt.Errorf("missing event line")
			}
			events = append(events, strings.TrimSpace(scanner.Text()))
		}
		cases = append(cases, TestCase{n: n, parent: parent, events: events})
		scanner.Scan() // consume blank line
	}
	return cases, nil
}

func buildAdj(n int, parent []int) [][]int {
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		if parent[i] > 0 {
			p := parent[i] - 1
			adj[i] = append(adj[i], p)
			adj[p] = append(adj[p], i)
		}
	}
	return adj
}

func path(adj [][]int, a, b int) []int {
	n := len(adj)
	prev := make([]int, n)
	for i := range prev {
		prev[i] = -1
	}
	q := []int{a}
	prev[a] = a
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		if v == b {
			break
		}
		for _, to := range adj[v] {
			if prev[to] == -1 {
				prev[to] = v
				q = append(q, to)
			}
		}
	}
	var path []int
	cur := b
	for cur != a {
		path = append(path, cur)
		cur = prev[cur]
	}
	path = append(path, a)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// expectedResults runs the reference solution 226E on the given test case to
// obtain the correct answers.
func expectedResults(tc TestCase) ([]int, error) {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.parent {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')
	input.WriteString(fmt.Sprintf("%d\n", len(tc.events)))
	for _, e := range tc.events {
		input.WriteString(e)
		input.WriteByte('\n')
	}
	cmd := exec.Command("./226E")
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ref run error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	res := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("bad ref output: %v", err)
		}
		res[i] = v
	}
	return res, nil
}

func runCase(bin string, tc TestCase, idx int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.parent {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')
	input.WriteString(fmt.Sprintf("%d\n", len(tc.events)))
	for _, e := range tc.events {
		input.WriteString(e)
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	exp, err := expectedResults(tc)
	if err != nil {
		return err
	}
	if len(fields) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(fields))
	}
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if v != exp[i] {
			return fmt.Errorf("at pos %d expected %d got %d", i, exp[i], v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		if err := runCase(bin, tc, i+1); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
