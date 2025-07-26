package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

const letters = 26

func expectedD(words []string) string {
	out := make([]int, letters)
	inDeg := make([]int, letters)
	used := make([]bool, letters)
	for i := 0; i < letters; i++ {
		out[i] = -1
		inDeg[i] = -1
	}
	for _, s := range words {
		seen := make([]bool, letters)
		for i := 0; i < len(s); i++ {
			c := int(s[i] - 'a')
			used[c] = true
			if seen[c] {
				return "NO"
			}
			seen[c] = true
		}
		for i := 0; i+1 < len(s); i++ {
			u := int(s[i] - 'a')
			v := int(s[i+1] - 'a')
			if out[u] != -1 && out[u] != v {
				return "NO"
			}
			if inDeg[v] != -1 && inDeg[v] != u {
				return "NO"
			}
			out[u] = v
			inDeg[v] = u
		}
	}
	// detect cycles
	state := make([]int, letters)
	var bad bool
	var dfs func(int)
	dfs = func(v int) {
		state[v] = 1
		if out[v] != -1 {
			u := out[v]
			if state[u] == 1 {
				bad = true
				return
			}
			if state[u] == 0 {
				dfs(u)
				if bad {
					return
				}
			}
		}
		state[v] = 2
	}
	for i := 0; i < letters; i++ {
		if used[i] && state[i] == 0 {
			dfs(i)
			if bad {
				return "NO"
			}
		}
	}
	visited := make([]bool, letters)
	var chains []string
	for i := 0; i < letters; i++ {
		if used[i] && inDeg[i] == -1 {
			cur := i
			var b []byte
			for cur != -1 && !visited[cur] {
				visited[cur] = true
				b = append(b, byte(cur+'a'))
				cur = out[cur]
			}
			chains = append(chains, string(b))
		}
	}
	for i := 0; i < letters; i++ {
		if used[i] && !visited[i] {
			return "NO"
		}
	}
	sort.Strings(chains)
	res := strings.Join(chains, "")
	return res
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	m := rng.Intn(3) + 1
	words := make([]string, n)
	used := map[string]struct{}{}
	for i := 0; i < n; i++ {
		for {
			var sb strings.Builder
			for j := 0; j < m; j++ {
				sb.WriteByte(byte('a' + rng.Intn(4)))
			}
			s := sb.String()
			if _, ok := used[s]; !ok {
				used[s] = struct{}{}
				words[i] = s
				break
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, w := range words {
		sb.WriteString(w)
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String(), expected: expectedD(words)}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if out != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d: expected %s got %s\ninput:\n%s", i+1, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
