package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	n   int
	adj [][]int
}

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func possibleKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func main() {
	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(bytes.NewReader(inputData))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read t: %v\n", err)
		os.Exit(1)
	}

	cases := make([]testCase, t)
	for idx := 0; idx < t; idx++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to read n: %v\n", idx+1, err)
			os.Exit(1)
		}
		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		cases[idx] = testCase{n: n, adj: adj}
	}

	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	candOut, err := runProgram(target, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	outReader := bufio.NewReader(strings.NewReader(candOut))

	for idx, tc := range cases {
		var k int
		if _, err := fmt.Fscan(outReader, &k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to read k: %v\n", idx+1, err)
			os.Exit(1)
		}
		if k < 0 || k > 3*tc.n {
			fmt.Fprintf(os.Stderr, "case %d: invalid number of instructions %d (limit %d)\n", idx+1, k, 3*tc.n)
			os.Exit(1)
		}
		destroyed := make([]bool, tc.n+1)
		possible := map[int]bool{1: true}
		prevType := 0

		for step := 0; step < k; step++ {
			var typ int
			if _, err := fmt.Fscan(outReader, &typ); err != nil {
				fmt.Fprintf(os.Stderr, "case %d: failed to read instruction type at step %d: %v\n", idx+1, step+1, err)
				os.Exit(1)
			}
			if typ == 1 {
				newPossible := make(map[int]bool)
				for node := range possible {
					if destroyed[node] {
						fmt.Fprintf(os.Stderr, "case %d: cat could be on destroyed node %d\n", idx+1, node)
						os.Exit(1)
					}
					moved := false
					for _, nb := range tc.adj[node] {
						if !destroyed[nb] {
							moved = true
							newPossible[nb] = true
						}
					}
					if !moved {
						newPossible[node] = true
					}
				}
				possible = newPossible
				prevType = 1
			} else if typ == 2 {
				var u int
				if _, err := fmt.Fscan(outReader, &u); err != nil {
					fmt.Fprintf(os.Stderr, "case %d: failed to read node for delete at step %d: %v\n", idx+1, step+1, err)
					os.Exit(1)
				}
				if prevType == 2 {
					fmt.Fprintf(os.Stderr, "case %d: consecutive delete instructions at step %d\n", idx+1, step+1)
					os.Exit(1)
				}
				if u < 1 || u > tc.n {
					fmt.Fprintf(os.Stderr, "case %d: delete node %d out of range\n", idx+1, u)
					os.Exit(1)
				}
				if !destroyed[u] {
					if possible[u] {
						fmt.Fprintf(os.Stderr, "case %d: cat might be on node %d when it is deleted\n", idx+1, u)
						os.Exit(1)
					}
					destroyed[u] = true
				}
				prevType = 2
			} else {
				fmt.Fprintf(os.Stderr, "case %d: invalid instruction type %d at step %d\n", idx+1, typ, step+1)
				os.Exit(1)
			}
		}

		if len(possible) != 1 || !possible[tc.n] {
			fmt.Fprintf(os.Stderr, "case %d: cat may end on nodes %v instead of %d\n", idx+1, possibleKeys(possible), tc.n)
			os.Exit(1)
		}
		if destroyed[tc.n] {
			fmt.Fprintf(os.Stderr, "case %d: node %d was destroyed\n", idx+1, tc.n)
			os.Exit(1)
		}
	}

	fmt.Println("all tests passed")
}
