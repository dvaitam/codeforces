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

func solveD(n int, a, b []int) (int, []int) {
	count := 0
	var ans []int
	for pos0 := 0; pos0 < n; pos0++ {
		perm := make([]int, n)
		used := make([]bool, n)
		ok := true
		for i := 0; i < n; i++ {
			v := a[i] ^ pos0
			if v < 0 || v >= n || used[v] {
				ok = false
				break
			}
			perm[i] = v
			used[v] = true
		}
		if !ok {
			continue
		}
		p0 := perm[0]
		inv := make([]int, n)
		for j := 0; j < n; j++ {
			inv[j] = b[j] ^ p0
		}
		for i := 0; i < n && ok; i++ {
			if inv[perm[i]] != i {
				ok = false
			}
		}
		if ok {
			count++
			if ans == nil {
				ans = perm
			}
		}
	}
	return count, ans
}

func runCase(exe, input string, expect string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	filtered := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if strings.HasPrefix(ln, "?") || ln == "" {
			continue
		}
		filtered = append(filtered, ln)
	}
	got := strings.Join(filtered, "\n")
	exp := strings.TrimSpace(expect)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		a := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			a[i], _ = strconv.Atoi(scan.Text())
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			b[i], _ = strconv.Atoi(scan.Text())
		}
		inputBuilder := &strings.Builder{}
		fmt.Fprintf(inputBuilder, "%d\n", n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(inputBuilder, "%d\n", a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fprintf(inputBuilder, "%d\n", b[i])
		}
		cnt, perm := solveD(n, a, b)
		expectBuilder := &strings.Builder{}
		fmt.Fprintln(expectBuilder, "!")
		fmt.Fprintln(expectBuilder, cnt)
		for i := 0; i < n; i++ {
			if i > 0 {
				expectBuilder.WriteByte(' ')
			}
			fmt.Fprintf(expectBuilder, "%d", perm[i])
		}
		expectBuilder.WriteByte('\n')
		if err := runCase(exe, inputBuilder.String(), expectBuilder.String()); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
