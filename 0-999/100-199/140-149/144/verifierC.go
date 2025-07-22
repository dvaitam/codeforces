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

func expected(s, p string) string {
	n := len(s)
	m := len(p)
	if m > n {
		return "0"
	}
	var countP [26]int
	for i := 0; i < m; i++ {
		countP[p[i]-'a']++
	}
	var countS [26]int
	question := 0
	for i := 0; i < m; i++ {
		ch := s[i]
		if ch == '?' {
			question++
		} else {
			countS[ch-'a']++
		}
	}
	check := func() bool {
		total := 0
		for c := 0; c < 26; c++ {
			if countS[c] > countP[c] {
				return false
			}
			total += countP[c] - countS[c]
		}
		return total == question
	}
	ans := 0
	if check() {
		ans++
	}
	for i := m; i < n; i++ {
		old := s[i-m]
		if old == '?' {
			question--
		} else {
			countS[old-'a']--
		}
		ch := s[i]
		if ch == '?' {
			question++
		} else {
			countS[ch-'a']++
		}
		if check() {
			ans++
		}
	}
	return fmt.Sprintf("%d", ans)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	m := rng.Intn(n) + 1
	alphabet := []byte("abcde")
	sbS := strings.Builder{}
	for i := 0; i < n; i++ {
		if rng.Intn(5) == 0 {
			sbS.WriteByte('?')
		} else {
			sbS.WriteByte(alphabet[rng.Intn(len(alphabet))])
		}
	}
	sbP := strings.Builder{}
	for i := 0; i < m; i++ {
		sbP.WriteByte(alphabet[rng.Intn(len(alphabet))])
	}
	s := sbS.String()
	p := sbP.String()
	input := s + "\n" + p + "\n"
	exp := expected(s, p)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
