package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)
const testcasesARaw = `
7 14
aeihgmn
8 11
fjdicebj
9 10
jcebfhibf
6 15
jkodih
9 16
ibacmapkh
6 25
wcgshe
8 19
ckqpdjrd
6 19
rgsjoc
7 21
kshjfgb
5 21
pcebc
9 24
vmwqihgvs
10 15
ehklomfbjh
6 20
ghaidh
3 13
fga
3 5
bae
9 22
tvcadugts
7 5
aceabda
1 23
r
10 15
bnebdkefgc
9 3
bacbacbac
4 23
bvfk
5 18
dofap
7 23
sqjulmv
3 10
iah
2 25
kx
9 3
bacbacbac
6 8
cdhgac
2 16
if
9 16
pabpkjobn
9 8
bcagfadbe
3 10
bhg
2 22
ai
2 16
ie
9 22
uldeiabgv
9 10
fjahgficd
10 14
eacefmlbfm
1 21
b
3 10
cje
7 13
icebhld
5 3
acbac
8 17
kdpnbjke
10 7
dgfagbfade
5 18
opgnc
4 13
ejmc
4 15
fbnl
2 2
ba
8 5
dcbaebad
6 14
incbjh
10 6
dfedfcdfbe
4 21
akbq
5 6
ebdec
8 24
cqbchebj
2 2
ba
8 6
cedaefad
4 26
jrtn
7 17
hafjqik
8 4
cdabcadb
10 17
celnbombdp
3 26
abt
3 21
ukd
9 24
ulgmpdbtw
10 16
kdjemjdgbm
6 16
golcbp
1 10
i
10 23
sghcuqwnqj
3 5
dea
7 5
adbadca
2 13
bf
1 24
l
3 13
adf
10 4
bacbdacbac
8 12
eikfcjbi
5 20
fmehk
4 18
hfjl
1 15
n
10 6
adfabdcedf
10 6
dcfabdfcea
7 14
aglfhdf
8 11
bcebijch
3 14
mgc
8 9
ficfhbde
2 2
ba
5 8
becgh
8 4
bdcadcba
2 3
bc
2 12
id
3 26
cnj
9 11
cjikdibgk
7 19
ijolsef
2 5
de
3 16
jlp
7 25
gpwqkpu
8 3
bacbacba
2 18
cm
2 2
ba
10 5
acbecbadcb
7 12
kgchlci
3 12
dch
7 13
gmhgldm
8 8
dagdbcfa
`


func solveA(n, p int, s string) string {
	a := []byte(s)
	for i := n - 1; i >= 0; i-- {
		for c := a[i] + 1; c < byte('a'+p); c++ {
			if i >= 1 && a[i-1] == c {
				continue
			}
			if i >= 2 && a[i-2] == c {
				continue
			}
			a[i] = c
			ok := true
			for j := i + 1; j < n; j++ {
				found := false
				for d := byte('a'); d < byte('a'+p); d++ {
					if j >= 1 && a[j-1] == d {
						continue
					}
					if j >= 2 && a[j-2] == d {
						continue
					}
					a[j] = d
					found = true
					break
				}
				if !found {
					ok = false
					break
				}
			}
			if ok {
				return string(a)
			}
		}
	}
	return "NO"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
	idx := 0
	for {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var n, p int
		fmt.Sscan(line, &n, &p)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing string on case %d\n", idx+1)
			os.Exit(1)
		}
		s := strings.TrimSpace(scanner.Text())
		idx++
		expected := solveA(n, p, s)
		input := fmt.Sprintf("%d %d\n%s\n", n, p, s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
