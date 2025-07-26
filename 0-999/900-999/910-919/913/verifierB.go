package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type treeInput struct {
	n       int
	parents []int
}

func randomTree(n int) []int {
	for {
		parents := make([]int, n-1)
		childCount := make([]int, n+1)
		for i := 2; i <= n; i++ {
			p := rand.Intn(i-1) + 1
			parents[i-2] = p
			childCount[p]++
		}
		if childCount[1] >= 2 {
			return parents
		}
	}
}

func solve(n int, parents []int) string {
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parents[i-2]
		children[p] = append(children[p], i)
	}
	for v := 1; v <= n; v++ {
		if len(children[v]) == 0 {
			continue
		}
		leafChildren := 0
		for _, ch := range children[v] {
			if len(children[ch]) == 0 {
				leafChildren++
			}
		}
		if leafChildren < 3 {
			return "No"
		}
	}
	return "Yes"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	const cases = 100
	for i := 0; i < cases; i++ {
		n := rand.Intn(20) + 3
		parents := randomTree(n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, p := range parents {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", p))
		}
		sb.WriteByte('\n')
		input := sb.String()
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			fmt.Printf("program output:\n%s\n", string(out))
			return
		}
		got := strings.TrimSpace(string(out))
		want := solve(n, parents)
		if got != want {
			fmt.Printf("case %d failed:\ninput:\n%sexpected %s got %s\n", i+1, input, want, got)
			return
		}
	}
	fmt.Printf("OK %d cases\n", cases)
}
