package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveB(s string) string {
	var lex []string
	n := len(s)
	for i := 0; i < n; {
		if s[i] == ' ' {
			i++
			continue
		}
		if s[i] == '"' {
			j := i + 1
			for j < n && s[j] != '"' {
				j++
			}
			lex = append(lex, s[i+1:j])
			i = j + 1
		} else {
			j := i
			for j < n && s[j] != ' ' {
				j++
			}
			lex = append(lex, s[i:j])
			i = j
		}
	}
	var sb strings.Builder
	for _, t := range lex {
		sb.WriteByte('<')
		sb.WriteString(t)
		sb.WriteString("\n>")
	}
	return strings.ReplaceAll(sb.String(), "\n>", ">\n")
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,?!")

func randomWord(rng *rand.Rand, l int) string {
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func genCaseB(rng *rand.Rand) (string, string) {
	m := rng.Intn(5) + 1
	parts := make([]string, m)
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			l := rng.Intn(5) + 1
			parts[i] = randomWord(rng, l)
		} else {
			l := rng.Intn(4)
			var sb strings.Builder
			for j := 0; j < l; j++ {
				if j > 0 && rng.Intn(3) == 0 {
					sb.WriteByte(' ')
				}
				sb.WriteByte(letters[rng.Intn(len(letters))])
			}
			parts[i] = "\"" + sb.String() + "\""
		}
	}
	cmd := strings.Join(parts, " ")
	return cmd + "\n", solveB(cmd)
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCaseB(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
