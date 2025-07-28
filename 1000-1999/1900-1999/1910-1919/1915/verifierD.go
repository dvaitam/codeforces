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

func runBinary(path, input string) (string, error) {
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

func splitSyllables(s string) string {
	isVowel := func(c byte) bool { return c == 'a' || c == 'e' }
	n := len(s)
	vPos := []int{}
	for i := 0; i < n; i++ {
		if isVowel(s[i]) {
			vPos = append(vPos, i)
		}
	}
	boundaries := make(map[int]bool)
	for i := 0; i+1 < len(vPos); i++ {
		diff := vPos[i+1] - vPos[i]
		if diff == 2 {
			boundaries[vPos[i]] = true
		} else if diff == 3 {
			boundaries[vPos[i]+1] = true
		}
	}
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		buf.WriteByte(s[i])
		if boundaries[i] {
			buf.WriteByte('.')
		}
	}
	return buf.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	vowels := []byte{'a', 'e'}
	cons := []byte{'b', 'c', 'd'}
	syl := rng.Intn(5) + 1
	var word strings.Builder
	for i := 0; i < syl; i++ {
		word.WriteByte(cons[rng.Intn(len(cons))])
		word.WriteByte(vowels[rng.Intn(len(vowels))])
		if rng.Intn(2) == 0 {
			word.WriteByte(cons[rng.Intn(len(cons))])
		}
	}
	w := word.String()
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n%s\n", len(w), w)
	input := sb.String()
	expected := splitSyllables(w)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := generateCase(rng)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
