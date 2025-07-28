package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type edge struct {
	to int
	w  int64
}

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func solveCase(n int, edgesList [][3]int64) string {
	adj := make([][]edge, n+1)
	for _, e := range edgesList {
		a := int(e[0])
		b := int(e[1])
		d := e[2]
		adj[b] = append(adj[b], edge{to: a, w: d})
		adj[a] = append(adj[a], edge{to: b, w: -d})
	}
	pos := make([]int64, n+1)
	vis := make([]bool, n+1)
	queue := make([]int, 0)
	ok := true
	for i := 1; i <= n && ok; i++ {
		if !vis[i] {
			vis[i] = true
			pos[i] = 0
			queue = append(queue, i)
			for len(queue) > 0 && ok {
				u := queue[0]
				queue = queue[1:]
				for _, e := range adj[u] {
					v := e.to
					val := pos[u] + e.w
					if !vis[v] {
						vis[v] = true
						pos[v] = val
						queue = append(queue, v)
					} else if pos[v] != val {
						ok = false
						break
					}
				}
			}
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesH.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var inputs []string
	var exps []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		edgesList := make([][3]int64, m)
		idx := 2
		for i := 0; i < m; i++ {
			a, _ := strconv.Atoi(parts[idx])
			b, _ := strconv.Atoi(parts[idx+1])
			d, _ := strconv.Atoi(parts[idx+2])
			edgesList[i] = [3]int64{int64(a), int64(b), int64(d)}
			idx += 3
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n", n, m)
		for i := 0; i < m; i++ {
			fmt.Fprintf(&sb, "%d %d %d\n", edgesList[i][0], edgesList[i][1], edgesList[i][2])
		}
		inputs = append(inputs, sb.String())
		exps = append(exps, solveCase(n, edgesList))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}

	for idx, input := range inputs {
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		got := strings.TrimSpace(strings.ToUpper(out))
		if got != exps[idx] {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", idx+1, input, exps[idx], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
