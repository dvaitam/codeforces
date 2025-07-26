package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// parse feed string and return comments by depth
func parseFeed(s string) ([][]string, int) {
	var res [][]string
	maxDepth := 0
	var stack []int
	i := 0
	n := len(s)
	for i < n {
		j := i
		for j < n && s[j] != ',' {
			j++
		}
		text := s[i:j]
		j++
		kStart := j
		for j < n && s[j] != ',' {
			j++
		}
		kStr := s[kStart:j]
		k, _ := strconv.Atoi(kStr)
		if j < n && s[j] == ',' {
			j++
		}
		i = j
		depth := len(stack) + 1
		if depth > maxDepth {
			maxDepth = depth
		}
		for len(res) < depth {
			res = append(res, []string{})
		}
		res[depth-1] = append(res[depth-1], text)
		if len(stack) > 0 {
			stack[len(stack)-1]--
			for len(stack) > 0 && stack[len(stack)-1] == 0 {
				stack = stack[:len(stack)-1]
			}
		}
		if k > 0 {
			stack = append(stack, k)
		}
	}
	return res, maxDepth
}

func encodeTree(tree [][]int, texts []string) string {
	var sb strings.Builder
	var dfs func(int)
	dfs = func(v int) {
		sb.WriteString(texts[v])
		sb.WriteByte(',')
		children := tree[v]
		sb.WriteString(strconv.Itoa(len(children)))
		for _, c := range children {
			sb.WriteByte(',')
			dfs(c)
		}
	}
	dfs(0)
	return sb.String()
}

func randomFeed() string {
	maxNodes := 8
	texts := make([]string, maxNodes)
	for i := range texts {
		texts[i] = fmt.Sprintf("c%d", i)
	}
	tree := make([][]int, maxNodes)
	used := 1
	for v := 0; v < used && used < maxNodes; v++ {
		children := rand.Intn(3)
		for j := 0; j < children && used < maxNodes; j++ {
			tree[v] = append(tree[v], used)
			used++
		}
	}
	return encodeTree(tree[:used], texts)
}

func expectedOutput(feed string) string {
	levels, depth := parseFeed(feed)
	var out strings.Builder
	out.WriteString(fmt.Sprintf("%d\n", depth))
	for d := 0; d < depth; d++ {
		for j, txt := range levels[d] {
			if j > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(txt)
		}
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

func runCase(bin, feed string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(feed + "\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	expected := expectedOutput(feed)
	got := strings.TrimSpace(string(out))
	if expected != got {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	for tcase := 1; tcase <= 120; tcase++ {
		feed := randomFeed()
		if err := runCase(bin, feed); err != nil {
			fmt.Printf("Test %d failed: %v\n", tcase, err)
			return
		}
	}
	fmt.Println("OK")
}
