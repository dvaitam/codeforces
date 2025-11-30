package main

import (
	"bufio"
	"bytes"
	"fmt"
	"hash/fnv"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const embeddedTestcasesH = `interactive test 1
interactive test 2
interactive test 3
interactive test 4
interactive test 5
interactive test 6
interactive test 7
interactive test 8
interactive test 9
interactive test 10
interactive test 11
interactive test 12
interactive test 13
interactive test 14
interactive test 15
interactive test 16
interactive test 17
interactive test 18
interactive test 19
interactive test 20
interactive test 21
interactive test 22
interactive test 23
interactive test 24
interactive test 25
interactive test 26
interactive test 27
interactive test 28
interactive test 29
interactive test 30
interactive test 31
interactive test 32
interactive test 33
interactive test 34
interactive test 35
interactive test 36
interactive test 37
interactive test 38
interactive test 39
interactive test 40
interactive test 41
interactive test 42
interactive test 43
interactive test 44
interactive test 45
interactive test 46
interactive test 47
interactive test 48
interactive test 49
interactive test 50
interactive test 51
interactive test 52
interactive test 53
interactive test 54
interactive test 55
interactive test 56
interactive test 57
interactive test 58
interactive test 59
interactive test 60
interactive test 61
interactive test 62
interactive test 63
interactive test 64
interactive test 65
interactive test 66
interactive test 67
interactive test 68
interactive test 69
interactive test 70
interactive test 71
interactive test 72
interactive test 73
interactive test 74
interactive test 75
interactive test 76
interactive test 77
interactive test 78
interactive test 79
interactive test 80
interactive test 81
interactive test 82
interactive test 83
interactive test 84
interactive test 85
interactive test 86
interactive test 87
interactive test 88
interactive test 89
interactive test 90
interactive test 91
interactive test 92
interactive test 93
interactive test 94
interactive test 95
interactive test 96
interactive test 97
interactive test 98
interactive test 99
interactive test 100`

// hashString mirrors the interactive judge: hash(s) = sum(ord(s[i]) * p^i) mod m.
func hashString(s string, p, m int64) int64 {
	res := int64(0)
	mul := int64(1)
	for i := 0; i < len(s); i++ {
		val := int64(s[i] - 'a' + 1)
		res = (res + val*mul) % m
		mul = (mul * p) % m
	}
	return res
}

// referenceSolve embeds the reference solution from 1994H.go and is used as a sanity
// check that our generated (p, m) pairs are solvable with the intended strategy.
func referenceSolve(p, m int64) (int64, int64) {
	pResp := hashString("aa", p, m)
	xResp := hashString("zzzzzzzzzz", p, m)
	hash := func(s string) int64 {
		return hashString(s, p, m)
	}
	return embedded1994HSolver(pResp, xResp, hash)
}

// embedded1994HSolver mirrors the logic in 1994H.go but runs fully in memory.
func embedded1994HSolver(pResp, xResp int64, hash func(string) int64) (int64, int64) {
	pVal := pResp - 1 // solver decrements immediately after reading
	hs := int64(0)
	y := xResp + 1
	var o int64
	an := make([]int64, 11)

	for i := 0; i < 10; i++ {
		hs = hs*pVal + 26
		an[i] = 26 - y%pVal
		y /= pVal
	}

	sBytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		if an[i] < 1 {
			an[i] = 26
			an[i+1]--
		}
		sBytes[i] = byte('a' + an[i] - 1)
	}
	s := string(sBytes)

	for i := 9; i >= 0; i-- {
		o = o*pVal + int64(sBytes[i]-'a'+1)
	}

	mResp := hash(s)
	ans := hs - xResp - o + mResp
	return pVal, ans
}

func deriveParams(seed string) (int64, int64) {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte(seed))
	h := hasher.Sum64()
	p := int64(27 + h%24)                   // 27..50 satisfies 26 < p <= 50
	m := int64(1_000_000_000 + h%1_000_000) // <= 1_000_999_999 < 2e9
	if m <= p+1 {
		m = p + 2
	}
	return p, m
}

func parseTestcases() []string {
	lines := strings.Split(embeddedTestcasesH, "\n")
	var out []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func nextToken(scanner *bufio.Scanner, context string) (string, error) {
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("%s: failed to read token: %w", context, err)
	}
	return "", fmt.Errorf("%s: unexpected EOF", context)
}

func asInt(tok, context string) (int64, error) {
	val, err := strconv.ParseInt(tok, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: expected integer, got %q", context, tok)
	}
	return val, nil
}

func runCase(bin string, p, m int64, idx int) error {
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("case %d: stdout pipe failed: %w", idx, err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("case %d: stdin pipe failed: %w", idx, err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("case %d: failed to start: %w", idx, err)
	}

	timer := time.AfterFunc(5*time.Second, func() {
		_ = cmd.Process.Kill()
	})
	done := false
	defer func() {
		timer.Stop()
		if !done {
			_ = cmd.Process.Kill()
			cmd.Wait()
		}
	}()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 0, 1024), 1<<20)

	queries := 0
	for {
		tok, err := nextToken(scanner, fmt.Sprintf("case %d", idx))
		if err != nil {
			return fmt.Errorf("%v (stderr: %s)", err, stderr.String())
		}
		switch tok {
		case "?":
			query, err := nextToken(scanner, fmt.Sprintf("case %d after '?'", idx))
			if err != nil {
				return fmt.Errorf("%v (stderr: %s)", err, stderr.String())
			}
			queries++
			if queries > 3 {
				return fmt.Errorf("case %d: too many queries (>3). stderr: %s", idx, stderr.String())
			}
			answer := hashString(query, p, m)
			if _, err := fmt.Fprintf(stdin, "%d\n", answer); err != nil {
				return fmt.Errorf("case %d: failed to send reply: %w", idx, err)
			}
		case "!":
			baseTok, err := nextToken(scanner, fmt.Sprintf("case %d after '!'", idx))
			if err != nil {
				return fmt.Errorf("%v (stderr: %s)", err, stderr.String())
			}
			modTok, err := nextToken(scanner, fmt.Sprintf("case %d after '!'", idx))
			if err != nil {
				return fmt.Errorf("%v (stderr: %s)", err, stderr.String())
			}
			guessP, err := asInt(baseTok, fmt.Sprintf("case %d", idx))
			if err != nil {
				return err
			}
			guessM, err := asInt(modTok, fmt.Sprintf("case %d", idx))
			if err != nil {
				return err
			}
			_ = stdin.Close()
			if err := cmd.Wait(); err != nil {
				return fmt.Errorf("case %d: program exited with error after answer: %v (stderr: %s)", idx, err, stderr.String())
			}
			done = true
			if guessP != p || guessM != m {
				return fmt.Errorf("case %d: wrong answer, expected %d %d got %d %d", idx, p, m, guessP, guessM)
			}
			if queries == 0 {
				return fmt.Errorf("case %d: no queries made before answering", idx)
			}
			return nil
		default:
			return fmt.Errorf("case %d: unexpected token %q (stderr: %s)", idx, tok, stderr.String())
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	if len(testcases) == 0 {
		fmt.Println("no embedded testcases found")
		os.Exit(1)
	}

	for i, seed := range testcases {
		p, m := deriveParams(seed)
		refP, refM := referenceSolve(p, m)
		if refP != p || refM != m {
			fmt.Fprintf(os.Stderr, "internal error: reference solution failed on case %d\n", i+1)
			os.Exit(1)
		}
		if err := runCase(bin, p, m, i+1); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
