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

func solve(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(r, &n); err != nil {
		return ""
	}
	var g [4][4]int
	for i := 0; i < 4; i++ {
		fmt.Fscan(r, &g[i][0], &g[i][1], &g[i][2], &g[i][3])
	}
	for i := 0; i < 4; i++ {
		m1 := g[i][0]
		if g[i][1] < m1 {
			m1 = g[i][1]
		}
		m2 := g[i][2]
		if g[i][3] < m2 {
			m2 = g[i][3]
		}
		if m1+m2 <= n {
			return fmt.Sprintf("%d %d %d\n", i+1, m1, n-m1)
		}
	}
	return "-1\n"
}

func genTest() (string, string) {
	n := rand.Intn(200) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < 4; i++ {
		a := rand.Intn(200) + 1
		b := rand.Intn(200) + 1
		c := rand.Intn(200) + 1
		d := rand.Intn(200) + 1
		fmt.Fprintf(&sb, "%d %d %d %d\n", a, b, c, d)
	}
	inp := sb.String()
	out := solve(inp)
	return inp, out
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genTest()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%sOutput:\n%s\n", i+1, err, in, got)
			return
		}
		if got != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
