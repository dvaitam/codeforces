package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"
)

func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func mentions(text, user string) bool {
	for i := 0; ; {
		idx := strings.Index(text[i:], user)
		if idx == -1 {
			return false
		}
		idx += i
		beforeOK := idx == 0 || !isWordChar(rune(text[idx-1]))
		afterIdx := idx + len(user)
		afterOK := afterIdx == len(text) || !isWordChar(rune(text[afterIdx]))
		if beforeOK && afterOK {
			return true
		}
		i = idx + 1
	}
}

type chatCase struct {
	names  []string
	prefix []string
	texts  []string
}

func randomWord(rng *rand.Rand) string {
	l := rng.Intn(5) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func randomChat(rng *rand.Rand) chatCase {
	n := rng.Intn(3) + 2
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = randomWord(rng)
	}
	m := rng.Intn(5) + 1
	prefix := make([]string, m)
	texts := make([]string, m)
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			prefix[i] = "?"
		} else {
			prefix[i] = names[rng.Intn(n)]
		}
		words := rng.Intn(4) + 1
		var sb strings.Builder
		for j := 0; j < words; j++ {
			w := randomWord(rng)
			if rng.Intn(4) == 0 {
				w = names[rng.Intn(n)]
			}
			sb.WriteString(w)
			if j+1 < words {
				sb.WriteByte(' ')
			}
		}
		texts[i] = sb.String()
		if prefix[i] != "?" && mentions(texts[i], prefix[i]) {
			texts[i] = randomWord(rng)
		}
	}
	return chatCase{names, prefix, texts}
}

func solveChat(cc chatCase) ([]string, bool) {
	n := len(cc.names)
	m := len(cc.prefix)
	nameIndex := make(map[string]int)
	for i, nm := range cc.names {
		nameIndex[nm] = i
	}
	allowed := make([][]int, m)
	for i := 0; i < m; i++ {
		if cc.prefix[i] != "?" {
			if id, ok := nameIndex[cc.prefix[i]]; ok {
				allowed[i] = []int{id}
			} else {
				return nil, false
			}
		} else {
			cand := []int{}
			for j, nm := range cc.names {
				if !mentions(cc.texts[i], nm) {
					cand = append(cand, j)
				}
			}
			if len(cand) == 0 {
				return nil, false
			}
			allowed[i] = cand
		}
	}
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}
	for _, u := range allowed[0] {
		dp[0][u] = -2
	}
	for i := 1; i < m; i++ {
		for _, u := range allowed[i] {
			for v := 0; v < n; v++ {
				if dp[i-1][v] != -1 && v != u {
					dp[i][u] = v
					break
				}
			}
		}
	}
	endUser := -1
	for u := 0; u < n; u++ {
		if dp[m-1][u] != -1 {
			endUser = u
			break
		}
	}
	if endUser == -1 {
		return nil, false
	}
	resIdx := make([]int, m)
	cur := endUser
	for i := m - 1; i >= 0; i-- {
		resIdx[i] = cur
		cur = dp[i][cur]
	}
	res := make([]string, m)
	for i := 0; i < m; i++ {
		res[i] = cc.names[resIdx[i]] + ":" + cc.texts[i]
	}
	return res, true
}

func verifyOutput(cc chatCase, outLines []string) error {
	if len(outLines) == 1 && strings.TrimSpace(outLines[0]) == "Impossible" {
		if _, ok := solveChat(cc); ok {
			return fmt.Errorf("should be possible")
		}
		return nil
	}
	if len(outLines) != len(cc.prefix) {
		return fmt.Errorf("wrong number of lines")
	}
	prev := ""
	for i, line := range outLines {
		idx := strings.IndexByte(line, ':')
		if idx == -1 {
			return fmt.Errorf("line %d missing colon", i+1)
		}
		user := line[:idx]
		text := line[idx+1:]
		found := false
		for _, nm := range cc.names {
			if nm == user {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown user %s", user)
		}
		if cc.prefix[i] != "?" && cc.prefix[i] != user {
			return fmt.Errorf("line %d wrong user", i+1)
		}
		if user == prev {
			return fmt.Errorf("line %d same as previous", i+1)
		}
		if mentions(text, user) {
			return fmt.Errorf("user %s mentions himself", user)
		}
		prev = user
	}
	return nil
}

func runCase(bin string, cc chatCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cc.names))
	for i, nm := range cc.names {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(nm)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", len(cc.prefix))
	for i := 0; i < len(cc.prefix); i++ {
		sb.WriteString(cc.prefix[i])
		sb.WriteByte(':')
		sb.WriteString(cc.texts[i])
		sb.WriteByte('\n')
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	return verifyOutput(cc, lines)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []chatCase{}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomChat(rng))
	}
	for i, cc := range cases {
		if err := runCase(bin, cc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
