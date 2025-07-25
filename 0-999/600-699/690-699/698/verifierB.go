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

func runCandidate(bin, input string) (string, error) {
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

func solve(a []int) (int, []int) {
	n := len(a)
	parent := make([]int, n+1)
	size := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(x, y int) {
		px := find(x)
		py := find(y)
		if px == py {
			return
		}
		if size[px] > size[py] {
			parent[py] = px
			size[px] += size[py]
		} else {
			parent[px] = py
			size[py] += size[px]
		}
	}
	root := -1
	selfloops := make([]int, 0)
	for i := 1; i <= n; i++ {
		if a[i-1] == i && root == -1 {
			root = i
		} else {
			u := find(i)
			v := find(a[i-1])
			if u == v {
				selfloops = append(selfloops, i)
			} else {
				union(u, v)
			}
		}
	}
	changes := len(selfloops)
	if root == -1 {
		root = selfloops[len(selfloops)-1]
		selfloops = selfloops[:len(selfloops)-1]
	}
	a[root-1] = root
	for _, idx := range selfloops {
		a[idx-1] = root
	}
	return changes, a
}

func expected(a []int) string {
	c, arr := solve(a)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(c))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	file, err := os.Open("testcasesB.txt")
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
		n, err := strconv.Atoi(fields[0])
		if err != nil || n != len(fields)-1 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[i+1])
			nums[i] = v
		}
		want := expected(append([]int(nil), nums...))
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(fields[1:], " "))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
