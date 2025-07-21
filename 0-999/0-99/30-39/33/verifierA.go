package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveA(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return ""
	}
	caps := make([]int64, m+1)
	const inf = int64(1) << 60
	for i := 1; i <= m; i++ {
		caps[i] = inf
	}
	for i := 0; i < n; i++ {
		var r int
		var c int64
		fmt.Fscan(reader, &r, &c)
		if c < caps[r] {
			caps[r] = c
		}
	}
	var total int64
	for i := 1; i <= m; i++ {
		total += caps[i]
	}
	if total > k {
		total = k
	}
	return fmt.Sprintf("%d\n", total)
}

func genTestA() (string, string) {
	m := rand.Intn(5) + 1
	n := m + rand.Intn(5)
	k := rand.Int63n(100)
	type tooth struct {
		r int
		c int64
	}
	teeth := make([]tooth, n)
	for i := 0; i < m; i++ {
		teeth[i].r = i + 1
		teeth[i].c = rand.Int63n(20)
	}
	for i := m; i < n; i++ {
		teeth[i].r = rand.Intn(m) + 1
		teeth[i].c = rand.Int63n(20)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", teeth[i].r, teeth[i].c)
	}
	input := sb.String()
	out := solveA(input)
	return input, out
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go <binary>")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genTestA()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%s\nOutput:\n%s\n", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
