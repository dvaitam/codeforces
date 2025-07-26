package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func decode(b string) string {
	set := make(map[rune]struct{})
	for _, ch := range b {
		set[ch] = struct{}{}
	}
	letters := make([]rune, 0, len(set))
	for ch := range set {
		letters = append(letters, ch)
	}
	sort.Slice(letters, func(i, j int) bool { return letters[i] < letters[j] })
	m := make(map[rune]rune)
	n := len(letters)
	for i, ch := range letters {
		m[ch] = letters[n-1-i]
	}
	res := make([]rune, len(b))
	for i, ch := range b {
		res[i] = m[ch]
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append(os.Args[:1], os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		sb := make([]rune, n)
		for j := 0; j < n; j++ {
			sb[j] = letters[rng.Intn(len(letters))]
		}
		s := string(sb)
		b := decode(s)
		input := fmt.Sprintf("1\n%d\n%s\n", n, b)
		expected := s
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
