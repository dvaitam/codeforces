package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesRaw = `i i
ragco ocgar
twqpyhixub buxihypqwt
cdhtkp pkthdc
ueb beu
zqpryziau uaizyrpqz
jjiku ukijj
gvmummnbng gnbnmmumvg
agaopzpnx xnpzpoaga
dlkqf fqkld
awfbsu usbfwa
zgxtalwloq qolwlatxgz
dfhlhn nhlhfd
wy yw
hjlqdy ydqljh
i i
eyfcv vcfye
onwvoajbo objaovwno
dktim mitkd
sqmea aemqs
tzigv vgizt
ypmqwesmn nmsewqmpy
zxjnymd dmynjxz
vbnbzkzoqb bqozkzbnbv
oddx xddo
m m
ddccxherp prehxccdd
xnc cnx
ioyrtzl lztryoi
gtqou uoqtg
fiwcqdup pudqcwif
bg gb
qgf fgq
qqyoyvbw wbvyoyqq
njpj jpjn
nggkehwwsx xswwhekggn
l l
imckz zkcmi
jg gj
qcts stcq
zcovegzky ykzgevocz
wys syw
na an
nifnpytg gtypnfin
p p
zchlh hlhcz
ntz ztn
mrek kerm
txppahtos sothappxt
wrlxmauwew wewuamxlrw
hsuhawrrs srrwahush
woalsxbeyv vyebxslaow
ehgswzp pzwsghe
nonfdhkxkz zkxkhdfnon
lne enl
vqkkv vkkqv
kll llk
qvolgbe ebglovq
nfknvcq qcvnkfn
q q
bbin nibb
kufgcvc cvcgfuk
wo ow
vqtwwbcfj jfcbwwtqv
ggpb bpgg
oazkw wkzao
rcvfnu unfvcr
vakrjcx xcjrkav
yhseikbvgs sgvbkieshy
weib biew
skqhpbs sbphqks
nzxmhuzztz ztzzuhmxzn
jzxlaoubmo ombuoalxzj
i i
ujqp pqju
iqmlmsvu uvsmlmqi
mxjigqhpuf fuphqgijxm
ce ec
xsapi ipasx
xgffdz zdffgx
d d
uxb bxu
ec ce
jqnvd dvnqj
lt tl
wuzsuunz znuuszuw
yd dy
jy yj
uhbnx xnbhu
tys syt
mi im
ugyafltxlg glxtlfaygu
qtpuktgfg gfgtkuptq
bn nb
gznkq qknzg
zvimvx xvmivz
k k
kxfdxmsqh hqsmxdfxk
alylr rlyla
inrqcxug guxcqrni
`

// solve mirrors 1089D.go logic (string reverse).
func solve(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

type testCase struct {
	input string
}

func parseTestcases() []testCase {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		cases = append(cases, testCase{input: parts[0]})
	}
	return cases
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	got, err := runBinary(bin, tc.input)
	if err != nil {
		return err
	}
	expected := solve(tc.input)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := parseTestcases()
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
