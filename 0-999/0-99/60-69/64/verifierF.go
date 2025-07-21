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

func isDomain(s string) bool {
	if len(s) == 0 || s[0] == '.' || s[len(s)-1] == '.' {
		return false
	}
	prevDot := false
	lastDot := -1
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '.' {
			if prevDot {
				return false
			}
			prevDot = true
			lastDot = i
		} else if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			prevDot = false
		} else {
			return false
		}
	}
	var lastLen int
	if lastDot == -1 {
		lastLen = len(s)
	} else {
		lastLen = len(s) - lastDot - 1
	}
	return lastLen == 2 || lastLen == 3
}

func randToken(rng *rand.Rand, min, max int) string {
	n := rng.Intn(max-min+1) + min
	b := make([]byte, n)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = byte('a' + rng.Intn(26))
		} else {
			b[i] = byte('0' + rng.Intn(10))
		}
	}
	return string(b)
}

func generateValid(rng *rand.Rand) string {
	parts := rng.Intn(3) + 1
	var sb strings.Builder
	for i := 0; i < parts; i++ {
		if i > 0 {
			sb.WriteByte('.')
		}
		sb.WriteString(randToken(rng, 1, 5))
	}
	tldLen := 2 + rng.Intn(2)
	sb.WriteByte('.')
	sb.WriteString(randToken(rng, tldLen, tldLen))
	return sb.String()
}

func mutateInvalid(rng *rand.Rand, s string) string {
	choice := rng.Intn(4)
	switch choice {
	case 0:
		return "." + s
	case 1:
		return s + "."
	case 2:
		return strings.ReplaceAll(s, ".", "..")
	default:
		return s + "#"
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		valid := rng.Intn(2) == 0
		dom := generateValid(rng)
		if !valid {
			dom = mutateInvalid(rng, dom)
		}
		input := dom + "\n"
		expected := "NO"
		if isDomain(dom) {
			expected = "YES"
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
