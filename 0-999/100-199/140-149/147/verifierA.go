package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"
)

type token struct {
	word  string
	punct string
}

func randomWord(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func randomPunct(rng *rand.Rand) string {
	p := []string{".", ",", "!", "?"}
	return p[rng.Intn(len(p))]
}

func generateCase(rng *rand.Rand) (string, string) {
	// build tokens
	num := rng.Intn(10) + 1
	toks := make([]token, 0, num)
	toks = append(toks, token{word: randomWord(rng)})
	for len(toks) < num {
		prev := toks[len(toks)-1]
		if prev.punct != "" {
			// must add a word
			toks = append(toks, token{word: randomWord(rng)})
		} else {
			if rng.Intn(3) == 0 {
				toks = append(toks, token{punct: randomPunct(rng)})
			} else {
				toks = append(toks, token{word: randomWord(rng)})
			}
		}
	}
	// ensure last is word
	if toks[len(toks)-1].punct != "" {
		toks[len(toks)-1] = token{word: randomWord(rng)}
	}
	var inputBuilder strings.Builder
	for i, t := range toks {
		if i > 0 {
			for j := 0; j < rng.Intn(3); j++ {
				inputBuilder.WriteByte(' ')
			}
		}
		if t.word != "" {
			inputBuilder.WriteString(t.word)
		} else {
			inputBuilder.WriteString(t.punct)
		}
		if i < len(toks)-1 {
			for j := 0; j < rng.Intn(3); j++ {
				inputBuilder.WriteByte(' ')
			}
		}
	}
	input := inputBuilder.String()
	expect := solveCase(input)
	input += "\n"
	return input, expect
}

func solveCase(line string) string {
	var tokens []string
	var curr strings.Builder
	for _, r := range line {
		if unicode.IsLetter(r) {
			curr.WriteRune(r)
		} else if r == '.' || r == ',' || r == '!' || r == '?' {
			if curr.Len() > 0 {
				tokens = append(tokens, curr.String())
				curr.Reset()
			}
			tokens = append(tokens, string(r))
		} else {
			if curr.Len() > 0 {
				tokens = append(tokens, curr.String())
				curr.Reset()
			}
		}
	}
	if curr.Len() > 0 {
		tokens = append(tokens, curr.String())
	}
	var out strings.Builder
	for i, tok := range tokens {
		isPunct := len(tok) == 1 && strings.ContainsAny(tok, ".,!?")
		if isPunct {
			out.WriteString(tok)
			out.WriteByte(' ')
		} else {
			out.WriteString(tok)
			if i+1 < len(tokens) {
				next := tokens[i+1]
				nextPunct := len(next) == 1 && strings.ContainsAny(next, ".,!?")
				if !nextPunct {
					out.WriteByte(' ')
				}
			}
		}
	}
	return out.String()
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
