package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var alphabet = []byte("abcdefghijklmnopqrstuvwxyz")

type testCaseA struct {
	a, b int
	l, r int64
}

func build(a, b int, ch byte) ([]byte, []byte) {
	seq := make([]byte, a)
	for i := 0; i < a; i++ {
		seq[i] = alphabet[i]
	}
	seen := map[string]int{}
	cycle := 0
	var start, period int
	for {
		for i := 0; i < b; i++ {
			seq = append(seq, ch)
		}
		suffix := seq[len(seq)-a:]
		used := make(map[byte]bool)
		for _, c := range suffix {
			used[c] = true
		}
		t := make([]byte, 0, a)
		for i := 0; i < 26 && len(t) < a; i++ {
			if !used[alphabet[i]] {
				t = append(t, alphabet[i])
			}
		}
		seq = append(seq, t...)
		suffix = seq[len(seq)-a:]
		cycle++
		key := string(suffix)
		if val, ok := seen[key]; ok {
			start = val
			period = cycle - val
			break
		}
		seen[key] = cycle
	}
	prefixLen := a + (start-1)*(a+b)
	periodLen := period * (a + b)
	prefix := append([]byte{}, seq[:prefixLen]...)
	periodStr := append([]byte{}, seq[prefixLen:prefixLen+periodLen]...)
	return prefix, periodStr
}

func uniqueCount(prefix, period []byte, l, r int64) int {
	ans := map[byte]struct{}{}
	prefixLen := int64(len(prefix))
	periodLen := int64(len(period))
	if l <= prefixLen {
		end := r
		if end > prefixLen {
			end = prefixLen
		}
		for i := l - 1; i < end; i++ {
			ans[prefix[i]] = struct{}{}
		}
		if r <= prefixLen {
			return len(ans)
		}
		l = prefixLen + 1
	}
	if periodLen == 0 {
		return len(ans)
	}
	if r-l+1 >= periodLen {
		for _, c := range period {
			ans[c] = struct{}{}
		}
	} else {
		start := (l - prefixLen - 1) % periodLen
		for i := int64(0); i < r-l+1; i++ {
			ans[period[(start+i)%periodLen]] = struct{}{}
		}
	}
	return len(ans)
}

func solveA(a, b int, l, r int64) int {
	best := 27
	for i := 0; i < 26; i++ {
		prefix, period := build(a, b, alphabet[i])
		val := uniqueCount(prefix, period, l, r)
		if val < best {
			best = val
		}
	}
	return best
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCaseA {
	rand.Seed(42)
	tests := make([]testCaseA, 100)
	for i := range tests {
		a := rand.Intn(5) + 1
		b := rand.Intn(5) + 1
		l := rand.Int63n(80) + 1
		r := l + rand.Int63n(20)
		tests[i] = testCaseA{a, b, l, r}
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d %d %d\n", tc.a, tc.b, tc.l, tc.r)
		expected := solveA(tc.a, tc.b, tc.l, tc.r)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", i+1, err)
			return
		}
		var got int
		fmt.Sscan(output, &got)
		if got != expected {
			fmt.Printf("test %d failed:\ninput: %sexpected %d got %s\n", i+1, input, expected, output)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
