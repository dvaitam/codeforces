package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const solution1336CSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var S, T string
	if _, err := fmt.Fscan(in, &S); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	n := len(S)
	m := len(T)
	s := []byte(S)
	t := []byte(T)

	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, n)
	}

	if m <= n {
		for i := 0; i < n; i++ {
			if i >= m || t[i] == s[0] {
				dp[i][i] = 2
			}
		}

		for length := 1; length < n; length++ {
			c := s[length]
			for l := 0; l+length-1 < n; l++ {
				r := l + length - 1
				val := dp[l][r]
				if val == 0 {
					continue
				}
				if l > 0 {
					if l-1 >= m || t[l-1] == c {
						dp[l-1][r] = (dp[l-1][r] + val) % MOD
					}
				}
				if r+1 < n {
					if r+1 >= m || t[r+1] == c {
						dp[l][r+1] = (dp[l][r+1] + val) % MOD
					}
				}
			}
		}

		ans := 0
		for i := m - 1; i < n; i++ {
			ans = (ans + dp[0][i]) % MOD
		}
		fmt.Println(ans)
	}
}
`

// Keep the embedded reference solution reachable.
var _ = solution1336CSource

const MOD int = 998244353

type testCase struct {
	s string
	t string
}

const testcasesRaw = `dvee astu
dwyuwg tutqd
xtf mcz
d l
kdvols wivp
rr c
aufqn tv
wv mi
m f
mzcu s
hqob ptd
pru mpi
rl f
eoc cpm
rciph djel
v b
srgabm rw
wwdpr lykd
hhpjih a
qkcc jsn
kbsc hm
wlsrl gqdr
zxo xmikm
vcs rnvbz
ki i
numzpn pssxt
wxfoms x
ee kqgz
ch yhg
uf y
afqoev bpw
bnebsb ryb
fwtodg det
uznpy fu
oseg ufy
kka fu
gvmjzq ct
fssnahh f
pl zg
c r
cbnir tn
rc qsk
lzfwk jc
dijkyr ugjv
f qitah
xe c
qr hjiq
h iyqy
ijqge hh
iex ekmlq
bhlcyzg a
vls ytasl
mqkoha xe
smzvhf rsgsh
tocygn iadkiy
vo bd
exuxz lcqs
oi of
rbftauel lrzetv
qamlc zcps
epivh xntw
qnj a
bt uawuxs
v buphy
xq ry
rybirbw x
wz zdxhq
dfvv j
zclnk hmo
nlygtyi igu
gu rhb
tntkj jmcbq
srvjaa xwbk
wdfy p
bmeyp lgkru
wj t
nklguh snrufac
vwe mpkmi
jp zilrc
nofr l
jr cc
qnc wexlikx
chs jorkcuy
s fxoc
x v
zppj sjvq
q pr
s lpzm
vj nyqyh
lmvaf kdaky
lsv hcknxm
n pltc
pxypbc imq
bb o
yn nf
vvihz modgggb
j hro
`

func parseTestcases() []testCase {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	res := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		res = append(res, testCase{s: parts[0], t: parts[1]})
	}
	return res
}

func expected(S, T string) int {
	n := len(S)
	m := len(T)
	s := []byte(S)
	t := []byte(T)
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, n)
	}
	if m <= n {
		for i := 0; i < n; i++ {
			if i >= m || t[i] == s[0] {
				dp[i][i] = 2
			}
		}
		for length := 1; length < n; length++ {
			c := s[length]
			for l := 0; l+length-1 < n; l++ {
				r := l + length - 1
				val := dp[l][r]
				if val == 0 {
					continue
				}
				if l > 0 && (l-1 >= m || t[l-1] == c) {
					dp[l-1][r] = (dp[l-1][r] + val) % MOD
				}
				if r+1 < n && (r+1 >= m || t[r+1] == c) {
					dp[l][r+1] = (dp[l][r+1] + val) % MOD
				}
			}
		}
		ans := 0
		for i := m - 1; i < n; i++ {
			ans = (ans + dp[0][i]) % MOD
		}
		return ans
	}
	return 0
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%s\n%s\n", tc.s, tc.t)
}

func runCase(bin string, idx int, tc testCase) error {
	expect := expected(tc.s, tc.t)
	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("case %d failed: %v\nstderr: %s\ninput:\n%s", idx, err, string(out), input)
	}
	gotStr := strings.TrimSpace(string(out))
	if gotStr == "" {
		gotStr = "0"
	}
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("case %d failed: invalid output %q\ninput:\n%s", idx, gotStr, input)
	}
	if got != expect {
		return fmt.Errorf("case %d failed: expected %d got %d\ninput:\n%s", idx, expect, got, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	for i, tc := range testcases {
		if err := runCase(bin, i+1, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
