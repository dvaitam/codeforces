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

func solve(s string) string {
	n := len(s)
	type state struct {
		ok  bool
		str string
	}
	dp := make([][2]state, n+1)
	dp[0][0] = state{true, ""}
	update := func(i, used int, cand string) {
		st := &dp[i][used]
		if !st.ok || len(cand) < len(st.str) || (len(cand) == len(st.str) && cand < st.str) {
			st.ok = true
			st.str = cand
		}
	}
	for i := 0; i < n; i++ {
		for used := 0; used < 2; used++ {
			cur := dp[i][used]
			if !cur.ok {
				continue
			}
			c := s[i]
			update(i+1, used, cur.str+string(c))
			if used == 0 && i > 0 && i+2 < n {
				if s[i] == 'a' && s[i+1] == 't' {
					update(i+2, 1, cur.str+"@")
				}
			}
			if i > 0 && i+3 < n {
				if s[i] == 'd' && s[i+1] == 'o' && s[i+2] == 't' {
					update(i+3, used, cur.str+".")
				}
			}
		}
	}
	if dp[n][1].ok {
		return dp[n][1].str
	}
	return ""
}

func genEmail(rng *rand.Rand) string {
	parts := rng.Intn(3) + 1
	var sb strings.Builder
	for i := 0; i < parts; i++ {
		l := rng.Intn(5) + 1
		for j := 0; j < l; j++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		if i+1 < parts {
			sb.WriteByte('.')
		}
	}
	local := sb.String()
	sb.Reset()
	domParts := rng.Intn(2) + 1
	for i := 0; i < domParts; i++ {
		l := rng.Intn(5) + 1
		for j := 0; j < l; j++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		if i+1 < domParts {
			sb.WriteByte('.')
		}
	}
	domain := sb.String()
	return local + "@" + domain
}

func describe(email string) string {
	email = strings.ReplaceAll(email, ".", "dot")
	email = strings.ReplaceAll(email, "@", "at")
	return email
}

func generateCase(rng *rand.Rand) (string, string) {
	email := genEmail(rng)
	desc := describe(email)
	exp := solve(desc)
	return desc + "\n", exp
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
