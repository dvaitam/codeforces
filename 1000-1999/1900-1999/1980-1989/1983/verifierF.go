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

type node struct {
	child [2]*node
}

func (n *node) insert(x int) {
	cur := n
	for i := 30; i >= 0; i-- {
		b := (x >> i) & 1
		if cur.child[b] == nil {
			cur.child[b] = &node{}
		}
		cur = cur.child[b]
	}
}

func (n *node) minXor(x int) int {
	cur := n
	res := 0
	for i := 30; i >= 0; i-- {
		b := (x >> i) & 1
		if cur.child[b] != nil {
			cur = cur.child[b]
		} else {
			res |= 1 << i
			cur = cur.child[1-b]
		}
	}
	return res
}

func expected(n, k int, a []int) string {
	var vals []int
	for l := 0; l < n; l++ {
		root := &node{}
		root.insert(a[l])
		minVal := int(^uint(0) >> 1)
		for r := l + 1; r < n; r++ {
			v := root.minXor(a[r])
			if v < minVal {
				minVal = v
			}
			root.insert(a[r])
			vals = append(vals, minVal)
		}
	}
	sort.Ints(vals)
	if k <= len(vals) {
		return fmt.Sprintf("%d", vals[k-1])
	}
	return "-1"
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Println("failed to open testcasesF.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid testcases file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Printf("missing n for case %d\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Printf("missing k for case %d\n", caseNum)
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			arr[i], _ = strconv.Atoi(scan.Text())
		}
		expect := expected(n, k, arr)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		sb.WriteByte('\n')
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %q got %q\n", caseNum, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
