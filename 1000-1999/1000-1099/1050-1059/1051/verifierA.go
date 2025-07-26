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

func solveCase(s string) string {
	b := []byte(s)
	var U, D, L []int
	for i, c := range b {
		switch {
		case c >= 'A' && c <= 'Z':
			U = append(U, i)
		case c >= '0' && c <= '9':
			D = append(D, i)
		default:
			L = append(L, i)
		}
	}
	missingU := len(U) == 0
	missingD := len(D) == 0
	missingL := len(L) == 0
	missing := 0
	if missingU {
		missing++
	}
	if missingD {
		missing++
	}
	if missingL {
		missing++
	}
	switch missing {
	case 0:
		return s
	case 1:
		var rep byte
		if missingU {
			rep = 'A'
		} else if missingD {
			rep = '0'
		} else {
			rep = 'a'
		}
		if len(U) > 1 {
			b[U[0]] = rep
		} else if len(D) > 1 {
			b[D[0]] = rep
		} else {
			b[L[0]] = rep
		}
	case 2:
		miss := make([]byte, 0, 2)
		if missingU {
			miss = append(miss, 'A')
		}
		if missingD {
			miss = append(miss, '0')
		}
		if missingL {
			miss = append(miss, 'a')
		}
		if len(U) >= 2 {
			b[U[0]] = miss[0]
			b[U[1]] = miss[1]
		} else if len(D) >= 2 {
			b[D[0]] = miss[0]
			b[D[1]] = miss[1]
		} else {
			b[L[0]] = miss[0]
			b[L[1]] = miss[1]
		}
	}
	return string(b)
}

func randChar(rng *rand.Rand, cat int) byte {
	switch cat {
	case 0:
		return byte('a' + rng.Intn(26))
	case 1:
		return byte('A' + rng.Intn(26))
	default:
		return byte('0' + rng.Intn(10))
	}
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 3 // 3..10
	chars := make([]byte, n)
	for i := 0; i < n; i++ {
		cat := rng.Intn(3)
		chars[i] = randChar(rng, cat)
	}
	return string(chars)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := genCase(rng)
		in := fmt.Sprintf("1\n%s\n", s)
		expect := solveCase(s)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, s)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
