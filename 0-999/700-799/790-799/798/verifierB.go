package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const testcasesBRaw = `2 bbc aac
1 acaa
4 cbc acb bcb aab
3 abbab cbbab bbaaa
3 c c c
2 babb aabb
5 bb cc ab ab ac
4 baaba cbcca bcbac bacbb
1 b
1 b
4 ccbba cbbbc acaaa bbcba
5 aa bb bc ac cc
2 abccb cbbbb
5 aa ca bb cb ac
3 cbc bab bcc
3 cabb acba caab
4 acac accc baba acab
3 ac ba ab
4 bab bcc bab cac
3 aa bb bb
4 caab acca baac aaba
2 cb ca
5 bbb acc ccc ccb cbb
3 abba aaab abca
2 b c
2 bcaba bcaca
3 cbb bac aab
1 cb
2 bc cc
3 cbb bbc aab
5 a b b c b
1 bbba
1 a
1 ba
4 abcba aabaa bbbba cbcca
3 ccbc cbba aacc
1 cb
2 bc ca
5 baabc baccb acccc cbaca acbcb
2 cbb caa
1 cbbcc
1 ba
2 aaabc caaca
2 bcca cbac
1 cbbbb
1 cc
2 cbcb caba
5 cabb baaa bcca bbcc ccab
3 ba bc ca
2 abc aaa
1 a
3 a b c
4 ab ba bb ac
3 babb cccc cbac
1 baac
2 cccca abaaa
1 cc
5 ca bc ac cb ca
5 abbcc bacab abbca bacba cbbbb
4 b c b c
4 c c c b
1 bcab
3 aabc abcc ccbc
2 cbbba acbac
5 aab cac aca bbb cbb
4 cbaa acab abac bccb
5 ac bb ac aa cc
2 acccc bbbbb
1 cbbc
5 bbacb aacac ccbbb acbcc aabca
2 ccb acc
2 ababa cbacc
4 caaab abcbb bbbcb ababa
3 bb bc cc
4 c b b c
2 cbbab bcaaa
5 baa aca ccb aba bba
3 caca bbca aaab
1 acccc
4 bc aa ca ba
3 ccb bcb abc
2 c c
2 abaac bcabc
2 cab cac
1 baaab
5 cbb bcb caa baa caa
3 aabca cbbac cabab
1 bccb
3 bbaab babab aaaba
2 ccca bbcc
4 aacb bacc aabb abaa
5 aab bac caa aac cac
3 b c c
4 a a c b
3 ab aa bb
3 abcaa cacbc aaaab
5 a a a c c
2 cb ab
4 cbcaa bcacb aabcc babba
4 abac bbca bbcb cacc
`

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "798B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)


	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("test %d invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != n+1 {
			fmt.Printf("test %d wrong number of strings\n", idx)
			os.Exit(1)
		}
		strs := fields[1:]
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(strs, "\n"))
		expected, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
