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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type block struct {
	tryLine   int
	catchLine int
	exType    string
	msg       string
}

func solveCFromInput(input string) (string, error) {
	r := bufio.NewScanner(strings.NewReader(input))
	r.Scan()
	var n int
	fmt.Sscan(r.Text(), &n)
	lines := make([]string, 0, n)
	for r.Scan() {
		lines = append(lines, strings.TrimSpace(r.Text()))
	}
	if len(lines) != n {
		return "", fmt.Errorf("wrong line count")
	}
	tryStack := []int{}
	blocks := []block{}
	throwLine := -1
	throwType := ""
	for i, line := range lines {
		s := strings.TrimSpace(line)
		if strings.HasPrefix(s, "try") {
			tryStack = append(tryStack, i)
		} else if strings.HasPrefix(s, "catch") {
			p1 := strings.Index(s, "(")
			p2 := strings.LastIndex(s, ")")
			inner := s[p1+1 : p2]
			ci := strings.Index(inner, ",")
			typ := strings.TrimSpace(inner[:ci])
			msg := strings.TrimSpace(inner[ci+1:])
			if len(msg) >= 2 && msg[0] == '"' && msg[len(msg)-1] == '"' {
				msg = msg[1 : len(msg)-1]
			}
			t := tryStack[len(tryStack)-1]
			tryStack = tryStack[:len(tryStack)-1]
			blocks = append(blocks, block{tryLine: t, catchLine: i, exType: typ, msg: msg})
		} else if strings.HasPrefix(s, "throw") {
			p1 := strings.Index(s, "(")
			p2 := strings.LastIndex(s, ")")
			inner := strings.TrimSpace(s[p1+1 : p2])
			throwType = inner
			throwLine = i
		}
	}
	bestCatch := n + 1
	bestMsg := ""
	for _, b := range blocks {
		if b.tryLine < throwLine && throwLine < b.catchLine && b.exType == throwType {
			if b.catchLine < bestCatch {
				bestCatch = b.catchLine
				bestMsg = b.msg
			}
		}
	}
	if bestMsg == "" {
		return "Unhandled Exception", nil
	}
	return bestMsg, nil
}

func randomIdent(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	letters := "abcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rng.Intn(len(letters))])
	}
	return sb.String()
}

func randomMsg(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	letters := "abcde"
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rng.Intn(len(letters))])
	}
	return sb.String()
}

func generateCaseC(rng *rand.Rand) string {
	depth := rng.Intn(5) + 1
	var lines []string
	for i := 0; i < depth; i++ {
		lines = append(lines, "try")
	}
	thrType := randomIdent(rng)
	lines = append(lines, fmt.Sprintf("throw(%s)", thrType))
	for i := 0; i < depth; i++ {
		typ := randomIdent(rng)
		msg := randomMsg(rng)
		lines = append(lines, fmt.Sprintf("catch(%s, \"%s\")", typ, msg))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(lines))
	for _, l := range lines {
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseC(rng)
		expect, err := solveCFromInput(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
