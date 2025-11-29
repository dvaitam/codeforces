package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
ea eeaa
dbad ddbbaaad
dcbeac dcccbbbeaaccc
da d
dc bdccc
cacede cccacceedddee
d dd
bcec bbbcccec
dcdad dddcccdddaaadd
addbeaa addddbbeaaa
eabbcb eb
dcededca dddcccdeeeedccaaa
beab bbbbbbbaaaaaaa
ccdeedcb cccccccccdddddddddeeeeeeeeeeebbbbbb
cbeea cccceeeebbbbb
dbdaeb dddbbbbbbaaa
c c
bec cdde
abddbbaacc aaabddbbbbaaccc
ebdb bbbbdddde
ebdcac eeeebbcccaaa
eecce eeeeeeeccc
acdbd bb
a b
eabaeaa eeeeeeaaaaa
eddecbbbba cdddddeeeecbbbba
cabd ccccaaaabbdddd
adacb abb
ecca eeeeccccaaaaaa
bcdaaeaaaa bbcdaaaaeeaaaa
dbceea dddbbbcccceeeaa
bacdacb bbbbbbaaaaccddddaaaacccbbb
ceceacbed ccccceeeeebbbbbbddddaa
dda cccccccc
b bbbbb
ddadbe dccee
adbaaeccaa dcccee
acdbccccab acccccccdbcccccccaaaaaabbbbbbb
eccbccde ccccccc
badcbad bbbbbaaadddddccccccbbbbbaaaaaddd
ebcda bbbbcccdaaaaa
aeecaa eeeeeeeccaaaaa
cbbd ccccbbbdddd
becaaa bbbb
ecccbae eeecccccbaae
bad bbbbbaaaaaadd
edabeeabb edbbeeeaaabb
eaba cddeaa
abbdcccaea abbbcccccccccaaeeeeeaaa
becddeecc aaaaa
deeeac daaabbbccc
bbd eee
adcad ddddcccaaa
dcaa dad
caee dddaaeee
aaecccaa bb
aaaad bbbbbaaaaa
ddaadcbbcd dddaaacccccbbbcd
eaeaa eee
ebc cb
baccbaeaaaa bbcccbaeeaaaa
abbabbcba aaaabbbaabbbbccca
eccdcbac eeeeecccccddbbaaaaccc
dabdb de
abbbecbbb ccccbbbbbaaaaaaccc
dddbdcc dddbddccccccc
decbbaecea dccbbeeceeaa
ccecaa cceecaaaaa
ddeaeeadcc cccc
eddebddaaa eeeeedddddbbbddaaaa
bbabceecacc bbabceeeeaacc
abcbadbca abbbbbbcaaa
bde b
ceabaaae ccccccaaaaaee
bebabb eeeee
beabdadcca cccccaaaaaaa
cccdee de
aeeca b
ccd ed
dcbbadbab cbcbddd
cbddcdad eeeeeebbbbbbbaaaaaaadd
cbbacca cbbaaaccccccccaaaaaaa
eeaddd eeecccaaaaddd
eda bbbcccccaaaaddddd
eabeee eeeeeeeaaa
eacadd bdcc
dddcabcccab dddcccabcccccbbb
cbaca ccccbbbaaaaacc
abcbcccbcccc abc
eaccdaccbe bcccccccaaaaa
cdcbaab bbb
eacbaaa bc
ccedcaaac eeeeee
cbdc cb
e eeee
a a
dcedbecaec cd
d cccc
bbdbbbbcdd bbdbbbbcdd
dbcaadbeee dbcaadbeee`

// Embedded reference logic from 1185B.go.
func referenceMatch(s, t string) bool {
	i, j := 0, 0
	n, m := len(s), len(t)
	for i < n && j < m {
		if s[i] != t[j] {
			return false
		}
		ch := s[i]
		cs, ct := 0, 0
		for i+cs < n && s[i+cs] == ch {
			cs++
		}
		for j+ct < m && t[j+ct] == ch {
			ct++
		}
		if ct < cs {
			return false
		}
		i += cs
		j += ct
	}
	return i == n && j == m
}

func runCase(bin, s, t string) error {
	expect := "NO"
	if referenceMatch(s, t) {
		expect = "YES"
	}
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(s)
	input.WriteByte(' ')
	input.WriteString(t)
	input.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scanner.Text())
	caseNum := 0
	for caseNum < t {
		if !scanner.Scan() {
			fmt.Printf("case %d: missing s\n", caseNum+1)
			os.Exit(1)
		}
		s := scanner.Text()
		if !scanner.Scan() {
			fmt.Printf("case %d: missing t\n", caseNum+1)
			os.Exit(1)
		}
		tStr := scanner.Text()
		if err := runCase(bin, s, tStr); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum+1, err)
			os.Exit(1)
		}
		caseNum++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", caseNum)
}
