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

type node struct {
	child [2]int
}

func expected(arr []uint64) uint64 {
	n := len(arr)
	P := make([]uint64, n+1)
	for i := 1; i <= n; i++ {
		P[i] = P[i-1] ^ arr[i-1]
	}
	S := make([]uint64, n+2)
	for i := n - 1; i >= 0; i-- {
		S[i+1] = S[i+2] ^ arr[i]
	}
	nodes := make([]node, 1, (n+1)*41)
	insert := func(x uint64) {
		cur := 0
		for bit := 39; bit >= 0; bit-- {
			b := int((x >> uint(bit)) & 1)
			if nodes[cur].child[b] == 0 {
				nodes = append(nodes, node{})
				nodes[cur].child[b] = len(nodes) - 1
			}
			cur = nodes[cur].child[b]
		}
	}
	query := func(x uint64) uint64 {
		cur := 0
		var res uint64
		for bit := 39; bit >= 0; bit-- {
			b := int((x >> uint(bit)) & 1)
			opp := b ^ 1
			if nodes[cur].child[opp] != 0 {
				res |= 1 << uint(bit)
				cur = nodes[cur].child[opp]
			} else {
				cur = nodes[cur].child[b]
			}
		}
		return res
	}
	insert(P[0])
	var ans uint64
	for j := 1; j <= n+1; j++ {
		val := query(S[j])
		if val > ans {
			ans = val
		}
		if j <= n {
			insert(P[j])
		}
	}
	return ans
}

func runCase(bin string, line string, idx int) error {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("test %d: invalid n", idx)
	}
	if len(fields) != 1+n {
		return fmt.Errorf("test %d: expected %d numbers got %d", idx, 1+n, len(fields))
	}
	arr := make([]uint64, n)
	for i := 0; i < n; i++ {
		v, _ := strconv.Atoi(fields[1+i])
		arr[i] = uint64(v)
	}
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", n)
	for i, v := range arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", v)
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var ans uint64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &ans); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(arr)
	if ans != exp {
		return fmt.Errorf("expected %d got %d", exp, ans)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open testcasesE.txt: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		if err := runCase(bin, line, idx); err != nil {
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
