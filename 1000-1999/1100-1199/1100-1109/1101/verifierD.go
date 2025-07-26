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

type caseD struct {
	n     int
	a     []int
	edges [][2]int
}

func readTestcasesD(path string) ([]caseD, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("bad file")
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([]caseD, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("bad file")
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			arr[j], _ = strconv.Atoi(scan.Text())
		}
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			scan.Scan()
			u, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			edges[j] = [2]int{u - 1, v - 1}
		}
		cases[i] = caseD{n, arr, edges}
	}
	return cases, nil
}

func solveCaseD(cs caseD) int {
	n := cs.n
	a := cs.a
	ma := 0
	for _, x := range a {
		if x > ma {
			ma = x
		}
	}
	g0 := make([][]int, ma+1)
	for i := 2; i <= ma; i++ {
		for j := i; j <= ma; j += i {
			g0[j] = append(g0[j], i)
		}
	}
	divs := make([][]int, n)
	divSet := make([]map[int]bool, n)
	for i := 0; i < n; i++ {
		divs[i] = g0[a[i]]
		m := make(map[int]bool, len(divs[i]))
		for _, p := range divs[i] {
			m[p] = true
		}
		divSet[i] = m
	}
	adj := make([][]int, n)
	for _, e := range cs.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	f1 := make([]map[int]int, n)
	f2 := make([]map[int]int, n)
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}
	type nodeState struct {
		u       int
		visited bool
	}
	stack := []nodeState{{0, false}}
	parent[0] = -2
	ans := 0
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		u := cur.u
		if !cur.visited {
			stack = append(stack, nodeState{u, true})
			for _, v := range adj[u] {
				if parent[v] == -1 {
					parent[v] = u
					stack = append(stack, nodeState{v, false})
				}
			}
		} else {
			f1[u] = make(map[int]int)
			f2[u] = make(map[int]int)
			for _, v := range adj[u] {
				if parent[v] != u {
					continue
				}
				for _, p := range divs[v] {
					if !divSet[u][p] {
						continue
					}
					l := f1[v][p] + 1
					if f1[u][p] < l {
						f2[u][p] = f1[u][p]
						f1[u][p] = l
					} else if f2[u][p] < l {
						f2[u][p] = l
					}
				}
			}
			for _, p := range divs[u] {
				v1 := f1[u][p]
				v2 := f2[u][p]
				if v1+v2+1 > ans {
					ans = v1 + v2 + 1
				}
			}
		}
	}
	return ans
}

func run(bin string, input string) (string, error) {
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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) > 2 {
		bin = os.Args[2]
	}
	cases, err := readTestcasesD("testcasesD.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for idx, cs := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", cs.n))
		for i, v := range cs.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, e := range cs.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
		}
		input := sb.String()
		expected := solveCaseD(cs)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output: %v", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", idx+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
