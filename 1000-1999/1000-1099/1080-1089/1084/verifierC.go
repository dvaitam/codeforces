package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const solution1084CSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 1000000007

func main() {
	reader := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for reader.Scan() {
		s := reader.Text()
		// prepend a 'b' to simplify calculation
		s = "b" + s
		var ans int64 = 1
		var cnt int64 = 0
		// iterate from end to start
		for i := len(s) - 1; i >= 0; i-- {
			switch s[i] {
			case 'a':
				cnt++
			case 'b':
				cnt++
				ans = ans * cnt % MOD
				cnt = 0
			}
		}
		// subtract the empty sequence
		res := (ans - 1) % MOD
		if res < 0 {
			res += MOD
		}
		fmt.Fprintln(writer, res)
	}
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1084CSource

const MOD int64 = 1000000007

var testcases = []string{
	"aaba",
	"cbbcacaccabcbccbcbcbaabbbbbcacaaaaabaaccbcccabbc",
	"bcbbbabccbccabbbccbcbbbccccbbcabca",
	"bbbbccccccccbbcabcbccabcaacaaccabcacacab",
	"aabcaabbaacaaaaa",
	"cab",
	"aacaccabcaaaaabcc",
	"cabbbabbcccabbccabaaccbaabaccbbcbabbbcbcaccacaba",
	"aaabcacca",
	"acbabacacccbbcbb",
	"aaabbaaccaaaaaaaaaaccbbbccbacacbbc",
	"cc",
	"bcca",
	"cbbacac",
	"bcbbacbaabacababcbbaccaa",
	"abbcaabaacbbcbabbab",
	"aacbcbac",
	"accbacc",
	"bbbabbbcbccbbcbbaabcbcbcccacccaa",
	"acabaaccaabbacccbbacbaa",
	"babcaccbaabcabbccabbcacbcabcabbbbcc",
	"ccbbcbcccacaccabcbcabbccbca",
	"abacccb",
	"abcabcccabcaac",
	"bbbccabcaababcccacbcbbbcbaacbabbcccaaabbbcba",
	"b",
	"babacbcbacbbcabcccbabbbbcbccbcaccabc",
	"accabaabcacbaaacbbabbabbbabbbbacbc",
	"bbaaabacbabaaa",
	"ccbbcbbccccccbcbcc",
	"cbccccabbbbabbbbaabbaaababaccbacbacbcbbca",
	"babaaabbabbbabaacaccabbc",
	"caaaabbbabccabcccbaabbbcbcccaccccbbcbabbcbaabcaacc",
	"abbbcabacbaa",
	"ab",
	"abcbabaabacaacccaaacacbccaabbaacaaabaacacbb",
	"cbccccacaacaababbbababbbacabbbcbbabcaccc",
	"accaaab",
	"acacacbccbcacaacbacaaacccaacbcbaabcbcc",
	"bcac",
	"caabb",
	"baabc",
	"baababbcabccba",
	"bcaaaaccccbbccbbabb",
	"cbbcbcbbcacbbc",
	"bcbcacbabccbccccabb",
	"a",
	"cacaabbccbabaccbcbccabacbabccabcacc",
	"abcbccbcccbaa",
	"bcbcacbccbbacccbaaabccac",
	"abbbcaabacccabcbabaabaac",
	"accbabcabbbcbaab",
	"cbacacaaabbbabbcaccacbcaabaabbac",
	"abbababacaa",
	"caabaacabbccbabbacbacbbabcb",
	"baaabbbaaccabcbcccbaaccaa",
	"aaaaccacc",
	"bcaaaccc",
	"accaaaabaaabbbbacbaaaccccaaabbcba",
	"accac",
	"cc",
	"babcabcccbbbaabbabcbcbba",
	"abcccaccabcbcaccaaa",
	"cbbaaaccaaabcabc",
	"aabcaabcbbbacacc",
	"babcac",
	"cbcabaccbbbac",
	"ccbcbacccacabbbacb",
	"abbc",
	"aacbcbbbcccca",
	"caccbbcbc",
	"cbaccbb",
	"bb",
	"aacacaaa",
	"bbaabbcaacacbacbbababbc",
	"cbaaacaccac",
	"caac",
	"bcb",
	"c",
	"ca",
	"cab",
	"ccabbbaccccaa",
	"cabaaabcbcbabcaccacaca",
	"bbaca",
	"acacbaccbaaccb",
	"accb",
	"aba",
	"babbabacb",
	"aabac",
	"cccbb",
	"a",
	"bcaababbccbbbbac",
	"caccacc",
	"bbbbaababbaccbabbcba",
	"abcbabab",
	"bccbaba",
	"aaaaaca",
	"bccba",
	"cccbbccbabbbbbcacca",
	"caaabbbbaccccc",
	"acacbbcacbb",
	"cacacbacaaabccccccaccb",
}

func solveCase(s string) int64 {
	ans := int64(1)
	cnt := int64(0)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 'a':
			cnt++
		case 'b':
			ans = ans * (cnt + 1) % MOD
			cnt = 0
		}
	}
	ans = ans * (cnt + 1) % MOD
	res := ans - 1
	if res < 0 {
		res += MOD
	}
	return res % MOD
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	cmd = exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for caseNum, s := range testcases {
		input := s + "\n"
		want := fmt.Sprintf("%d", solveCase(s))
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
