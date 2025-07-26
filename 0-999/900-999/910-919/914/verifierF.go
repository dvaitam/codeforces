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

type query struct {
	typ int
	a   int
	b   int
	s   string
}

type testCase struct {
	str string
	qs  []query
}

func expected(tc testCase) []int {
	s := []byte(tc.str)
	res := []int{}
	for _, q := range tc.qs {
		if q.typ == 1 {
			s[q.a-1] = q.s[0]
		} else {
			l, r := q.a, q.b
			pat := q.s
			cnt := 0
			for i := l - 1; i+len(pat) <= r; i++ {
				match := true
				for j := 0; j < len(pat); j++ {
					if s[i+j] != pat[j] {
						match = false
						break
					}
				}
				if match {
					cnt++
				}
			}
			res = append(res, cnt)
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, testCase) {
	length := rng.Intn(8) + 1
	bytesArr := make([]byte, length)
	for i := range bytesArr {
		bytesArr[i] = byte('a' + rng.Intn(3))
	}
	str := string(bytesArr)
	q := rng.Intn(5) + 1
	qs := make([]query, q)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\n%d\n", str, q))
	for i := 0; i < q; i++ {
		typ := rng.Intn(2) + 1
		if typ == 1 {
			pos := rng.Intn(length) + 1
			ch := byte('a' + rng.Intn(3))
			sb.WriteString(fmt.Sprintf("1 %d %c\n", pos, ch))
			qs[i] = query{typ: 1, a: pos, s: string(ch)}
			s := []byte(str)
			s[pos-1] = ch
			str = string(s)
		} else {
			l := rng.Intn(length) + 1
			r := rng.Intn(length-l+1) + l
			patLen := rng.Intn(3) + 1
			patBytes := make([]byte, patLen)
			for j := range patBytes {
				patBytes[j] = byte('a' + rng.Intn(3))
			}
			pat := string(patBytes)
			sb.WriteString(fmt.Sprintf("2 %d %d %s\n", l, r, pat))
			qs[i] = query{typ: 2, a: l, b: r, s: pat}
		}
	}
	return sb.String(), testCase{str: string(bytesArr), qs: qs}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, tc := generateCase(rng)
		expSlice := expected(tc)
		expStrings := make([]string, len(expSlice))
		for j, v := range expSlice {
			expStrings[j] = fmt.Sprintf("%d", v)
		}
		exp := strings.Join(expStrings, "\n")
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
