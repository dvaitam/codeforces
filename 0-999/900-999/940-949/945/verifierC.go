package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `cb bc
c c
aad daa
a a
bcddb bddcb
edc cde
abadd ddaba
eeebc cbeee
ecbb bbce
beaa aaeb
eeced decee
ebbe ebbe
cddcd dcddc
cd dc
cde edc
abc cba
cd dc
cba abc
addca acdda
eae eae
dc cd
bd db
eb be
c c
cbd dbc
bccce ecccb
dccea aeccd
e e
aecab bacea
dbbe ebbd
ac ca
bbb bbb
a a
ccab bacc
ada ada
e e
bb bb
bd db
ee ee
e e
aacce eccaa
baadd ddaab
e e
cedbc cbdec
acbcc ccbca
ad da
cc cc
beacd dcaeb
ccd dcc
ebe ebe
dc cd
dde edd
adaad daada
ed de
ab ba
cebbb bbbec
dcdb bdcd
acb bca
cdee eedc
adde edda
e e
bad dab
bbbd dbbb
c c
ec ce
b b
cce ecc
de ed
cdcda adcdc
bbea aebb
bea aeb
bbea aebb
d d
d d
aaa aaa
c c
dcb bcd
e e
dddaa aaddd
abcee eecba
deabd dbaed
bccce ecccb
cdcca accdc
dd dd
cdb bdc
cab bac
bddd dddb
cdcea aecdc
ee ee
a a
ed de
bdad dadb
ee ee
e e
bbaa aabb
abcdc cdcba
b b
c c
eacb bcae
ccadd ddacc`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	caseNum := 0
	passed := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Println("invalid test case:", line)
			continue
		}
		input := strings.Join(parts[:len(parts)-1], " ")
		expected := parts[len(parts)-1]
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input + "\n")
		out, err := cmd.CombinedOutput()
		result := strings.TrimSpace(string(out))
		caseNum++
		if err != nil {
			fmt.Printf("Case %d: runtime error: %v\n", caseNum, err)
			fmt.Printf("Output: %s\n", result)
			continue
		}
		if result == expected {
			passed++
		} else {
			fmt.Printf("Case %d failed: expected %s got %s\n", caseNum, expected, result)
		}
	}
	fmt.Printf("%d/%d cases passed\n", passed, caseNum)
}
