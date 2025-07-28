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

func reverseStr(s string) string {
	bs := []byte(s)
	for i, j := 0, len(bs)-1; i < j; i, j = i+1, j-1 {
		bs[i], bs[j] = bs[j], bs[i]
	}
	return string(bs)
}

func solve(words []string) string {
	set2 := make(map[string]bool)
	set3 := make(map[string]bool)
	for _, w := range words {
		if len(w) == 1 {
			return "YES"
		}
		rev := reverseStr(w)
		if len(w) == 2 {
			if w[0] == w[1] {
				return "YES"
			}
			if set2[rev] || set3[rev] {
				return "YES"
			}
			for c := byte('a'); c <= byte('c'); c++ {
				key := rev + string(c)
				if set3[key] {
					return "YES"
				}
			}
			set2[w] = true
		} else if len(w) == 3 {
			if w[0] == w[2] {
				return "YES"
			}
			if set3[rev] {
				return "YES"
			}
			key := rev[:2]
			if set2[key] {
				return "YES"
			}
			set3[w] = true
		}
	}
	return "NO"
}

func buildInput(words []string) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(words)))
	for i, w := range words {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(w)
	}
	sb.WriteByte('\n')
	return sb.String()
}

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

func main() {
	var bin string
	if len(os.Args) == 2 {
		bin = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		bin = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []byte("abc")
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 1
		words := make([]string, n)
		for j := range words {
			l := rng.Intn(3) + 1
			b := make([]byte, l)
			for k := range b {
				b[k] = letters[rng.Intn(len(letters))]
			}
			words[j] = string(b)
		}
		input := buildInput(words)
		exp := solve(words)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("case %d wrong answer\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
