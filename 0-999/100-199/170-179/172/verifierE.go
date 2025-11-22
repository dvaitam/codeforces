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

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type node struct {
	tag      string
	children []*node
}

// parse the BHTML document into a tree
func buildTree(doc string) *node {
	root := &node{}
	stack := []*node{root}
	n := len(doc)

	for i := 0; i < n; {
		if doc[i] != '<' {
			i++
			continue
		}

		j := i + 1
		closing := false
		if j < n && doc[j] == '/' {
			closing = true
			j++
		}

		start := j
		for j < n && doc[j] != '>' && doc[j] != '/' {
			j++
		}
		name := doc[start:j]

		selfClose := false
		if !closing && j < n && doc[j] == '/' {
			selfClose = true
			j++
		}
		if j < n && doc[j] == '>' {
			j++
		}

		if closing {
			// pop matching opening tag
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
			i = j
			continue
		}

		cur := &node{tag: name}
		parent := stack[len(stack)-1]
		parent.children = append(parent.children, cur)
		if !selfClose {
			stack = append(stack, cur)
		}
		i = j
	}

	return root
}

// compute expected output using an internal implementation of the rules
func expected(doc string, queries []string) ([]int, error) {
	tree := buildTree(doc)

	// map tags to integer ids for faster comparison
	tagToID := make(map[string]int)
	getID := func(tag string) int {
		if id, ok := tagToID[tag]; ok {
			return id
		}
		id := len(tagToID)
		tagToID[tag] = id
		return id
	}

	qIDs := make([][]int, len(queries))
	for i, q := range queries {
		parts := strings.Fields(q)
		if len(parts) == 0 {
			qIDs[i] = nil
			continue
		}
		qIDs[i] = make([]int, len(parts))
		for j, p := range parts {
			qIDs[i][j] = getID(p)
		}
	}

	results := make([]int, len(queries))

	var dfs func(*node, []int)
	dfs = func(u *node, states []int) {
		var nextStates []int
		if u.tag != "" {
			tagID := getID(u.tag)
			nextStates = make([]int, len(states))
			for i, idx := range states {
				pat := qIDs[i]
				if pat == nil {
					continue
				}
				plen := len(pat)

				// Check if current node matches the target (last element)
				// and we have matched all ancestors (idx == plen-1)
				if idx == plen-1 && pat[plen-1] == tagID {
					results[i]++
				}

				// Update state for children (ancestor matching)
				// Greedy match: if current node matches the expected ancestor, advance index
				if idx < plen-1 && pat[idx] == tagID {
					nextStates[i] = idx + 1
				} else {
					nextStates[i] = idx
				}
			}
		} else {
			// Virtual root, pass states as is
			nextStates = states
		}

		for _, ch := range u.children {
			dfs(ch, nextStates)
		}
	}

	initialStates := make([]int, len(queries))
	dfs(tree, initialStates)

	return results, nil
}

func matchChain(path, query []int) bool {
	// Deprecated, kept/removed as needed. 
	// Since I am replacing the block, I should probably remove it or replace it with nothing if possible.
	// The tool requires replacing exact text. 
	// I will include it in old_string and remove it in new_string if I can capture it.
	// But simpler to just define expected and let matchChain be unused (or remove it if I target it).
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil || len(parts) != m+2 {
			fmt.Fprintf(os.Stderr, "test %d bad m\n", idx)
			os.Exit(1)
		}
		doc := parts[0]
		queries := parts[2:]
		want, err := expected(doc, queries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d reference error: %v\n", idx, err)
			os.Exit(1)
		}
		var input strings.Builder
		input.WriteString(doc)
		input.WriteByte('\n')
		input.WriteString(fmt.Sprintf("%d\n", m))
		for _, q := range queries {
			input.WriteString(q)
			input.WriteByte('\n')
		}
		gotStr, err := runProg(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		scan := bufio.NewScanner(strings.NewReader(gotStr))
		outVals := []int{}
		for scan.Scan() {
			v, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
			outVals = append(outVals, v)
		}
		if len(outVals) != m {
			fmt.Fprintf(os.Stderr, "case %d wrong output length\n", idx)
			os.Exit(1)
		}
		for i := 0; i < m; i++ {
			if outVals[i] != want[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at query %d: expected %d got %d\n", idx, i+1, want[i], outVals[i])
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
