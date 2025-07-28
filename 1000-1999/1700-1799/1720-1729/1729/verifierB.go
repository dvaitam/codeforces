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

func encodeString(s string) string {
	var sb strings.Builder
	for _, ch := range s {
		idx := int(ch - 'a' + 1)
		if idx >= 10 {
			sb.WriteString(fmt.Sprintf("%d0", idx))
		} else {
			sb.WriteString(fmt.Sprintf("%d", idx))
		}
	}
	return sb.String()
}

func decode(t string) string {
	res := make([]byte, 0)
	for i := len(t) - 1; i >= 0; {
		if t[i] == '0' {
			num := int(t[i-2]-'0')*10 + int(t[i-1]-'0')
			res = append(res, byte('a'+num-1))
			i -= 3
		} else {
			num := int(t[i] - '0')
			res = append(res, byte('a'+num-1))
			i--
		}
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return string(res)
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 0, 120)
	// simple deterministic strings
	base := []string{"a", "b", "z", "abc", "leetcode"}
	cases = append(cases, base...)
	for len(cases) < 120 {
		l := rng.Intn(10) + 1
		var sb strings.Builder
		for i := 0; i < l; i++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		t := encodeString(sb.String())
		if len(t) <= 50 {
			cases = append(cases, sb.String())
		}
	}

	for i, s := range cases {
		t := encodeString(s)
		input := fmt.Sprintf("1\n%d\n%s\n", len(t), t)
		expected := decode(t)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
