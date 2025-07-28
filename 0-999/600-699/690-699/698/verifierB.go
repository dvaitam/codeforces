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

func validateOutput(orig []int, out string, expectChanges int) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		return fmt.Errorf("output should contain two lines")
	}
	c, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("invalid first line: %v", err)
	}
	if c != expectChanges {
		return fmt.Errorf("reported %d changes, expected %d", c, expectChanges)
	}
	fields := strings.Fields(lines[1])
	if len(fields) != len(orig) {
		return fmt.Errorf("expected %d numbers, got %d", len(orig), len(fields))
	}
	b := make([]int, len(orig))
	diff := 0
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid integer on second line: %v", err)
		}
		if v < 1 || v > len(orig) {
			return fmt.Errorf("value %d out of range", v)
		}
		b[i] = v
		if v != orig[i] {
			diff++
		}
	}
	if diff != c {
		return fmt.Errorf("reported %d changes but sequence differs in %d positions", c, diff)
	}
	root := -1
	for i, v := range b {
		if v == i+1 {
			if root != -1 {
				return fmt.Errorf("multiple roots")
			}
			root = i + 1
		}
	}
	if root == -1 {
		return fmt.Errorf("no root")
	}
	n := len(b)
	for i := 1; i <= n; i++ {
		vis := make(map[int]bool)
		v := i
		for !vis[v] {
			vis[v] = true
			v = b[v-1]
		}
		if v != root {
			return fmt.Errorf("cycle not ending at root")
		}
	}
	return nil
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
		minChanges, _ := solve(append([]int(nil), nums...))
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(fields[1:], " "))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if err := validateOutput(nums, got, minChanges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
