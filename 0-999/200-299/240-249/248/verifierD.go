package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Testcases embedded from testcasesD.txt (one per line).
const rawTestcases = `4 22 .HS.
8 28 .H.HSS.H
4 19 .HSS
4 24 H.S.
1 3 H
10 11 SHSS..S.SS
10 24 HSHHHSHS.S
5 18 .S.SH
7 25 HS.HS..
3 13 H..
4 24 .SSH
2 17 HS
6 8 SHHSSS
2 3 H.
7 25 S.S.HHS
1 3 H
10 27 HHSS.SH.HS
6 17 HSSS.S
10 27 H..SS...HS
7 15 .S.SHS.
6 6 S...HH
6 20 S.S..H
8 8 .H.HSS.S
5 23 .SHSH
6 17 .SSSHH
10 14 S.H.SHSH.S
2 5 .H
4 18 HHS.
4 22 SSHH
1 17 H
6 24 HSS.H.
6 24 HSS.SS
6 26 SSS.SH
7 11 HHS..S.
4 5 .SH.
9 19 HH.SHHHH.
9 15 S.HHS.HH.
5 12 .H..S
1 20 H
6 10 S.SHSH
4 7 .HHH
5 9 HS..S
1 9 H
5 24 ..SHS
6 6 HHHHHH
8 9 .H..SSHS
2 13 SH
10 19 SSHSSHH.H.
7 9 .HHSS..
9 21 .H.SHS.S.
6 19 .SSHHH
9 17 ..HSHSHHS
6 23 SHS.H.
9 21 .H.S..H..
1 16 H
4 16 H.H.
2 14 HH
6 9 HH.S.S
10 19 HH....HH..
2 19 H.
6 24 HHHH.H
8 27 .SSS.SH.
2 14 .H
7 12 S..H...
3 23 SHH
2 17 .H
10 15 HSHH..SH..
5 18 ...SH
5 5 SSSHH
10 21 HSSHS.S..H
8 26 ..HS.HS.
1 15 H
4 24 ..HH
5 12 HSHH.
1 9 H
1 11 H
7 9 .HHHSSH
6 20 ...SHH
6 16 SSSHH.
10 25 SH.SHSH.SH
8 24 SH..S.SS
5 25 SHS..
9 25 H.S.SH.S.
10 15 ....HHS.SH
2 19 HS
10 23 .S.HSSSH.H
2 7 .H
6 9 SSSHHH
8 24 SH..SSSS
9 11 SSHHS.H.H
10 13 H.H.S..S.H
10 29 HSH..H.HHS
6 15 HSSSS.
10 19 .S...S.SH.
4 24 HH.H
8 13 ..SHH.HH
9 29 SH.S.HHSH
10 11 HS.H.HS.HS
6 16 S..HH.
1 3 H`

func parseTestcases() []string {
	lines := strings.Split(strings.TrimSpace(rawTestcases), "\n")
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}

func solveCase(n int, t int64, s string) string {
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		val := 0
		switch s[i-1] {
		case 'S':
			val = 1
		case 'H':
			val = -1
		}
		a[i] = a[i-1] + val
	}

	check := func(x int) bool {
		if int64(a[n])+int64(x) < 0 {
			return false
		}
		to := 0
		for i := 1; i <= n; i++ {
			if s[i-1] == 'H' || int64(a[i-1])+int64(x) < 0 {
				to = i
			}
		}
		var res int64
		for i := 0; i <= n; i++ {
			dist := to - i
			if dist < 0 {
				dist = 0
			}
			if res+int64(dist)*2 <= t {
				return true
			}
			if i == n {
				break
			}
			cost := int64(1)
			if int64(a[i])+int64(x) < 0 {
				cost += 2
			}
			res += cost
		}
		return false
	}

	l, r := 0, n+1
	for l < r {
		mid := (l + r) >> 1
		if check(mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	if l <= n {
		return strconv.Itoa(l)
	}
	return "-1"
}

func parseCase(line string) (int, int64, string, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) != 3 {
		return 0, 0, "", fmt.Errorf("invalid line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, "", err
	}
	t, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return 0, 0, "", err
	}
	s := fields[2]
	if len(s) != n {
		return 0, 0, "", fmt.Errorf("expected string length %d got %d", n, len(s))
	}
	return n, t, s, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	for idx, line := range testcases {
		n, t, s, err := parseCase(line)
		if err != nil {
			fmt.Printf("case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected := solveCase(n, t, s)
		input := fmt.Sprintf("%d %d\n%s\n", n, t, s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
