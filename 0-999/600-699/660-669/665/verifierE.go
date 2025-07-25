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

const maxBit = 30

type node struct {
	child [2]*node
	cnt   int
}

func insert(root *node, x int) {
	cur := root
	for b := maxBit; b >= 0; b-- {
		bit := (x >> b) & 1
		if cur.child[bit] == nil {
			cur.child[bit] = &node{}
		}
		cur = cur.child[bit]
		cur.cnt++
	}
}

func countLess(root *node, x, k int) int {
	cur := root
	res := 0
	for b := maxBit; b >= 0; b-- {
		if cur == nil {
			break
		}
		xBit := (x >> b) & 1
		kBit := (k >> b) & 1
		if kBit == 1 {
			if cur.child[xBit] != nil {
				res += cur.child[xBit].cnt
			}
			cur = cur.child[xBit^1]
		} else {
			cur = cur.child[xBit]
		}
	}
	return res
}

func expected(n, k int, arr []int) string {
	prefix := 0
	root := &node{}
	insert(root, 0)
	total := int64(0)
	inserted := 1
	for i := 0; i < n; i++ {
		prefix ^= arr[i]
		less := countLess(root, prefix, k)
		total += int64(inserted - less)
		insert(root, prefix)
		inserted++
	}
	return fmt.Sprintf("%d", total)
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		kVal, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			arr[i], _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, kVal))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(arr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(n, kVal, arr) + "\n"
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
