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

// Embedded source for the reference solution (was 1243B2.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       var sStr, tStr string
       fmt.Fscan(reader, &sStr)
       fmt.Fscan(reader, &tStr)
       s := []byte(sStr)
       t := []byte(tStr)

       // check total character counts
       cnt := make([]int, 26)
       for i := 0; i < n; i++ {
           cnt[s[i]-'a']++
           cnt[t[i]-'a']++
       }
       possible := true
       for i := 0; i < 26; i++ {
           if cnt[i]%2 != 0 {
               possible = false
               break
           }
       }
       if !possible {
           fmt.Fprintln(writer, "No")
           continue
       }
       fmt.Fprintln(writer, "Yes")
       // perform swaps
       type pair struct{ x, y int }
       ops := make([]pair, 0, 2*n)
       for i := 0; i < n; i++ {
           if s[i] != t[i] {
               // try to find in s
               found := false
               for j := i + 1; j < n; j++ {
                   if s[j] == s[i] {
                       // swap s[j] with t[i]
                       ops = append(ops, pair{j, i})
                       s[j], t[i] = t[i], s[j]
                       found = true
                       break
                   }
               }
               if found {
                   continue
               }
               // find in t
               for j := i + 1; j < n; j++ {
                   if t[j] == s[i] {
                       // swap s[j] with t[j]
                       ops = append(ops, pair{j, j})
                       s[j], t[j] = t[j], s[j]
                       // swap s[j] with t[i]
                       ops = append(ops, pair{j, i})
                       s[j], t[i] = t[i], s[j]
                       break
                   }
               }
           }
       }
       // output operations
       fmt.Fprintln(writer, len(ops))
       for _, op := range ops {
           // convert to 1-based
           fmt.Fprintln(writer, op.x+1, op.y+1)
       }
   }
}
`

const testcasesRaw = `100
2
cc
cf
6
itgtbs
itgtbf
8
umzxqlro
umzxblro
2
lo
lm
8
qfrfhhaf
qfrfhfaf
4
qqlq
qolq
8
xqylzsll
xqylzslf
8
wxouqhpi
wxouqhpq
10
zlvoolsxrx
zlvoolsprx
5
kwfti
kwfji
6
zwqrqq
zwqrqt
8
jxgpqlvt
jzgpqlvt
7
xagxdbs
xagxdis
5
vdyqe
vdhqe
5
bnwyb
lnwyb
7
fhvacdc
bhvacdc
2
li
fi
4
qwam
zwam
5
ebalt
jbalt
7
pajoryt
pajorbt
6
ymtwep
yctwep
7
daozeqs
daozeqm
9
qkekiitnu
wkekiitnu
10
evbibeffdo
evbqbeffdo
2
hh
hc
6
cshtzz
cshtwz
7
ivniqya
ibniqya
8
nfdqxchd
nadqxchd
4
yhdg
qhdg
9
ojrumgvyg
ojrumgnyg
10
assbnqsfdv
assbnqsldv
2
qd
qj
7
javnddj
jyvnddj
2
zo
no
9
ogstcajal
ogstxajal
3
hyp
dyp
7
mwoeylm
iwoeylm
3
dct
dut
8
gwdatvpy
xwdatvpy
9
jloezlipq
jloezlixq
8
pvjmhfpt
pvjmrfpt
8
wvwcsxsd
wlwcsxsd
4
rezn
zezn
3
vzu
ezu
6
mhwvvk
mhwfvk
10
jderyndkqh
jderyndkih
4
fowh
fowl
4
ooxa
ooxx
4
mqbp
mqmp
6
wxnwup
wxrwup
7
wxvcyxh
wxvctxh
5
mvmua
mvoua
10
woufdamgxs
woufdamgxm
5
dmryz
diryz
5
pztea
pztev
8
piqsfowg
plqsfowg
2
pr
yr
9
vkoiqoact
vkoiqfact
8
ivuzxebf
ivuzxebm
9
vjeajroal
rjeajroal
8
sogvjpue
sogvjpuw
10
wjcikjkuzj
wjcikjquzj
3
qug
qtg
10
ezqucjbhor
ezqqcjbhor
6
bddvzm
bdgvzm
7
lckolfp
lckjlfp
9
ewougikfd
ewopgikfd
5
yvlfl
yzlfl
4
hizr
hizm
7
ixtqswx
ixxqswx
8
ywwuywjr
ylwuywjr
6
mpfilo
mpfclo
4
kmea
lmea
4
lcxy
lcxa
10
khtmrjpuel
khtmrgpuel
9
dezgkienl
dezgcienl
7
ghwhxtb
ghlhxtb
2
ef
nf
9
yiekqsdku
yiekqshku
2
my
mp
7
rttcsqr
rttcspr
8
wofnmqob
wsfnmqob
4
dvqf
mvqf
6
ozwaid
ozwail
5
faenv
kaenv
9
bphcperae
bphcperar
2
bg
qg
7
vqhelpa
vrhelpa
3
hdo
zdo
2
tg
tk
8
wqqyvfqd
wquyvfqd
5
fmgjk
fmgek
8
emkzjzdr
epkzjzdr
6
jqypih
jqywih
4
wrvd
trvd
10
ygggmsbueu
xgggmsbueu
6
wxprbx
weprbx
7
bwgdeuw
bwgdfuw
3
vou
vgu
4
kwiq
nwiq`

var _ = solutionSource

type testCase struct {
	n int
	s string
	t string
}

func parseTests() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return nil, fmt.Errorf("invalid test data")
	}
	t, _ := strconv.Atoi(scanner.Text())
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected EOF at case %d", i+1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing s at case %d", i+1)
		}
		s := scanner.Text()
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing t at case %d", i+1)
		}
		tt := scanner.Text()
		if len(s) != n || len(tt) != n {
			return nil, fmt.Errorf("length mismatch at case %d", i+1)
		}
		cases = append(cases, testCase{n: n, s: s, t: tt})
	}
	return cases, nil
}

func possible(tc testCase) bool {
	cnt := make([]int, 26)
	for i := 0; i < tc.n; i++ {
		cnt[tc.s[i]-'a']++
		cnt[tc.t[i]-'a']++
	}
	for _, c := range cnt {
		if c%2 != 0 {
			return false
		}
	}
	return true
}

func validate(tc testCase, output string) error {
	out := strings.Fields(output)
	if len(out) == 0 {
		return fmt.Errorf("empty output")
	}
	poss := possible(tc)
	if strings.ToLower(out[0]) == "no" {
		if poss {
			return fmt.Errorf("claimed impossible but transformation exists")
		}
		return nil
	}
	if strings.ToLower(out[0]) != "yes" {
		return fmt.Errorf("first token not Yes/No")
	}
	if !poss {
		return fmt.Errorf("claimed possible but counts mismatch")
	}
	if len(out) < 2 {
		return fmt.Errorf("missing operation count")
	}
	k, err := strconv.Atoi(out[1])
	if err != nil {
		return fmt.Errorf("invalid operation count")
	}
	if len(out) != 2+2*k {
		return fmt.Errorf("expected %d indices got %d", 2*k, len(out)-2)
	}
	if k > 2*tc.n+2 {
		return fmt.Errorf("too many operations")
	}
	s := []byte(tc.s)
	t := []byte(tc.t)
	for i := 0; i < k; i++ {
		u, err1 := strconv.Atoi(out[2+2*i])
		v, err2 := strconv.Atoi(out[2+2*i+1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid op value")
		}
		if u < 1 || u > tc.n || v < 1 || v > tc.n {
			return fmt.Errorf("op out of range")
		}
		u--
		v--
		s[u], t[v] = t[v], s[u]
	}
	if string(s) != string(t) {
		return fmt.Errorf("strings not equal after ops")
	}
	return nil
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
		fmt.Println("Usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		sb.WriteString(tc.t)
		sb.WriteByte('\n')
		input := sb.String()
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := validate(tc, got); err != nil {
			fmt.Printf("case %d failed: %v\noutput:\n%s\n", i+1, err, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
