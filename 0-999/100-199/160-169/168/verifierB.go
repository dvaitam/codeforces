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

func solveB(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var out strings.Builder
	group := make([]string, 0)
	flush := func() {
		if len(group) == 0 {
			return
		}
		for _, line := range group {
			for i := 0; i < len(line); i++ {
				if line[i] != ' ' {
					out.WriteByte(line[i])
				}
			}
		}
		out.WriteByte('\n')
		group = group[:0]
	}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		isAmp := false
		for i := 0; i < len(line); i++ {
			if line[i] == ' ' {
				continue
			}
			if line[i] == '#' {
				isAmp = true
			}
			break
		}
		if isAmp {
			flush()
			out.WriteString(line)
			out.WriteByte('\n')
		} else {
			group = append(group, line)
		}
	}
	flush()
	return out.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func randomCase(rng *rand.Rand) string {
	lines := rng.Intn(10) + 1
	var sb strings.Builder
	hasChar := false
	for i := 0; i < lines; i++ {
		l := rng.Intn(20)
		for j := 0; j < l; j++ {
			switch rng.Intn(4) {
			case 0:
				sb.WriteByte(byte('a' + rng.Intn(26)))
				hasChar = true
			case 1:
				sb.WriteByte(' ')
			case 2:
				sb.WriteByte('#')
				hasChar = true
			default:
				sb.WriteByte(byte('A' + rng.Intn(26)))
				hasChar = true
			}
		}
		sb.WriteByte('\n')
	}
	if !hasChar {
		sb.WriteString("a\n")
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{
		"#abc\n  de f\n#x\n",
		"hello world\n",
		"#line1\nline2\nline3\n#end\n",
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		expect := solveB(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:\n%sgot:\n%s", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
