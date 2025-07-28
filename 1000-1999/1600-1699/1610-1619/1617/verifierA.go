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

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return out.String() + errb.String(), err
	}
	return out.String(), nil
}

func solve(s, perm string) string {
	cnt := make([]int, 26)
	for i := 0; i < len(s); i++ {
		cnt[s[i]-'a']++
	}
	var sb strings.Builder
	if perm == "abc" && cnt[0] > 0 && cnt[1] > 0 && cnt[2] > 0 {
		for i := 0; i < cnt[0]; i++ {
			sb.WriteByte('a')
		}
		for i := 0; i < cnt[2]; i++ {
			sb.WriteByte('c')
		}
		for i := 0; i < cnt[1]; i++ {
			sb.WriteByte('b')
		}
		for ch := byte('d'); ch <= 'z'; ch++ {
			for i := 0; i < cnt[ch-'a']; i++ {
				sb.WriteByte(ch)
			}
		}
	} else {
		for ch := byte('a'); ch <= 'z'; ch++ {
			for i := 0; i < cnt[ch-'a']; i++ {
				sb.WriteByte(ch)
			}
		}
	}
	return sb.String()
}

func genTest(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + rng.Intn(26)))
	}
	letters := []byte{'a', 'b', 'c'}
	rng.Shuffle(3, func(i, j int) { letters[i], letters[j] = letters[j], letters[i] })
	return sb.String(), string(letters)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const tests = 100
	var input bytes.Buffer
	input.WriteString(fmt.Sprintf("%d\n", tests))
	expected := make([]string, tests)
	for i := 0; i < tests; i++ {
		s, perm := genTest(rng)
		expected[i] = solve(s, perm)
		input.WriteString(s)
		input.WriteByte('\n')
		input.WriteString(perm)
		input.WriteByte('\n')
	}
	out, err := run(bin, input.String())
	if err != nil {
		fmt.Printf("runtime error: %v\n%s", err, out)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != tests {
		fmt.Printf("expected %d lines of output got %d\n", tests, len(lines))
		os.Exit(1)
	}
	for i := 0; i < tests; i++ {
		if strings.TrimSpace(lines[i]) != expected[i] {
			fmt.Printf("test %d failed\nexpected:%s got:%s\n", i+1, expected[i], strings.TrimSpace(lines[i]))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
