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

type Info struct {
	name string
	p    int
	h    int
}

func expected(n int, arr []Info) string {
	sort.Slice(arr, func(i, j int) bool { return arr[i].p < arr[j].p })
	id := make([]bool, n+5)
	for i := n - 1; i >= 0; i-- {
		x := (i - arr[i].p) + 1
		p := 1
		if x <= 0 {
			return "-1"
		}
		for x > 0 || id[p] {
			if !id[p] {
				x--
			}
			if x == 0 {
				break
			}
			p++
		}
		id[p] = true
		arr[i].h = p
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%s %d\n", arr[i].name, arr[i].h))
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Println("failed to open testcasesC.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	tests := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 1 {
			continue
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != 1+2*n {
			fmt.Printf("test %d: invalid line\n", tests+1)
			os.Exit(1)
		}
		arr := make([]Info, n)
		idx := 1
		for i := 0; i < n; i++ {
			arr[i].name = parts[idx]
			idx++
			arr[i].p, _ = strconv.Atoi(parts[idx])
			idx++
		}
		expect := expected(n, append([]Info(nil), arr...))
		// build input
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			input.WriteString(fmt.Sprintf("%s %d\n", arr[i].name, arr[i].p))
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tests+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", tests+1, expect, got)
			os.Exit(1)
		}
		tests++
	}
	if err := sc.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", tests)
}
