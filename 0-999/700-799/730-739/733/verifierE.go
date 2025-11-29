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

const testcasesE = `
100
6
DDUUDU
8
UUUUDDDU
7
DDUUUDD
7
DUDDUUU
1
D
2
DU
4
DUUD
7
UUDUUUD
7
UUDDUUU
8
DUDDUDDU
1
D
1
U
3
DUD
9
UUDDUDDUU
2
UU
6
DDDUDD
9
DDDUDDUUU
4
DUDD
6
DDUUDD
3
UDD
6
DDDUUD
6
UDUUDU
10
UDUUDUDDUD
5
UDDDD
7
DDDUDDD
6
UUUDUU
3
UDD
4
UUDD
5
UUUUD
2
UD
3
DUD
9
UDUDDDUDD
9
DDUUDDUUD
3
DUU
7
UUDUDDU
6
DDUUDU
1
D
3
UUU
7
UDUUUUD
2
DU
10
UDDDUDDDUU
1
D
1
U
2
DD
3
UDD
2
UU
4
UUDD
7
UUDUUUD
1
D
3
UUU
3
UUD
2
DU
10
UDDUUUDDUD
6
UDUUDU
1
D
1
D
4
DUUU
8
UUUUUDDU
3
UDD
6
DDDUUU
5
UUDDU
8
UUUUDUDD
10
UUDUDDUDUU
10
DDDDUDUDDU
9
DUDUDUUDU
3
DDU
2
UU
4
DDDD
4
UUUD
8
UUDDDDDD
6
UUDDDU
9
DUUDDUUUU
5
DDUDU
6
DDUUDD
3
UDD
2
UU
7
UUDUUUU
5
DDDDD
7
UDDDDUU
5
DUUDU
7
DUDDDUU
6
DDUUDU
10
DUDUUDUUUU
3
UUU
6
UUUUDD
6
UDDDDU
5
DUUDU
10
DDUDDUDUDU
5
DUDDD
9
DDUDDUUUU
9
UDUDDDUDD
7
UDUDUDU
3
DUD
4
DUUU
7
UDDDDDD
9
DDDUUUDDD
7
UUDUDDU
2
UU
2
DU
3
DUU
`

func simulate(orig []byte, start int) int {
	n := len(orig)
	s := make([]byte, n)
	copy(s, orig)
	pos := start
	steps := 0
	limit := n*n*2 + 5
	for steps <= limit {
		if pos < 1 || pos > n {
			return steps
		}
		if s[pos-1] == 'U' {
			s[pos-1] = 'D'
			pos++
		} else {
			s[pos-1] = 'U'
			pos--
		}
		steps++
	}
	return -1
}

func expectedCase(n int, str string) string {
	arr := []byte(str)
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(simulate(arr, i)))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	scan := bufio.NewScanner(strings.NewReader(testcasesE))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		s := scan.Text()
		input := fmt.Sprintf("%d\n%s\n", n, s)
		exp := expectedCase(n, s)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
