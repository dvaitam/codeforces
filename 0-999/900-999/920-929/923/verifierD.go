package main

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func applyRules(s string) []string {
	res := make([]string, 0)
	n := len(s)
	for i := 0; i < n; i++ {
		switch s[i] {
		case 'A':
			res = append(res, s[:i]+"BC"+s[i+1:])
		case 'B':
			res = append(res, s[:i]+"AC"+s[i+1:])
		case 'C':
			res = append(res, s[:i]+"AB"+s[i+1:])
		}
	}
	for i := 0; i+3 <= n; i++ {
		if s[i:i+3] == "AAA" {
			res = append(res, s[:i]+s[i+3:])
		}
	}
	return res
}

func canTransform(src, tgt string) bool {
	const maxLen = 10
	vis := map[string]bool{src: true}
	q := list.New()
	q.PushBack(src)
	for q.Len() > 0 {
		cur := q.Remove(q.Front()).(string)
		if cur == tgt {
			return true
		}
		if len(cur) > maxLen {
			continue
		}
		for _, nxt := range applyRules(cur) {
			if len(nxt) > maxLen {
				continue
			}
			if !vis[nxt] {
				vis[nxt] = true
				q.PushBack(nxt)
			}
		}
	}
	return false
}

type testD struct {
	S       string
	T       string
	queries [][4]int
}

func genString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	letters := []byte{'A', 'B', 'C'}
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(3)]
	}
	return string(b)
}

func genTest(rng *rand.Rand) testD {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	s := genString(rng, n)
	t := genString(rng, m)
	q := rng.Intn(3) + 1
	queries := make([][4]int, q)
	for i := 0; i < q; i++ {
		a := rng.Intn(n) + 1
		b := a + rng.Intn(n-a+1)
		c := rng.Intn(m) + 1
		d := c + rng.Intn(m-c+1)
		queries[i] = [4]int{a, b, c, d}
	}
	return testD{s, t, queries}
}

func solveD(tc testD) string {
	var sb strings.Builder
	for _, q := range tc.queries {
		sSub := tc.S[q[0]-1 : q[1]]
		tSub := tc.T[q[2]-1 : q[3]]
		if canTransform(sSub, tSub) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

func formatInput(tc testD) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s\n%s\n%d\n", tc.S, tc.T, len(tc.queries))
	for _, q := range tc.queries {
		fmt.Fprintf(&sb, "%d %d %d %d\n", q[0], q[1], q[2], q[3])
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genTest(rng)
		input := formatInput(tc)
		expected := solveD(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
