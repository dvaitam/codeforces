package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveD(n int, w []int64, edges [][2]int) string {
	deg := make([]int, n+1)
	for _, e := range edges {
		deg[e[0]]++
		deg[e[1]]++
	}
	total := int64(0)
	for i := 1; i <= n; i++ {
		total += w[i-1]
	}
	extras := make([]int64, 0)
	for i := 1; i <= n; i++ {
		for j := 1; j < deg[i]; j++ {
			extras = append(extras, w[i-1])
		}
	}
	// sort descending
	for i := 0; i < len(extras); i++ {
		for j := i + 1; j < len(extras); j++ {
			if extras[j] > extras[i] {
				extras[i], extras[j] = extras[j], extras[i]
			}
		}
	}
	ans := make([]int64, n)
	ans[1] = total
	for k := 2; k < n; k++ {
		ans[k] = ans[k-1] + extras[k-2]
	}
	res := strings.Builder{}
	for k := 1; k < n; k++ {
		if k > 1 {
			res.WriteByte(' ')
		}
		res.WriteString(fmt.Sprint(ans[k]))
	}
	return res.String()
}

func genCases() []string {
	rand.Seed(4)
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 2
		w := make([]int64, n)
		for j := 0; j < n; j++ {
			w[j] = int64(rand.Intn(10) + 1)
		}
		// build random tree edges
		edges := make([][2]int, n-1)
		for j := 1; j < n; j++ {
			p := rand.Intn(j) + 1
			edges[j-1] = [2]int{p, j + 1}
		}
		sb := strings.Builder{}
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(w[j]))
		}
		sb.WriteByte('\n')
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		cases[i] = sb.String()
	}
	return cases
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var n int
		fmt.Sscan(lines[1], &n)
		wFields := strings.Fields(lines[2])
		w := make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Sscan(wFields[j], &w[j])
		}
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			fmt.Sscan(lines[3+j], &edges[j][0], &edges[j][1])
		}
		want := solveD(n, w, edges)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected: %s Got: %s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
