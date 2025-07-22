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

func expected(s, t string) string {
	n, m := len(s), len(t)
	var cnt0 [26]int
	for i := 0; i < n; i++ {
		cnt0[s[i]-'a']++
	}
	maxPre := m
	if n < m {
		maxPre = n
	}
	cnt := make([][26]int, maxPre+1)
	ok := make([]bool, maxPre+1)
	cnt[0] = cnt0
	ok[0] = true
	for j := 0; j < maxPre; j++ {
		if !ok[j] {
			break
		}
		c := t[j] - 'a'
		if cnt[j][c] > 0 {
			cnt[j+1] = cnt[j]
			cnt[j+1][c]--
			ok[j+1] = true
		} else {
			break
		}
	}
	if n > m && ok[m] {
		res := make([]byte, n)
		for i := 0; i < m; i++ {
			res[i] = t[i]
		}
		for k := 0; k < 26; k++ {
			if cnt[m][k] > 0 {
				res[m] = byte('a' + k)
				cnt[m][k]--
				break
			}
		}
		idx := m + 1
		for k := 0; k < 26; k++ {
			for cnt[m][k] > 0 {
				res[idx] = byte('a' + k)
				cnt[m][k]--
				idx++
			}
		}
		return string(res)
	}
	for j := maxPre; j >= 1; j-- {
		if !ok[j-1] {
			continue
		}
		base := cnt[j-1]
		x := t[j-1] - 'a'
		kSel := -1
		for k := int(x) + 1; k < 26; k++ {
			if base[k] > 0 {
				kSel = k
				break
			}
		}
		if kSel < 0 {
			continue
		}
		res := make([]byte, n)
		for i := 0; i < j-1; i++ {
			res[i] = t[i]
		}
		res[j-1] = byte('a' + kSel)
		base[kSel]--
		idx := j
		for kk := 0; kk < 26; kk++ {
			for base[kk] > 0 {
				res[idx] = byte('a' + kk)
				base[kk]--
				idx++
			}
		}
		return string(res)
	}
	return "-1"
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

func runCase(bin, s, t string) error {
	input := fmt.Sprintf("%s\n%s\n", s, t)
	exp := expected(s, t)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != exp {
		return fmt.Errorf("expected %s got %s", exp, strings.TrimSpace(out))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := [][2]string{{"a", ""}, {"abc", "ab"}, {"ab", "ba"}, {"a", "a"}}

	letters := "abc"
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(n) + 1
		sb := make([]byte, n)
		for j := range sb {
			sb[j] = letters[rng.Intn(len(letters))]
		}
		tb := make([]byte, m)
		copy(tb, sb[:m])
		// random shuffle tb to produce some random t (not necessarily minimal)
		for j := 0; j < m; j++ {
			k := rng.Intn(j + 1)
			tb[j], tb[k] = tb[k], tb[j]
		}
		cases = append(cases, [2]string{string(sb), string(tb)})
	}

	for idx, c := range cases {
		if err := runCase(bin, c[0], c[1]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n%s\n", idx+1, err, c[0], c[1])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
