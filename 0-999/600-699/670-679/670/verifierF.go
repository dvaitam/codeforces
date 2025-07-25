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

func numDigits(x int) int {
	if x == 0 {
		return 1
	}
	c := 0
	for x > 0 {
		c++
		x /= 10
	}
	return c
}

func expected(s, t string) string {
	n := len(s)
	length := 1
	for d := 1; d <= 7; d++ {
		l := n - d
		if l >= 1 && numDigits(l) == d {
			length = l
			break
		}
	}
	cnt := make([]int, 10)
	for i := 0; i < len(s); i++ {
		cnt[s[i]-'0']++
	}
	strL := fmt.Sprintf("%d", length)
	for i := 0; i < len(strL); i++ {
		cnt[strL[i]-'0']--
	}
	cntT := make([]int, 10)
	for i := 0; i < len(t); i++ {
		cntT[t[i]-'0']++
	}
	rem := make([]int, 10)
	for i := 0; i < 10; i++ {
		rem[i] = cnt[i] - cntT[i]
	}
	build := func(arr []int) string {
		first := int(t[0] - '0')
		var less, equal, greater []byte
		for d := 0; d < 10; d++ {
			for i := 0; i < arr[d]; i++ {
				if d < first {
					less = append(less, byte('0'+d))
				} else if d == first {
					equal = append(equal, byte('0'+d))
				} else {
					greater = append(greater, byte('0'+d))
				}
			}
		}
		cand1 := string(append(append(append([]byte{}, less...), []byte(t)...), append(equal, greater...)...))
		cand2 := string(append(append(append([]byte{}, less...), equal...), append([]byte(t), greater...)...))
		if cand1 < cand2 {
			return cand1
		}
		return cand2
	}
	ans := ""
	have := false
	if t[0] != '0' || length == 1 {
		var rest []byte
		for d := 0; d < 10; d++ {
			for i := 0; i < rem[d]; i++ {
				rest = append(rest, byte('0'+d))
			}
		}
		cand := t + string(rest)
		ans = cand
		have = true
	}
	for d := 1; d <= 9; d++ {
		if rem[d] > 0 {
			rem[d]--
			cand := string('0'+byte(d)) + build(rem)
			rem[d]++
			if !have || cand < ans {
				ans = cand
				have = true
			}
			break
		}
	}
	if !have {
		ans = t
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(1000000)
	s0 := fmt.Sprintf("%d", n)
	k := len(s0)
	digits := append([]byte(s0), []byte(fmt.Sprintf("%d", k))...)
	rng.Shuffle(len(digits), func(i, j int) { digits[i], digits[j] = digits[j], digits[i] })
	s := string(digits)
	start := rng.Intn(len(s0))
	end := rng.Intn(len(s0)-start) + start
	t := s0[start : end+1]
	input := s + "\n" + t + "\n"
	exp := expected(s, t)
	return input, exp
}

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\n got: %s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
