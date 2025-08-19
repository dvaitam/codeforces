package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCaseB struct {
	s string
}

func parseTestcases(path string) ([]testCaseB, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewScanner(f)
	if !in.Scan() {
		return nil, fmt.Errorf("empty file")
	}
	var T int
	fmt.Sscan(in.Text(), &T)
	cases := make([]testCaseB, 0, T)
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		cases = append(cases, testCaseB{s: line})
		if len(cases) == T {
			break
		}
	}
	if err := in.Err(); err != nil {
		return nil, err
	}
	if len(cases) != T {
		return nil, fmt.Errorf("expected %d cases got %d", T, len(cases))
	}
	return cases, nil
}

func expectedAnswer(s string) string {
	cnt := map[rune]int{'R': 0, 'P': 0, 'S': 0}
	for _, ch := range s {
		cnt[ch]++
	}
	var ans strings.Builder
	if cnt['R'] == cnt['P'] && cnt['P'] == cnt['S'] {
		for _, ch := range s {
			switch ch {
			case 'R':
				ans.WriteByte('P')
			case 'P':
				ans.WriteByte('S')
			case 'S':
				ans.WriteByte('R')
			}
		}
	} else {
		cmx := 'R'
		maxc := cnt['R']
		if cnt['P'] > maxc {
			cmx = 'P'
			maxc = cnt['P']
		}
		if cnt['S'] > maxc {
			cmx = 'S'
			maxc = cnt['S']
		}
		var play byte
		switch cmx {
		case 'R':
			play = 'P'
		case 'P':
			play = 'S'
		case 'S':
			play = 'R'
		}
		for i := 0; i < len(s); i++ {
			ans.WriteByte(play)
		}
	}
	return ans.String()
}

func run(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return out.String(), errBuf.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		input := fmt.Sprintf("1\n%s\n", tc.s)
		outStr, errStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", idx+1, err, errStr)
			os.Exit(1)
		}
        ans := strings.TrimSpace(outStr)
        expected := expectedAnswer(tc.s)
        if ans != expected {
            // Also accept uniform answers using a best-response move to any most frequent symbol
            if !acceptUniformBest(tc.s, ans) {
                fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, ans)
                os.Exit(1)
            }
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}

func acceptUniformBest(s, ans string) bool {
    if len(ans) != len(s) || len(ans) == 0 {
        return false
    }
    // Ensure ans is uniform and consists of valid chars
    ch := ans[0]
    if ch != 'R' && ch != 'P' && ch != 'S' {
        return false
    }
    for i := 1; i < len(ans); i++ {
        if ans[i] != ch {
            return false
        }
    }
    // Count frequencies in s
    cntR, cntP, cntS := 0, 0, 0
    for i := 0; i < len(s); i++ {
        switch s[i] {
        case 'R':
            cntR++
        case 'P':
            cntP++
        case 'S':
            cntS++
        default:
            return false
        }
    }
    maxc := cntR
    if cntP > maxc { maxc = cntP }
    if cntS > maxc { maxc = cntS }
    // Allowed uniform plays are counters to any symbol achieving max frequency
    allowed := map[byte]bool{}
    if cntR == maxc { allowed['P'] = true }
    if cntP == maxc { allowed['S'] = true }
    if cntS == maxc { allowed['R'] = true }
    return allowed[ch]
}
