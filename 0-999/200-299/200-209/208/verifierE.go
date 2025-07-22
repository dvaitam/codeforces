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

func expectedAnswer(n int, parents []int, queries [][2]int) []int {
	children := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		p := parents[i-1]
		if p > 0 {
			children[p] = append(children[p], i)
		}
	}
	depth := make([]int, n+1)
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	depthList := make([][]int, n+2)
	t := 0
	var dfs func(int)
	dfs = func(u int) {
		t++
		tin[u] = t
		d := depth[u]
		if d >= len(depthList) {
			depthList = append(depthList, []int{})
		}
		depthList[d] = append(depthList[d], t)
		for _, v := range children[u] {
			depth[v] = d + 1
			dfs(v)
		}
		tout[u] = t
	}
	for i := 1; i <= n; i++ {
		if parents[i-1] == 0 {
			depth[i] = 0
			dfs(i)
		}
	}
	const LOG = 18
	up := make([][]int, LOG)
	up[0] = make([]int, n+1)
	for i := 1; i <= n; i++ {
		up[0][i] = parents[i-1]
	}
	for k := 1; k < LOG; k++ {
		up[k] = make([]int, n+1)
		for i := 1; i <= n; i++ {
			up[k][i] = up[k-1][up[k-1][i]]
		}
	}
	res := make([]int, len(queries))
	for qi, q := range queries {
		v, p := q[0], q[1]
		u := v
		for k := 0; k < LOG && u > 0; k++ {
			if (p>>k)&1 == 1 {
				u = up[k][u]
			}
		}
		if u == 0 {
			res[qi] = 0
			continue
		}
		D := depth[u] + p
		if D >= len(depthList) {
			res[qi] = 0
			continue
		}
		arr := depthList[D]
		l := sort.Search(len(arr), func(i int) bool { return arr[i] >= tin[u] })
		r := sort.Search(len(arr), func(i int) bool { return arr[i] > tout[u] })
		cnt := r - l
		if cnt > 0 {
			cnt--
		}
		res[qi] = cnt
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		pos := 0
		n, _ := strconv.Atoi(fields[pos])
		pos++
		if len(fields) < pos+n+1 {
			fmt.Printf("test %d: invalid parent count\n", idx)
			os.Exit(1)
		}
		parents := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[pos+i])
			parents[i] = v
		}
		pos += n
		m, _ := strconv.Atoi(fields[pos])
		pos++
		if len(fields) != pos+2*m {
			fmt.Printf("test %d: expected %d queries got %d\n", idx, m, (len(fields)-pos)/2)
			os.Exit(1)
		}
		queries := make([][2]int, m)
		for i := 0; i < m; i++ {
			v, _ := strconv.Atoi(fields[pos+2*i])
			p, _ := strconv.Atoi(fields[pos+2*i+1])
			queries[i] = [2]int{v, p}
		}
		answers := expectedAnswer(n, parents, queries)
		var sb strings.Builder
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for i, v := range parents {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprint(m))
		sb.WriteByte('\n')
		for _, q := range queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.Fields(strings.TrimSpace(out.String()))
		expectAns := make([]string, len(answers))
		for i, v := range answers {
			expectAns[i] = fmt.Sprint(v)
		}
		if strings.Join(got, " ") != strings.Join(expectAns, " ") {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, strings.Join(expectAns, " "), strings.Join(got, " "))
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
