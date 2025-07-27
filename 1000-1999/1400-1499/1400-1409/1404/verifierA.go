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
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return ""
	}
	var sb strings.Builder
	for ; t > 0; t-- {
		var n, k int
		var s string
		fmt.Fscan(reader, &n, &k)
		fmt.Fscan(reader, &s)
		tarr := make([]byte, k)
		for i := 0; i < k; i++ {
			tarr[i] = '?'
		}
		ok := true
		for i := 0; i < n; i++ {
			c := s[i]
			if c == '?' {
				continue
			}
			pos := i % k
			if tarr[pos] == '?' {
				tarr[pos] = c
			} else if tarr[pos] != c {
				ok = false
				break
			}
		}
		if !ok {
			sb.WriteString("NO\n")
			continue
		}
		cnt0, cnt1 := 0, 0
		for i := 0; i < k; i++ {
			if tarr[i] == '0' {
				cnt0++
			} else if tarr[i] == '1' {
				cnt1++
			}
		}
		if cnt0 > k/2 || cnt1 > k/2 {
			sb.WriteString("NO\n")
		} else {
			sb.WriteString("YES\n")
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for ; t > 0; t-- {
		n := rng.Intn(30) + 2
		if n%2 == 1 {
			n++ // make even sometimes? not necessary, k even only
		}
		k := rng.Intn(n/2)*2 + 2 // ensure even and <=n
		if k > n {
			k = n
			if k%2 == 1 {
				k--
			}
		}
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		// generate random string of 0/1/?
		data := make([]byte, n)
		for i := 0; i < n; i++ {
			r := rng.Intn(3)
			if r == 0 {
				data[i] = '0'
			} else if r == 1 {
				data[i] = '1'
			} else {
				data[i] = '?'
			}
		}
		fmt.Fprintf(&sb, "%s\n", string(data))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
