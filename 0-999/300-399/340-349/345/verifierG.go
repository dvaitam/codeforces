package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Node struct {
	children map[byte]*Node
	count    int
}

func insert(root *Node, s string) int {
	cur := root
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if cur.children == nil {
			cur.children = make(map[byte]*Node)
		}
		if cur.children[c] == nil {
			cur.children[c] = &Node{}
		}
		cur = cur.children[c]
		cur.count++
	}
	return cur.count
}

func solveCase(stringsIn []string) int {
	root := &Node{}
	ans := 0
	for _, s := range stringsIn {
		cnt := insert(root, s)
		if cnt > ans {
			ans = cnt
		}
	}
	return ans
}

func runCase(bin string, arr []string) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for _, s := range arr {
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	got, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solveCase(arr)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func randomString(rng *rand.Rand) string {
	l := rng.Intn(10) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		arr := make([]string, n)
		for j := 0; j < n; j++ {
			arr[j] = randomString(rng)
		}
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
