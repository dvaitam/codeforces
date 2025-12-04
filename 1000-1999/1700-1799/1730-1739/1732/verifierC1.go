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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func solve(reader *bufio.Reader) string {
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return ""
	}
	var sb strings.Builder
	for ; T > 0; T-- {
		var n, Q int
		fmt.Fscan(reader, &n, &Q)
		v := make([]int, n+2)
		psum := make([]int64, n+2)
		xsum := make([]int64, n+2)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &v[i])
			psum[i] = psum[i-1] + int64(v[i])
			xsum[i] = xsum[i-1] ^ int64(v[i])
		}
		nxt := make([]int, n+2)
		id := n + 1
		for i := n; i >= 1; i-- {
			nxt[i] = id
			if v[i] != 0 {
				id = i
			}
		}
		var a, b int
		fmt.Fscan(reader, &a, &b)
		getval := func(i, j int) int64 {
			return (psum[j] - psum[i-1]) - (xsum[j] ^ xsum[i-1])
		}
		mxval := getval(a, b)
		s := a
		if s <= n && v[s] == 0 {
			s = nxt[s]
		}
		ansA, ansB := a, b
		for i := 0; i < 31; i++ {
			if s > b {
				s = b
				break
			}
			if getval(s, b) != mxval {
				break
			}
			l, r := s-1, b
			for l+1 < r {
				m := (l + r) >> 1
				if getval(s, m) == mxval {
					r = m
				} else {
					l = m
				}
			}
			if ansB-ansA > r-s {
				ansA = s
				ansB = r
			}
			if s <= n {
				s = nxt[s]
			}
		}
		fmt.Fprintf(&sb, "%d %d\n", ansA, ansB)
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) string {
	T := rng.Intn(2) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", T)
	for ; T > 0; T-- {
		n := rng.Intn(50) + 1
		fmt.Fprintf(&sb, "%d 1\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(100))
		}
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "1 %d\n", n)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solve(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
