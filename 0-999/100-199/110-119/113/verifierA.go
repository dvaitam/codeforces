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

// solveA implements the logic from 113A.go
func solveA(input string) string {
	words := strings.Fields(input)
	if len(words) == 0 {
		return "NO"
	}
	type suffix struct {
		suf    string
		typ    int
		gender int
	}
	suffixes := []suffix{
		{"lios", 0, 0},
		{"liala", 0, 1},
		{"etr", 1, 0},
		{"etra", 1, 1},
		{"initis", 2, 0},
		{"inites", 2, 1},
	}
	types := make([]int, len(words))
	var gender *int
	nounCount := 0
	for i, w := range words {
		matched := false
		for _, s := range suffixes {
			if strings.HasSuffix(w, s.suf) {
				matched = true
				types[i] = s.typ
				if gender == nil {
					g := s.gender
					gender = &g
				} else if *gender != s.gender {
					return "NO"
				}
				if s.typ == 1 {
					nounCount++
				}
				break
			}
		}
		if !matched {
			return "NO"
		}
	}
	if len(words) == 1 {
		return "YES"
	}
	if nounCount != 1 {
		return "NO"
	}
	state := 0
	for _, t := range types {
		switch state {
		case 0:
			if t == 0 {
				continue
			} else if t == 1 {
				state = 1
			} else {
				return "NO"
			}
		case 1:
			if t == 1 {
				continue
			} else if t == 2 {
				state = 2
			} else {
				return "NO"
			}
		case 2:
			if t == 2 {
				continue
			} else {
				return "NO"
			}
		}
	}
	return "YES"
}

func randWord(rng *rand.Rand, end string) string {
	n := rng.Intn(3) // 0..2 letters prefix
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + rng.Intn(26)))
	}
	sb.WriteString(end)
	return sb.String()
}

func randInvalidWord(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + rng.Intn(26)))
	}
	sb.WriteString("zz")
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	endings := []string{"lios", "liala", "etr", "etra", "initis", "inites"}
	n := rng.Intn(6) + 1
	words := make([]string, n)
	for i := 0; i < n; i++ {
		if rng.Intn(5) == 0 { // 20% invalid
			words[i] = randInvalidWord(rng)
		} else {
			words[i] = randWord(rng, endings[rng.Intn(len(endings))])
		}
	}
	line := strings.Join(words, " ")
	return line, solveA(line)
}

func runCase(exe string, input string, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	edgeCases := []string{
		"lios",
		"foo",
		"etr etr",
		"liala etr inites",
		"aaaainites",
	}
	for i, tc := range edgeCases {
		if err := runCase(exe, tc, solveA(tc)); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\ninput:\n%s\n", i+1, err, tc)
			os.Exit(1)
		}
	}
	for i := 0; i < 100; i++ {
		tc, exp := generateCase(rng)
		if err := runCase(exe, tc, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
