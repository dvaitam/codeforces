package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func repeatsNeeded(s string, p string) int {
	k := len(p)
	pos := make([][]int, 26)
	for i := 0; i < k; i++ {
		c := p[i] - 'a'
		pos[c] = append(pos[c], i)
	}
	idx := -1
	for i := 0; i < len(s); i++ {
		arr := pos[s[i]-'a']
		mod := -1
		if idx >= 0 {
			mod = idx % k
		}
		found := false
		for _, j := range arr {
			if j > mod {
				idx = (idx/k)*k + j
				found = true
				break
			}
		}
		if !found {
			idx = ((idx/k)+1)*k + arr[0]
		}
	}
	return idx/k + 1
}

func solveE(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return ""
	}
	var str string
	fmt.Fscan(in, &str)
	s := []byte(str)
	var out bytes.Buffer
	for ; m > 0; m-- {
		var tp int
		fmt.Fscan(in, &tp)
		if tp == 1 {
			var l, r int
			var c string
			fmt.Fscan(in, &l, &r, &c)
			for i := l - 1; i <= r-1; i++ {
				s[i] = c[0]
			}
		} else {
			var perm string
			fmt.Fscan(in, &perm)
			val := repeatsNeeded(string(s), perm)
			if out.Len() > 0 {
				out.WriteByte('\n')
			}
			fmt.Fprintf(&out, "%d", val)
		}
	}
	return out.String()
}

func genTestE(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	k := rng.Intn(4) + 1
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d %d %d\n", n, m, k)
	for i := 0; i < n; i++ {
		buf.WriteByte(byte('a' + rng.Intn(k)))
	}
	buf.WriteByte('\n')
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			c := byte('a' + rng.Intn(k))
			fmt.Fprintf(&buf, "1 %d %d %c\n", l, r, c)
		} else {
			perm := make([]byte, k)
			used := make([]bool, k)
			for j := 0; j < k; j++ {
				for {
					x := rng.Intn(k)
					if !used[x] {
						perm[j] = byte('a' + x)
						used[x] = true
						break
					}
				}
			}
			fmt.Fprintf(&buf, "2 %s\n", string(perm))
		}
	}
	return buf.String()
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestE(rng)
		expect := solveE(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
