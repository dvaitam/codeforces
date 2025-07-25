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

func solveC(s string) string {
	n := len(s)
	sum := 0
	idxMod := make([][]int, 3)
	for i := 0; i < n; i++ {
		d := int(s[i] - '0')
		sum += d
		idxMod[d%3] = append(idxMod[d%3], i)
	}
	p := sum % 3
	remove := make([]bool, n)
	switch p {
	case 1:
		if len(idxMod[1]) >= 1 {
			pos := idxMod[1][len(idxMod[1])-1]
			remove[pos] = true
		} else if len(idxMod[2]) >= 2 {
			l := len(idxMod[2])
			remove[idxMod[2][l-1]] = true
			remove[idxMod[2][l-2]] = true
		} else {
			return "-1"
		}
	case 2:
		if len(idxMod[2]) >= 1 {
			pos := idxMod[2][len(idxMod[2])-1]
			remove[pos] = true
		} else if len(idxMod[1]) >= 2 {
			l := len(idxMod[1])
			remove[idxMod[1][l-1]] = true
			remove[idxMod[1][l-2]] = true
		} else {
			return "-1"
		}
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if !remove[i] {
			sb.WriteByte(s[i])
		}
	}
	t := sb.String()
	if len(t) == 0 {
		return "-1"
	}
	i := 0
	for i < len(t) && t[i] == '0' {
		i++
	}
	if i == len(t) {
		return "0"
	}
	return t[i:]
}

func genCaseC(rng *rand.Rand) (string, string) {
	l := rng.Intn(20) + 1
	digits := make([]byte, l)
	digits[0] = byte(rng.Intn(9) + 1 + '0')
	for i := 1; i < l; i++ {
		digits[i] = byte(rng.Intn(10) + '0')
	}
	s := string(digits)
	exp := solveC(s)
	return s + "\n", exp
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCaseC(bin, in, exp string) error {
	out, err := runCandidate(bin, []byte(in))
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %s got %s", exp, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseC(rng)
		if err := runCaseC(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
