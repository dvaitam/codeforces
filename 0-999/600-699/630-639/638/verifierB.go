package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solve(frags []string) string {
	var vis [256]int
	var edges [256]rune
	var ans []rune

	for _, s := range frags {
		for j := 0; j+1 < len(s); j++ {
			u := rune(s[j])
			v := rune(s[j+1])
			edges[u] = v
			vis[v] = 2
		}
		u0 := rune(s[0])
		if vis[u0] != 2 {
			vis[u0] = 1
		}
	}

	var dfs func(rune)
	dfs = func(u rune) {
		vis[u] = 3
		v := edges[u]
		if v >= 'a' && vis[v] != 3 {
			dfs(v)
		}
		ans = append(ans, u)
	}

	for u := 'a'; u <= 'z'; u++ {
		if vis[u] == 1 {
			dfs(u)
		}
	}

	for i, j := 0, len(ans)-1; i < j; i, j = i+1, j-1 {
		ans[i], ans[j] = ans[j], ans[i]
	}
	return string(ans)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		panic(err)
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
		parts := strings.Fields(line)
		if len(parts) < 1 {
			continue
		}
		n := 0
		fmt.Sscanf(parts[0], "%d", &n)
		frags := parts[1:]
		if len(frags) < n {
			// handle if generator produced less; but there should be exactly n
		}
		expect := solve(frags[:n])

		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for i := 0; i < n; i++ {
			fmt.Fprintln(&input, frags[i])
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expect {
			fmt.Printf("Test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
