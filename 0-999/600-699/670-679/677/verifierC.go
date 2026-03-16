package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "677C.go")
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesCRaw = `lk
vNGBeu3YV4Ie
U
OW2vwEDdwPwrUx5
UTU5uZ7F-J_T67
6-CPvI9NM0NA
UjRbynhgIDnrIFAh
ehUUwFdkoidfcV
quxaXfFte
S
oKR-dN5fHZt8ClOnd5qY
PsRHH1crhGequvm6
eFD4jGkD
UG2JateX0uolEnmcxDnB
7
NWBA32cg1xm9Ucp
ULVNc0mnNzc5h0-7AjaK
V
jC-yoVY7rS
pGppkQYBnd8f_
T6sVI991-L
Du-H2kmjTws1i
eqL
DQ4wKot2mQFGv
7EZTs7
dXxYg9JZG08UQkC
yZWbO77wmcZBXBmXzJ
y-rb39Gw7AjSa-i-Q6I
6dkSwZGrgu_W7LtbK
7aUeW4AN_r9MjHOMQN
lAYtlNfD6DJho
WUBO
jQ6Uv_4L7r4B
PumE8yVxT
rDIWZ
JPZLjVNY9wH
49lxOWqdnSvT
3bP
XK8tUOz_
sAQG
1UGlR
FEfRVhs
i34IqP
oRYDgY8-Ol_Z6v0X5
n5
qpwjYN6bGnSCwds2lR7
8E
9rd
hgzaR
ErV-aqoFn7BgBWRYu
ntAwWzMR3s2qYOMmm
IK-JD1rndzByYfrd
8gCsOeznr
xl7LAuPJi01e6MpIcB
QHYz3qv56SX8Gy
84y8RNjvV8CqNAMmbdz
hOGR4j18cKq
tuWi5Jk
9EtMDzQY0EBhHFrY
p6YY8WKBEChlag
3ZDImUU_i6CJd
9
qs
PEgsLnl
3reNI
gTRmUnSUJ9KtdfR3
S
gj321ExufcTxLcfFCZ
Uoi
DynaZkJC
Z2
t2q6V
x4
5u_qSscGxt0G48
y23ICTeYd2Md9HI
76U7Fu6K
1oFXo376lX5U
SvsDw05_w0HOZMHObZ
fA6mobUPOZwPk9ZE4mc
dM_rgbQZ9a9A
PvQMY-7J
BUE
VxECA7EYIAua
8Uxzx-arABak6
xIYcbUo
eU8Wpj7ws6
Is9
jLKdAj1qxN7zfR6htD
SNnytCdPpInVjrSq2zNc
R
jgr_1
9LMitoalsTtM
qlDDQlx3A0eQ2_Sxpfq
s1CIplwzjOqEE
Z6nODgSnYcoHTER
DfrvqvYtnM27ao
amey3ayWOtwIIG
DNN1_9XLYw1L62PjbG0
6l8JuBVnCB3qyXujN6s
lVxCf3eYMDuZCzDGJ`

	scanner := bufio.NewScanner(strings.NewReader(testcasesCRaw))
	idx := 0
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		idx++
		input := s + "\n"
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
