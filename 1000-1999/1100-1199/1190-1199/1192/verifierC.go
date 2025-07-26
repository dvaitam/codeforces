package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const numTestsC = 100
const mod = 998244353

var edges = [12][2]int{
	{0, 1}, {4, 5}, {2, 3}, {6, 7},
	{0, 2}, {1, 3}, {4, 6}, {5, 7},
	{0, 4}, {1, 5}, {2, 6}, {3, 7},
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verifC_bin")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func randWord(rng *rand.Rand, l int) string {
	letters := []byte("abc")
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func countCubes(words []string) int {
	n := len(words)
	L := len(words[0])
	first := make([]byte, n)
	last := make([]byte, n)
	for i, w := range words {
		first[i] = w[0]
		last[i] = w[L-1]
	}
	var vert [8]byte
	var used [8]bool
	count := 0
	var dfs func(int)
	dfs = func(e int) {
		if e == len(edges) {
			count = (count + 1) % mod
			return
		}
		v1 := edges[e][0]
		v2 := edges[e][1]
		for i := 0; i < n; i++ {
			s := first[i]
			t := last[i]
			if (!used[v1] || vert[v1] == s) && (!used[v2] || vert[v2] == t) {
				o1, o2 := vert[v1], vert[v2]
				u1, u2 := used[v1], used[v2]
				vert[v1], used[v1] = s, true
				vert[v2], used[v2] = t, true
				dfs(e + 1)
				vert[v1], used[v1] = o1, u1
				vert[v2], used[v2] = o2, u2
			}
			s2 := last[i]
			t2 := first[i]
			if (!used[v1] || vert[v1] == s2) && (!used[v2] || vert[v2] == t2) {
				o1, o2 := vert[v1], vert[v2]
				u1, u2 := used[v1], used[v2]
				vert[v1], used[v1] = s2, true
				vert[v2], used[v2] = t2, true
				dfs(e + 1)
				vert[v1], used[v1] = o1, u1
				vert[v2], used[v2] = o2, u2
			}
		}
	}
	dfs(0)
	return count % mod
}

func generateCaseC(rng *rand.Rand) (string, string) {
	n := rng.Intn(2) + 1
	l := rng.Intn(2) + 3
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = randWord(rng, l)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, w := range words {
		sb.WriteString(w)
		sb.WriteByte('\n')
	}
	input := sb.String()
	ans := countCubes(words)
	expected := strconv.Itoa(ans)
	return input, expected
}

func runCase(bin, input, expected string) error {
	got, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < numTestsC; i++ {
		in, exp := generateCaseC(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, in)
			return
		}
	}
	fmt.Println("All tests passed")
}
