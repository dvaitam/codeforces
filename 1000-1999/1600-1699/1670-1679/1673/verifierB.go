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

type testCase struct {
	s string
}

func generateTests() []testCase {
	r := rand.New(rand.NewSource(2))
	tests := []testCase{{"a"}, {"ab"}, {"abb"}, {"abc"}, {"aabb"}}
	for len(tests) < 100 {
		n := r.Intn(8) + 1
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = byte('a' + r.Intn(3))
		}
		tests = append(tests, testCase{s: string(b)})
	}
	return tests
}

func expected(s string) string {
	letters := make(map[byte]bool)
	for i := 0; i < len(s); i++ {
		letters[s[i]] = true
	}
	n := len(s)
	pref := make([][26]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1]
		pref[i][s[i-1]-'a']++
	}
	pos := make(map[byte][]int)
	for i := 0; i < n; i++ {
		c := s[i]
		pos[c] = append(pos[c], i+1)
	}
	balanced := true
	for _, arr := range pos {
		for j := 1; j < len(arr) && balanced; j++ {
			l := arr[j-1]
			r := arr[j]
			for ch := range letters {
				idx := ch - 'a'
				if pref[r][idx]-pref[l-1][idx] == 0 {
					balanced = false
					break
				}
			}
		}
		if !balanced {
			break
		}
	}
	if balanced {
		return "YES"
	}
	return "NO"
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("run failed: %v\n%s", err, errb.String())
	}
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("1\n%s\n", t.s)
		want := expected(t.s)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != want {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, want, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
