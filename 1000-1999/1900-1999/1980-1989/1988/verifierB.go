package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func generate() (string, string) {
	const T = 100
	var in strings.Builder
	var out strings.Builder
	fmt.Fprintf(&in, "%d\n", T)
	rand.Seed(2)
	for i := 0; i < T; i++ {
		n := rand.Intn(20) + 1
		sb := make([]byte, n)
		ones, groups := 0, 0
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				sb[j] = '0'
				if j == 0 || sb[j-1] != '0' {
					groups++
				}
			} else {
				sb[j] = '1'
				ones++
			}
		}
		s := string(sb)
		fmt.Fprintf(&in, "%d %s\n", n, s)
		if ones > groups {
			fmt.Fprintln(&out, "YES")
		} else {
			fmt.Fprintln(&out, "NO")
		}
	}
	return in.String(), out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	in, exp := generate()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("error running binary:", err)
		os.Exit(1)
	}
	got := buf.String()
	if strings.TrimSpace(got) != strings.TrimSpace(exp) {
		fmt.Println("wrong answer")
		fmt.Println("expected:\n" + exp)
		fmt.Println("got:\n" + got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
