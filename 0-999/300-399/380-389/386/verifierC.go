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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(s string) string {
	n := len(s)
	seen := make([]bool, 26)
	d := 0
	for i := 0; i < n; i++ {
		idx := s[i] - 'a'
		if !seen[idx] {
			seen[idx] = true
			d++
		}
	}
	f := make([]int64, d+1)
	for k := 1; k <= d; k++ {
		var cnt [26]int
		distinct := 0
		l := 0
		var res int64
		for r := 0; r < n; r++ {
			idx := s[r] - 'a'
			if cnt[idx] == 0 {
				distinct++
			}
			cnt[idx]++
			for distinct > k {
				idxl := s[l] - 'a'
				cnt[idxl]--
				if cnt[idxl] == 0 {
					distinct--
				}
				l++
			}
			res += int64(r - l + 1)
		}
		f[k] = res
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", d))
	for k := 1; k <= d; k++ {
		tk := f[k]
		if k > 1 {
			tk -= f[k-1]
		}
		sb.WriteString(fmt.Sprintf("%d\n", tk))
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(30) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(rng.Intn(26)) + 'a'
	}
	s := string(b)
	return fmt.Sprintf("%s\n", s), solveCase(s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct{ in, out string }{
		{"a\n", solveCase("a")},
		{"abca\n", solveCase("abca")},
	}
	for i := 0; i < 100; i++ {
		in, out := generateCase(rng)
		cases = append(cases, struct{ in, out string }{in, out})
	}

	for i, tc := range cases {
		got, err := runCandidate(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(tc.out) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.out, got, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
