package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB = `R140892C596854
FZWQ267460
ZUEX797927
R683245C398056
R98419C511555
AWVJH936711
R453790C636945
ZNVK729634
R756590C840776
R619870C991189
ZSAH422724
R858483C129907
R534610C931184
BLGS966392
R99116C115277
R922465C959347
R775066C696770
K718707
WQQX309181
OARL602500
R757005C61126
R37817C825816
R134378C152082
R556707C480438
R774593C95471
AF941901
R380727C787845
R765815C998211
R117199C36313
R783876C571279
R335791C547171
R20268C567806
ZUW796509
R761482C882267
R912498C194187
R206691C66147
R429008C37978
R526945C95344
R891467C519484
R41924C174012
R432372C286135
WDTC753892
R608052C631651
GRVE23233
R330381C1836
R62683C182990
R812602C930708
R897155C490525
HZPB285426
R260536C434071
R916734C801549
R621502C835907
R320057C3539
UX570394
R482855C523546
R976962C4761
R951360C175650
R389220C66342
R209806C21254
R283722C160722
R292156C23012
LJSU536972
R344079C929548
R727163C526057
R911653C954084
R302495C437580
R421156C545777
XQRP334029
R6666C70842
R455516C97458
R218615C51990
R473011C982892
R506785C117205
SWXE699453
R229216C62906
R174817C383219
R856211C61229
CC34410
R843647C507604
PHNJ155800
R802309C614494
R961153C368487
R99304C393447
R945724C101337
R738529C901657
R585598C852955
R991787C551778
R54017C91212
R301969C900814
R657805C761027
R806350C377822
BCHE257128
R230993C2180
R644996C574749
R635333C417719
XPYY254467
R483617C250643
R837149C376301
R181223C985708
R599644C698796
R534215C329015
R49149C448227
R280897C464526
R209034C368020`

func colToLetters(c int) string {
	s := ""
	for c > 0 {
		c--
		s = string(rune('A'+(c%26))) + s
		c /= 26
	}
	return s
}

func lettersToCol(s string) int {
	c := 0
	for _, ch := range s {
		c = c*26 + int(ch-'A'+1)
	}
	return c
}

func isRXCY(s string) bool {
	if len(s) < 4 || s[0] != 'R' {
		return false
	}
	idx := strings.IndexRune(s, 'C')
	if idx == -1 {
		return false
	}
	if idx == 1 || idx == len(s)-1 {
		return false
	}
	for i := 1; i < idx; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	for i := idx + 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func convert(s string) string {
	if isRXCY(s) {
		idx := strings.IndexRune(s, 'C')
		row := s[1:idx]
		col := s[idx+1:]
		c, _ := strconv.Atoi(col)
		return colToLetters(c) + row
	}
	i := 0
	for i < len(s) && s[i] >= 'A' && s[i] <= 'Z' {
		i++
	}
	col := lettersToCol(s[:i])
	row := s[i:]
	return fmt.Sprintf("R%vC%v", row, col)
}

type testCase struct {
	s string
}

func parseTests() []testCase {
	lines := strings.Split(strings.TrimSpace(testcasesB), "\n")
	tests := make([]testCase, len(lines))
	for i, line := range lines {
		tests[i] = testCase{s: strings.TrimSpace(line)}
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	tests := parseTests()
	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != len(tests) {
		fmt.Printf("expected %d outputs, got %d\n", len(tests), len(lines))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := convert(tc.s)
		if strings.TrimSpace(lines[i]) != want {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, want, lines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
