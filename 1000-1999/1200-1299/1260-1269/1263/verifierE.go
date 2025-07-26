package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveE(n int, s string) string {
	arr := make([]int, n+5)
	pos := 1
	maxPos := 1
	results := make([]int, n)
	for i := 0; i < n; i++ {
		c := s[i]
		switch c {
		case 'L':
			if pos > 1 {
				pos--
			}
		case 'R':
			pos++
			if pos > maxPos {
				maxPos = pos
			}
		case '(':
			arr[pos] = 1
			if pos > maxPos {
				maxPos = pos
			}
		case ')':
			arr[pos] = -1
			if pos > maxPos {
				maxPos = pos
			}
		default:
			arr[pos] = 0
			if pos > maxPos {
				maxPos = pos
			}
		}
		sum := 0
		minV := 0
		maxV := 0
		for j := 1; j <= maxPos; j++ {
			sum += arr[j]
			if sum < minV {
				minV = sum
			}
			if sum > maxV {
				maxV = sum
			}
		}
		if sum != 0 || minV < 0 {
			results[i] = -1
		} else {
			results[i] = maxV
		}
	}
	var sb strings.Builder
	for i, v := range results {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	cmds := make([]byte, n)
	alphabet := []byte{'L', 'R', '(', ')', 'a'}
	for i := 0; i < n; i++ {
		cmds[i] = alphabet[rng.Intn(len(alphabet))]
	}
	s := string(cmds)
	input := fmt.Sprintf("%d\n%s\n", n, s)
	return input, solveE(n, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\n", i+1, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
