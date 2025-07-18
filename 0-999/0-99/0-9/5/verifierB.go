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

func generateCase(rng *rand.Rand) (string, string) {
	lines := rng.Intn(5) + 1
	var inputLines []string
	maxLen := 0
	for i := 0; i < lines; i++ {
		length := rng.Intn(20) + 1
		// generate string of letters/digits/spaces not starting or ending with space
		b := make([]byte, length)
		for {
			for j := 0; j < length; j++ {
				t := rng.Intn(36 + 1) // letters+digits+space
				if t < 26 {
					b[j] = byte('a' + t)
				} else if t < 36 {
					b[j] = byte('0' + t - 26)
				} else {
					b[j] = ' '
				}
			}
			if b[0] != ' ' && b[length-1] != ' ' {
				break
			}
		}
		s := string(b)
		inputLines = append(inputLines, s)
		if len(s) > maxLen {
			maxLen = len(s)
		}
	}
	// compute expected output
	border := strings.Repeat("*", maxLen+2)
	var out []string
	out = append(out, border)
	leftFlag := true
	for _, line := range inputLines {
		spaces := maxLen - len(line)
		l := spaces / 2
		r := spaces - l
		if spaces%2 == 1 {
			if leftFlag {
				l++
			} else {
				r++
			}
			leftFlag = !leftFlag
		}
		out = append(out, "*"+strings.Repeat(" ", l)+line+strings.Repeat(" ", r)+"*")
	}
	out = append(out, border)
	input := strings.Join(inputLines, "\n") + "\n"
	expected := strings.Join(out, "\n")
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
