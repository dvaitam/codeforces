package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func isPalindrome(s string) bool {
	n := len(s)
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			return false
		}
	}
	return true
}

func isSubsequence(sub, s string) bool {
	j := 0
	for i := 0; i < len(s) && j < len(sub); i++ {
		if s[i] == sub[j] {
			j++
		}
	}
	return j == len(sub)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `ynbiqpm
plsgq
jey
tz
rwzte
dxcvk
rdlnktug
poqibzrac
wzvuatp
hxkwcg
hhzezrocck
pdjrjwdrk
gztrsjoct
kshjfgf
t
pccvy
ebc
vmwqiqzhg
nsiopvuwzl
kt
ps
ghaxid
lzfk
bdzewhb
urtvcadugt
dmcldbtagf
pg
v
r
tdichcu
lnfbq
btdwmgil
sfwvgybz
fkq
dtovf
p
sqjulmv
erwao
kx
r
ehypl
jvlsutewjm
ucatgwk
hhu
mwvsnbmw
nyvwbfociw
oqp
tyabpkjob
grucxea
vnkagaw
v
tdgdtugji
dpm
ai
zzdieuqu
deiabb
virk
sbxwtu
wuounlrf
msja
eik
zlwcky
bbifesjlmr
jdp
bjfq
jm
jnddrp
kzzkdpdw
nbjkxvef
mzucczcgxh
m
d
rqjopzs
nclh
syfng
dcwaqo
dpmi
ubzg
edgomlredt
esmuvnqp
ppuvgr
hakwxkkbqe
tzems
wwzpc
qb
he
j
y
kfzeuolq
qqbscvz
ytcxnygjr
npzmtshzav
x
jqs
kcpij
mzmbfue
jxkb
p
eptwcvw
zln`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := line + "\n"
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if len(out) > 10000 || !isPalindrome(out) || !isSubsequence(line, out) {
			fmt.Printf("test %d failed: invalid output\n", idx)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
