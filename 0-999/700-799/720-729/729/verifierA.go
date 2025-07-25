package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
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

func expected(n int, s string) string {
	var result []rune
	for i := 0; i < len(s); {
		if i+3 <= len(s) && s[i:i+3] == "ogo" {
			j := i + 3
			for j+1 < len(s) && s[j:j+2] == "go" {
				j += 2
			}
			result = append(result, '*', '*', '*')
			i = j
		} else {
			result = append(result, rune(s[i]))
			i++
		}
	}
	return string(result)
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	r := rand.New(rand.NewSource(1))
	for tc := 1; tc <= 100; tc++ {
		n := r.Intn(100) + 1
		letters := "abgoxyz"
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(letters[r.Intn(len(letters))])
		}
		s := sb.String()
		input := fmt.Sprintf("%d\n%s\n", n, s)
		expect := expected(n, s)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", tc, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %q\ngot: %q\n", tc, input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
