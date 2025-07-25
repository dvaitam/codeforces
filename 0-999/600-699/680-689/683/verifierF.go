package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "683F.go")
	return runProg(ref, input)
}

func randWord() string {
	l := rand.Intn(5) + 1
	b := make([]byte, l)
	for i := range b {
		v := rand.Intn(52)
		if v < 26 {
			b[i] = byte('a' + v)
		} else {
			b[i] = byte('A' + v - 26)
		}
	}
	return string(b)
}

func genCase() string {
	tokenCount := rand.Intn(15) + 1
	var tokens []string
	for i := 0; i < tokenCount; i++ {
		if rand.Intn(5) == 0 {
			if rand.Intn(2) == 0 {
				tokens = append(tokens, ".")
			} else {
				tokens = append(tokens, ",")
			}
		} else {
			tokens = append(tokens, randWord())
		}
	}
	var sb strings.Builder
	if rand.Intn(2) == 0 {
		sb.WriteByte(' ')
	}
	for i, t := range tokens {
		sb.WriteString(t)
		if i+1 < len(tokens) {
			sb.WriteString(strings.Repeat(" ", rand.Intn(3)))
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		in := genCase()
		expect, err := runRef(in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:%sexpected:%s\nactual:%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
