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

const testcasesBRaw = `4 4 abcb bbbb caca baac
5 2 cb bc ba ca ac
4 1 c b b a
3 1 a c a
2 2 cb aa
3 5 babcb cacbc acccb
3 4 aaca aaac ccbc
5 3 caa ccb cbb bcc cba
3 5 abccb aaacb acaba
3 4 aaac aacc ccca
1 1 c
2 5 cabab aacaa
2 1 b a
1 1 c
4 5 abaaa cbbba acbac acbab
1 5 bcaab
4 5 cbcbb cbacc abcab cacba
2 3 bcb caa
5 4 abcc bcca abab cabc bcca
1 4 bbba
4 2 cc ac aa bc
2 3 bca abb
1 1 b
4 1 b a c c
3 1 a b a
1 1 a
3 5 bbcac cccbc cbcbb
5 2 ab cb aa ab bb
3 1 b c a
1 3 aac
3 3 bca bab caa
3 2 cc ab bb
3 4 aacb bbba bacb
2 2 cc bc
1 1 a
2 2 ab aa
4 5 cbbbc ccaba babca babaa
1 5 bcaab
4 3 aca aca aab bbc
2 1 c b
2 5 caabc cbbac
2 3 cab cbc
4 1 a c a a
3 2 ac bb cb
5 5 bbabb bbbba acaab cbaab babcc
4 2 aa cc ac ba
1 1 a
4 3 acc cca aac ccc
4 5 baabc baaba abaca bbcba
4 3 cba baa bca bba
1 2 ba
5 2 aa ac cc ac ab
3 1 c a a
2 4 abbc baaa
1 5 cbaba
2 2 bc aa
1 4 acaa
4 3 cbc acb bcb aab
3 5 abbab cbbab bbaaa
4 5 baaba cbcca bcbac bacbb
1 1 b
1 1 b
4 5 ccbba cbbbc acaaa bbcba
2 5 abccb cbbbb
5 2 aa ca bb cb ac
3 3 cbc bab bcc
3 4 cabb acba caab
4 4 acac accc baba acab
3 2 ac ba ab
4 4 caab acca baac aaba
2 2 cb ca
5 3 bbb acc ccc ccb cbb
3 4 abba aaab abca
2 1 b c
3 3 cbb bac aab
1 2 cb
2 2 bc cc
3 3 cbb bbc aab
5 1 a b b c b
1 4 bbba
1 1 a
1 2 ba
3 4 ccbc cbba aacc
1 2 cb
2 2 bc ca
5 5 baabc baccb acccc cbaca acbcb
1 2 ba
2 5 aaabc caaca
2 4 bcca cbac
1 5 cbbbb
2 4 cbcb caba
3 2 ba bc ca
2 3 abc aaa
1 1 a
3 1 a b c
4 2 ab ba bb ac
3 4 babb cccc cbac
1 4 baac
2 5 cccca abaaa
5 2 ca bc ac cb ca
`

func expected(n, m int, grid []string) string {
	row := make([][26]int, n)
	col := make([][26]int, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			c := grid[i][j] - 'a'
			row[i][c]++
			col[j][c]++
		}
	}
	var ans []byte
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			c := grid[i][j] - 'a'
			if row[i][c] == 1 && col[j][c] == 1 {
				ans = append(ans, grid[i][j])
			}
		}
	}
	return string(ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+n {
			fmt.Printf("test %d: wrong number of rows\n", idx)
			os.Exit(1)
		}
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			if len(parts[2+i]) != m {
				fmt.Printf("test %d: row %d has wrong length\n", idx, i+1)
				os.Exit(1)
			}
			grid[i] = parts[2+i]
		}
		expect := expected(n, m, grid)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			input.WriteString(grid[i])
			input.WriteByte('\n')
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed:\nexpected: %q\n   got: %q\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
