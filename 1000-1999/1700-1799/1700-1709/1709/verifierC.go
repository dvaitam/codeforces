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

func runCandidate(bin, input string) (string, error) {
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

func isValid(s []byte) bool {
	bal := 0
	for _, c := range s {
		if c == '(' {
			bal++
		} else {
			bal--
		}
		if bal < 0 {
			return false
		}
	}
	return bal == 0
}

func solveCase(str string) string {
	n := len(str)
	openCnt, closeCnt := 0, 0
	for i := 0; i < n; i++ {
		if str[i] == '(' {
			openCnt++
		} else if str[i] == ')' {
			closeCnt++
		}
	}
	openNeed := n/2 - openCnt
	closeNeed := n/2 - closeCnt
	bytesArr := []byte(str)
	openPos := []int{}
	closePos := []int{}
	for i := 0; i < n; i++ {
		if bytesArr[i] == '?' {
			if openNeed > 0 {
				bytesArr[i] = '('
				openNeed--
				openPos = append(openPos, i)
			} else {
				bytesArr[i] = ')'
				closeNeed--
				closePos = append(closePos, i)
			}
		}
	}
	if len(openPos) == 0 || len(closePos) == 0 {
		if isValid(bytesArr) {
			return "YES"
		}
		return "NO"
	}
	i := openPos[len(openPos)-1]
	j := closePos[0]
	bytesArr[i] = ')'
	bytesArr[j] = '('
	if isValid(bytesArr) {
		return "NO"
	}
	return "YES"
}

func randRBS(rng *rand.Rand, n int) string {
	open := n / 2
	close := n / 2
	bal := 0
	b := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		if open == 0 {
			b = append(b, ')')
			close--
			bal--
		} else if bal == 0 {
			b = append(b, '(')
			open--
			bal++
		} else {
			if rng.Intn(open+close) < open {
				b = append(b, '(')
				open--
				bal++
			} else {
				b = append(b, ')')
				close--
				bal--
			}
		}
	}
	return string(b)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(9)*2 + 2 // even between 2 and 20
	base := randRBS(rng, n)
	bytesArr := []byte(base)
	hasQ := false
	for i := 0; i < n; i++ {
		if rng.Intn(3) == 0 { // replace with ?
			bytesArr[i] = '?'
			hasQ = true
		}
	}
	if !hasQ {
		bytesArr[rng.Intn(n)] = '?'
	}
	s := string(bytesArr)
	input := fmt.Sprintf("1\n%s\n", s)
	expect := solveCase(s)
	return input, expect
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
