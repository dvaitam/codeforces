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

type pair struct{ a, b string }

func solveB(n, m int, pairs []pair, lecture []string) string {
	dict := make(map[string]string, m)
	for _, pr := range pairs {
		dict[pr.a] = pr.b
	}
	out := make([]string, n)
	for i, w := range lecture {
		alt := dict[w]
		if len(alt) < len(w) {
			out[i] = alt
		} else {
			out[i] = w
		}
	}
	return strings.Join(out, " ") + "\n"
}

func parseCase(input string) (int, int, []pair, []string) {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m int
	fmt.Fscan(in, &n, &m)
	pairs := make([]pair, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &pairs[i].a, &pairs[i].b)
	}
	lecture := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &lecture[i])
	}
	return n, m, pairs, lecture
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func randWord(rng *rand.Rand) string {
	l := rng.Intn(10) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(50) + 1
	m := rng.Intn(50) + 1
	if m < n {
		m = n
	}
	pairs := make([]pair, m)
	used := make(map[string]bool)
	for i := 0; i < m; i++ {
		var a string
		for {
			a = randWord(rng)
			if !used[a] {
				used[a] = true
				break
			}
		}
		var b string
		for {
			b = randWord(rng)
			if !used[b] {
				used[b] = true
				break
			}
		}
		pairs[i] = pair{a: a, b: b}
	}
	lecture := make([]string, n)
	for i := 0; i < n; i++ {
		lecture[i] = pairs[rng.Intn(m)].a
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		fmt.Fprintf(&sb, "%s %s\n", pairs[i].a, pairs[i].b)
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(lecture[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []string
	tests = append(tests, "3 3\ncodeforces code\nhello hi\nlonger word\ncodeforces hello longer\n")
	for i := 0; i < 99; i++ {
		tests = append(tests, generateCase(rng))
	}
	for i, tc := range tests {
		n, m, pairs, lecture := parseCase(tc)
		expect := solveB(n, m, pairs, lecture)
		out, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(expect), strings.TrimSpace(out), tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
