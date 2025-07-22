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

func isVowel(b byte) bool {
	switch b {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

func getSuffix(s string, k int) string {
	cnt := 0
	for i := len(s) - 1; i >= 0; i-- {
		if isVowel(s[i]) {
			cnt++
			if cnt == k {
				return s[i:]
			}
		}
	}
	return ""
}

func expectedAnswer(lines []string, k int) string {
	n := len(lines) / 4
	suffixes := make([]string, len(lines))
	for i, line := range lines {
		suffixes[i] = getSuffix(line, k)
		if suffixes[i] == "" {
			return "NO"
		}
	}
	possibleAABB, possibleABAB, possibleABBA := true, true, true
	allAAAA := true
	for i := 0; i < n; i++ {
		a := suffixes[4*i]
		b := suffixes[4*i+1]
		c := suffixes[4*i+2]
		d := suffixes[4*i+3]
		if a == b && b == c && c == d {
			continue
		}
		allAAAA = false
		curAABB := (a == b && c == d)
		curABAB := (a == c && b == d)
		curABBA := (a == d && b == c)
		if !curAABB && !curABAB && !curABBA {
			return "NO"
		}
		if !curAABB {
			possibleAABB = false
		}
		if !curABAB {
			possibleABAB = false
		}
		if !curABBA {
			possibleABBA = false
		}
	}
	if allAAAA {
		return "aaaa"
	}
	count := 0
	scheme := ""
	if possibleAABB {
		count++
		scheme = "aabb"
	}
	if possibleABAB {
		count++
		scheme = "abab"
	}
	if possibleABBA {
		count++
		scheme = "abba"
	}
	if count == 1 {
		return scheme
	}
	return "NO"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	k := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	lines := make([]string, 4*n)
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	vowels := "aeiou"
	for i := 0; i < 4*n; i++ {
		l := rng.Intn(10) + 1
		var w strings.Builder
		for j := 0; j < l; j++ {
			if rng.Float64() < 0.4 {
				w.WriteByte(vowels[rng.Intn(len(vowels))])
			} else {
				w.WriteByte(letters[rng.Intn(len(letters))])
			}
		}
		lines[i] = w.String()
		fmt.Fprintln(&sb, lines[i])
	}
	exp := expectedAnswer(lines, k)
	return sb.String(), exp
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", t, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
