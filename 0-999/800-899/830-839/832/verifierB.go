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

type Test struct {
	good    string
	pattern string
	queries []string
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(t.good + "\n")
	sb.WriteString(t.pattern + "\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(t.queries)))
	for _, q := range t.queries {
		sb.WriteString(q + "\n")
	}
	return sb.String()
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "832B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func randString(alphabet string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func genTests() []Test {
	rand.Seed(time.Now().UnixNano())
	tests := make([]Test, 0, 100)
	letters := "abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < 100; i++ {
		// good letters
		gcount := rand.Intn(26) + 1
		goodSet := make(map[byte]bool)
		for len(goodSet) < gcount {
			ch := letters[rand.Intn(26)]
			goodSet[ch] = true
		}
		var good strings.Builder
		for ch := range goodSet {
			good.WriteByte(ch)
		}

		// pattern
		plen := rand.Intn(10) + 1
		var pattern strings.Builder
		hasStar := rand.Intn(4) == 0
		starPos := -1
		for j := 0; j < plen; j++ {
			if hasStar && starPos == -1 && rand.Intn(plen) == 0 {
				pattern.WriteByte('*')
				starPos = j
				continue
			}
			r := rand.Intn(26 + 1)
			if r == 26 {
				pattern.WriteByte('?')
			} else {
				pattern.WriteByte(letters[r])
			}
		}
		if hasStar && starPos == -1 {
			pattern.WriteByte('*')
		}

		// queries
		qnum := rand.Intn(5) + 1
		qs := make([]string, qnum)
		for j := 0; j < qnum; j++ {
			l := rand.Intn(12) + 1
			qs[j] = randString(letters, l)
		}
		tests = append(tests, Test{good: good.String(), pattern: pattern.String(), queries: qs})
	}
	// simple case
	tests = append(tests, Test{good: "abc", pattern: "a?c", queries: []string{"abc"}})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:%sExpected:%sGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
