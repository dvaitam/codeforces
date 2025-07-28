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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func inversions(perm []int) int {
	inv := 0
	for i := 0; i < len(perm); i++ {
		for j := i + 1; j < len(perm); j++ {
			if perm[i] > perm[j] {
				inv++
			}
		}
	}
	return inv
}

func expected(n, k int, p []int) int {
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i + 1
	}
	best := int(^uint(0) >> 1)
	var permute func(int)
	permute = func(pos int) {
		if pos == n {
			for i := 0; i < n; i++ {
				for j := i + 1; j < n; j++ {
					if p[idx[i]-1] > p[idx[j]-1]+k {
						return
					}
				}
			}
			inv := inversions(idx)
			if inv < best {
				best = inv
			}
			return
		}
		for i := pos; i < n; i++ {
			idx[pos], idx[i] = idx[i], idx[pos]
			permute(pos + 1)
			idx[pos], idx[i] = idx[i], idx[pos]
		}
	}
	permute(0)
	return best
}

func main() {
	arg := ""
	if len(os.Args) == 2 {
		arg = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		arg = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := arg
	file, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesF.txt: %v\n", err)
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
		parts := strings.Fields(line)
		if len(parts) < 3 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		k, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+n {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		p := make([]int, n)
		for i := 0; i < n; i++ {
			p[i], _ = strconv.Atoi(parts[2+i])
		}
		expect := fmt.Sprintf("%d", expected(n, k, p))
		input := fmt.Sprintf("%d %d\n%s\n", n, k, strings.Join(parts[2:], " "))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
