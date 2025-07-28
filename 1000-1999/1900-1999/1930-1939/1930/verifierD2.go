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

func computeSum(s string) int64 {
	var cnt [5]int64
	var sum [5]int64
	var total int64
	for i := 0; i < len(s); i++ {
		ch := s[i]
		var newCnt [5]int64
		var newSum [5]int64
		if ch == '0' {
			newCnt[0] = 1
		} else {
			newCnt[1] = 1
			newSum[1] = 1
		}
		for st := 0; st < 5; st++ {
			c := cnt[st]
			if c == 0 {
				continue
			}
			if ch == '0' {
				switch st {
				case 0:
					newCnt[0] += c
					newSum[0] += sum[st]
				case 1:
					newCnt[2] += c
					newSum[2] += sum[st]
				case 2:
					newCnt[3] += c
					newSum[3] += sum[st]
				default:
					newCnt[4] += c
					newSum[4] += sum[st]
				}
			} else {
				switch st {
				case 0:
					newCnt[1] += c
					newSum[1] += sum[st] + c
				case 1:
					newCnt[2] += c
					newSum[2] += sum[st]
				case 2:
					newCnt[3] += c
					newSum[3] += sum[st]
				case 3, 4:
					newCnt[1] += c
					newSum[1] += sum[st] + c
				}
			}
		}
		cnt = newCnt
		sum = newSum
		for j := 0; j < 5; j++ {
			total += sum[j]
		}
	}
	return total
}

func generateCaseD2(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	bytes := make([]byte, n)
	for i := range bytes {
		if rng.Intn(2) == 0 {
			bytes[i] = '0'
		} else {
			bytes[i] = '1'
		}
	}
	s := string(bytes)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n%s\n", n, s))
	input := sb.String()
	exp := fmt.Sprintf("%d\n", computeSum(s))
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
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseD2(rng)
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
