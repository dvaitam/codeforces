package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type caseD struct {
	n     int
	edges [][2]int
}

func parseCases(path string) ([]caseD, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	cases := []caseD{}
	for {
		if !sc.Scan() {
			break
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		n, _ := strconv.Atoi(line)
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("bad file")
			}
			var u, v int
			fmt.Sscan(sc.Text(), &u, &v)
			edges[i] = [2]int{u, v}
		}
		cases = append(cases, caseD{n: n, edges: edges})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func encode(adj [][]int, u, p int) string {
	subs := []string{}
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		subs = append(subs, encode(adj, v, u))
	}
	sort.Strings(subs)
	return "(" + strings.Join(subs, "") + ")"
}

func canonical(adj [][]int) string {
	n := len(adj)
	if n == 1 {
		return "()"
	}
	deg := make([]int, n)
	leaves := []int{}
	for i := 0; i < n; i++ {
		deg[i] = len(adj[i])
		if deg[i] <= 1 {
			leaves = append(leaves, i)
		}
	}
	removed := len(leaves)
	for removed < n {
		newLeaves := []int{}
		for _, leaf := range leaves {
			for _, nb := range adj[leaf] {
				deg[nb]--
				if deg[nb] == 1 {
					newLeaves = append(newLeaves, nb)
				}
			}
		}
		removed += len(newLeaves)
		if removed >= n {
			leaves = newLeaves
			break
		}
		leaves = newLeaves
	}
	best := ""
	for _, c := range leaves {
		enc := encode(adj, c, -1)
		if best == "" || enc < best {
			best = enc
		}
	}
	return best
}

func solve(tc caseD) int {
	adj := make([][]int, tc.n)
	deg := make([]int, tc.n)
	for _, e := range tc.edges {
		u := e[0] - 1
		v := e[1] - 1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		deg[u]++
		deg[v]++
	}
	forms := make(map[string]bool)
	for i := 0; i < tc.n; i++ {
		if deg[i] >= 4 {
			continue
		}
		newAdj := make([][]int, tc.n+1)
		for j := 0; j < tc.n; j++ {
			newAdj[j] = append([]int{}, adj[j]...)
		}
		newAdj[i] = append(newAdj[i], tc.n)
		newAdj[tc.n] = []int{i}
		form := canonical(newAdj)
		forms[form] = true
	}
	return len(forms)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases("testcasesD.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		exp := solve(tc)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		val, _ := strconv.Atoi(strings.TrimSpace(out))
		if val != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, exp, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
