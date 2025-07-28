package main

import (
	"bufio"
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func steps(s []byte, start int) int {
	n := len(s)
	b := make([]byte, n)
	copy(b, s)
	pos := start
	t := 0
	for pos >= 0 && pos < n {
		if b[pos] == '>' {
			b[pos] = '<'
			pos++
		} else {
			b[pos] = '>'
			pos--
		}
		t++
	}
	return t
}

func referenceB(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var sb strings.Builder
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var str string
		fmt.Fscan(in, &str)
		s := []byte(str)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(steps(s, i)))
		}
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func genCaseB(rng *rand.Rand) string {
	t := 1
	n := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%d\n", t, n))
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('<')
		} else {
			sb.WriteByte('>')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCaseB(rng)
		expect := referenceB(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
