package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ===== Embedded reference solver for 683F =====

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func refSolve(input string) string {
	s := strings.TrimRight(input, "\r\n")

	var b strings.Builder
	firstWord := true
	capNext := true

	for i := 0; i < len(s); {
		for i < len(s) && s[i] == ' ' {
			i++
		}
		if i >= len(s) {
			break
		}

		if isLetter(s[i]) {
			start := i
			for i < len(s) && isLetter(s[i]) {
				i++
			}
			word := []byte(strings.ToLower(s[start:i]))
			if capNext && len(word) > 0 && 'a' <= word[0] && word[0] <= 'z' {
				word[0] = word[0] - 'a' + 'A'
			}
			if !firstWord {
				b.WriteByte(' ')
			}
			b.Write(word)
			firstWord = false
			capNext = false
		} else if s[i] == '.' || s[i] == ',' {
			if !firstWord {
				b.WriteByte(s[i])
			}
			if s[i] == '.' {
				capNext = true
			}
			i++
		} else {
			i++
		}
	}

	return b.String()
}

// ===== Verifier infrastructure =====

func runProg(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
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
		expect := refSolve(in)
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
