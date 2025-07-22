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
	var k int
	fmt.Fscan(reader, &k)
	arr := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	ans := make([]int, 0, k)
	var y int
	f, f1 := false, false
	for _, x := range arr {
		if x == 0 || x == 100 {
			ans = append(ans, x)
		} else if x < 10 && !f {
			ans = append(ans, x)
			f = true
		} else if x%10 == 0 && !f1 {
			ans = append(ans, x)
			f1 = true
		} else {
			y = x
		}
	}
	if !f && !f1 && y != 0 {
		ans = append(ans, y)
	}
	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(ans))
	for i, v := range ans {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprint(&buf, v)
	}
	buf.WriteByte('\n')
	return strings.TrimSpace(buf.String())
}

func genTestA() string {
	k := rand.Intn(100) + 1
	used := make(map[int]bool)
	arr := make([]int, 0, k)
	for len(arr) < k {
		x := rand.Intn(101)
		if !used[x] {
			used[x] = true
			arr = append(arr, x)
		}
	}
	var buf bytes.Buffer
	fmt.Fprintln(&buf, k)
	for i, v := range arr {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprint(&buf, v)
	}
	buf.WriteByte('\n')
	return buf.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		in := genTestA()
		expected := solveA(in)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(in)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\noutput: %s\n", i, err, string(out))
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
