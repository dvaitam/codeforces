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

const MOD = 1000000007

func run(bin, input string) (string, error) {
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

func countInv(p []int) int {
	cnt := 0
	for i := 0; i < len(p); i++ {
		for j := i + 1; j < len(p); j++ {
			if p[i] > p[j] {
				cnt++
			}
		}
	}
	return cnt
}

func leq(a, b []int) bool {
	for i := range a {
		if a[i] < b[i] {
			return true
		}
		if a[i] > b[i] {
			return false
		}
	}
	return true
}

func solveCase(n int, perm []int) int {
	elems := make([]int, n)
	for i := 0; i < n; i++ {
		elems[i] = i + 1
	}
	used := make([]bool, n)
	var cur []int
	var ans int
	var dfs func(pos int)
	dfs = func(pos int) {
		if pos == n {
			if leq(cur, perm) {
				ans += countInv(cur)
			}
			return
		}
		for i := 0; i < n; i++ {
			if !used[i] {
				used[i] = true
				cur = append(cur, elems[i])
				dfs(pos + 1)
				cur = cur[:len(cur)-1]
				used[i] = false
			}
		}
	}
	dfs(0)
	return ans % MOD
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		perm := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			perm[j], _ = strconv.Atoi(scan.Text())
		}
		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		for j, v := range perm {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		expected := solveCase(n, perm)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		gotVal, err2 := strconv.Atoi(strings.TrimSpace(got))
		if err2 != nil || gotVal != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i+1, expected, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
