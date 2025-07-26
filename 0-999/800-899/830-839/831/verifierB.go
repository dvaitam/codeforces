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

type testCaseB struct {
	layout1 string
	layout2 string
	text    string
}

func randomLayout(rng *rand.Rand) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	rng.Shuffle(len(letters), func(i, j int) { letters[i], letters[j] = letters[j], letters[i] })
	return string(letters)
}

func randomText(rng *rand.Rand, length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(chars[rng.Intn(len(chars))])
	}
	return sb.String()
}

func generateCaseB(rng *rand.Rand) (string, testCaseB) {
	l1 := randomLayout(rng)
	l2 := randomLayout(rng)
	length := rng.Intn(1000) + 1
	text := randomText(rng, length)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s\n%s\n%s\n", l1, l2, text)
	return sb.String(), testCaseB{layout1: l1, layout2: l2, text: text}
}

func expectedB(tc testCaseB) string {
	mapping := make(map[byte]byte, 52)
	for i := 0; i < 26; i++ {
		c1 := tc.layout1[i]
		c2 := tc.layout2[i]
		mapping[c1] = c2
		mapping[c1-'a'+'A'] = c2 - 'a' + 'A'
	}
	res := make([]byte, len(tc.text))
	for i := 0; i < len(tc.text); i++ {
		ch := tc.text[i]
		if v, ok := mapping[ch]; ok {
			res[i] = v
		} else {
			res[i] = ch
		}
	}
	return string(res)
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, tc := generateCaseB(rng)
		expect := expectedB(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
