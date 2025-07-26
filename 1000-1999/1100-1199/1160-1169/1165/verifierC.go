package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const numTestsC = 100

func solveC(s string) (int, string) {
	ans := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if len(ans)%2 == 0 {
			ans = append(ans, s[i])
		} else if ans[len(ans)-1] != s[i] {
			ans = append(ans, s[i])
		}
	}
	if len(ans)%2 == 1 {
		ans = ans[:len(ans)-1]
	}
	deletions := len(s) - len(ans)
	return deletions, string(ans)
}

func run(binary, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(3)
	letters := []byte("abc")
	for t := 1; t <= numTestsC; t++ {
		n := rand.Intn(50) + 1
		b := make([]byte, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		s := string(b)
		input := fmt.Sprintf("%d\n%s\n", n, s)
		del, res := solveC(s)
		out, err := run(binary, input)
		if err != nil {
			fmt.Printf("test %d failed to run: %v\noutput:%s\n", t, err, out)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) == 0 {
			fmt.Printf("test %d: no output\n", t)
			os.Exit(1)
		}
		var gotDel int
		fmt.Sscanf(strings.Fields(lines[0])[0], "%d", &gotDel)
		gotRes := ""
		if len(lines) > 1 {
			gotRes = strings.TrimSpace(lines[1])
		}
		if gotDel != del || gotRes != res {
			fmt.Printf("test %d failed\ninput:%sexpected:%d %s\noutput:%s\n", t, input, del, res, out)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
