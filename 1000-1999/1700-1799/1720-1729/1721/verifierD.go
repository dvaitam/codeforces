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

func runExe(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(a, b []int) string {
	n := len(a)
	ans := 0
	for bit := 29; bit >= 0; bit-- {
		cand := ans | (1 << uint(bit))
		count := make(map[int]int)
		for i := 0; i < n; i++ {
			count[a[i]&cand]++
		}
		ok := true
		for i := 0; i < n; i++ {
			key := cand ^ (b[i] & cand)
			if count[key] == 0 {
				ok = false
				break
			}
			count[key]--
		}
		if ok {
			ans = cand
		}
	}
	return fmt.Sprintf("%d", ans)
}

func intsToStrs(a []int) []string {
	res := make([]string, len(a))
	for i, v := range a {
		res[i] = strconv.Itoa(v)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
	var t int
	fmt.Sscan(scan.Text(), &t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		var n int
		fmt.Sscan(scan.Text(), &n)
		a := make([]int, n)
		bArr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &a[i])
		}
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &bArr[i])
		}
		exp := solveCase(a, bArr)
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d\n", n)
		fmt.Fprintln(&input, strings.Join(intsToStrs(a), " "))
		fmt.Fprintln(&input, strings.Join(intsToStrs(bArr), " "))
		got, err := runExe(bin, input.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
