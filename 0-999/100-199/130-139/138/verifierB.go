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

type state struct {
	c1    [10]uint8
	c2    [10]uint8
	carry uint8
}

func encode(st state) string {
	var b bytes.Buffer
	for i := 0; i < 10; i++ {
		b.WriteByte(st.c1[i])
	}
	for i := 0; i < 10; i++ {
		b.WriteByte(st.c2[i])
	}
	b.WriteByte(st.carry)
	return b.String()
}

func dfs(c1, c2 [10]uint8, carry uint8, memo map[string]int) int {
	st := state{c1, c2, carry}
	key := encode(st)
	if v, ok := memo[key]; ok {
		return v
	}
	best := 0
	for i := 0; i < 10; i++ {
		if c1[i] == 0 {
			continue
		}
		c1[i]--
		for j := 0; j < 10; j++ {
			if c2[j] == 0 {
				continue
			}
			c2[j]--
			if int(i)+int(j)+int(carry)%10 == 0 {
				newCarry := uint8((int(i) + int(j) + int(carry)) / 10)
				val := 1 + dfs(c1, c2, newCarry, memo)
				if val > best {
					best = val
				}
			}
			c2[j]++
		}
		c1[i]++
	}
	memo[key] = best
	return best
}

func maxTrailingZeros(s string) int {
	var c1, c2 [10]uint8
	for i := 0; i < len(s); i++ {
		d := s[i] - '0'
		c1[d]++
		c2[d]++
	}
	memo := make(map[string]int)
	return dfs(c1, c2, 0, memo)
}

func trailingZerosOfSum(a, b string) int {
	i := len(a) - 1
	j := len(b) - 1
	carry := 0
	zeros := 0
	for i >= 0 || j >= 0 {
		da, db := 0, 0
		if i >= 0 {
			da = int(a[i] - '0')
			i--
		}
		if j >= 0 {
			db = int(b[j] - '0')
			j--
		}
		v := da + db + carry
		d := v % 10
		carry = v / 10
		if d == 0 {
			zeros++
		} else {
			return zeros
		}
	}
	if carry%10 == 0 && carry > 0 {
		zeros++
	}
	return zeros
}

func generateCase(rng *rand.Rand) (string, int) {
	l := rng.Intn(6) + 1
	var sb strings.Builder
	digits := make([]byte, l)
	for i := 0; i < l; i++ {
		digits[i] = byte('0' + rng.Intn(10))
	}
	s := string(digits)
	fmt.Fprintf(&sb, "%s\n", s)
	exp := maxTrailingZeros(s)
	return sb.String(), exp
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, in)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) < 2 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected two lines of output\n", t)
			os.Exit(1)
		}
		a := strings.TrimSpace(lines[0])
		b := strings.TrimSpace(lines[1])
		if len(a) != len(strings.TrimSpace(in))-1 || len(b) != len(strings.TrimSpace(in))-1 {
			fmt.Fprintf(os.Stderr, "case %d failed: output length mismatch\n", t)
			os.Exit(1)
		}
		var cntOrig [10]int
		for i := 0; i < len(in)-1; i++ {
			cntOrig[in[i]-'0']++
		}
		var cntA, cntB [10]int
		for i := 0; i < len(a); i++ {
			if a[i] < '0' || a[i] > '9' {
				fmt.Fprintf(os.Stderr, "case %d failed: invalid digit in output\n", t)
				os.Exit(1)
			}
			cntA[a[i]-'0']++
		}
		for i := 0; i < len(b); i++ {
			if b[i] < '0' || b[i] > '9' {
				fmt.Fprintf(os.Stderr, "case %d failed: invalid digit in output\n", t)
				os.Exit(1)
			}
			cntB[b[i]-'0']++
		}
		for d := 0; d < 10; d++ {
			if cntA[d] != cntOrig[d] || cntB[d] != cntOrig[d] {
				fmt.Fprintf(os.Stderr, "case %d failed: output is not a permutation\n", t)
				os.Exit(1)
			}
		}
		got := trailingZerosOfSum(a, b)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d trailing zeros, got %d\ninput:\n%soutput:\n%s", t, exp, got, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
