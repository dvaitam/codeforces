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

func solveD(s string) string {
	n := len(s)
	countEnd := make([]int64, n)
	countStart := make([]int64, n)
	for center := 0; center < n; center++ {
		l, r := center, center
		for l >= 0 && r < n && s[l] == s[r] {
			countStart[l]++
			countEnd[r]++
			l--
			r++
		}
	}
	for center := 0; center < n-1; center++ {
		l, r := center, center+1
		for l >= 0 && r < n && s[l] == s[r] {
			countStart[l]++
			countEnd[r]++
			l--
			r++
		}
	}
	suffix := make([]int64, n+1)
	for i := n - 1; i >= 0; i-- {
		suffix[i] = suffix[i+1] + countStart[i]
	}
	var ans int64
	for b := 0; b < n-1; b++ {
		if countEnd[b] > 0 {
			ans += countEnd[b] * suffix[b+1]
		}
	}
	return fmt.Sprintf("%d", ans)
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

func randString(rng *rand.Rand, n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := randString(rng, rng.Intn(8)+1)
		input := fmt.Sprintf("%s\n", s)
		expected := solveD(s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
