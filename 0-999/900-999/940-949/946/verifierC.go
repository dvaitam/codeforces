package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `100
cclf
zvjitgtbsvfnumzxqlroqibalokmnqfrfhhafkfeqqlqvrfo
xqylzsllofymwxouqhpipqqzlvo
lsxrxopvhkwftiypjjzwqrqqutsnjx
pqlvtczkxagxdb
ubishvdyqeihgbnwybbllfhvacdcabxaliefx
qwamsbzhebal
uxxdjkpajorytxbiymtwephcvvkdaozeqsympqke
iitnuawrevbibeffdouhqw
hhw
cicshtzztwlivniqyaebmnfdqxchd
afyhdga
voojrumgvygxznnqassbnqsfdvzplaqdtl
wljavnddjgyvazobnup
gstcajaljxchypgdslmwoeylmdiddc
kumgwdatvpybxwpjloezlipqpxxznpvjmhfptirn
vwcsxsdclfreznczcvzubejmhwvvkofqjderyndkqhwqi
fowhmlzysxe
oxaztmxfmqbpimiwxnwuplrkwxvcyx
rtgmvmuakoqwouf
amgxstm
dmryzgixssgpzt
atvnpiqsf
wgyclaprvvcyspvkoiqoactylfyyzm
vuzxebfpmovjeajro
l
rms
gvjpuepwrwjcikjkuzjuumqcqugmt
ezqucjbhorhqibddvzmlgklckolfpojoew
ugikfdhpgyvlflezehizrummzxkix
qswxkxmywwuywjrtuvcljmpfilopcfkmeadlflc
yunarkhtmrjpuelkgpdezgkienlickghwhxtbkluytbefcn
yiekqsdkuywtmhbmypptkrttcsqrv
mwofnmqobdosedvqfcmjozwaidvlhfae
vckuobphcperaewqrbbgraqkvqhe
paerdhdogzbtgumktumwqqyv
qdeugfmgjkn
nemkzjzdrd
ijqypihnwewrvdatrygggmsbueuaxiw
prbxyhetkbwgdeuwrfycvoujgfkwiqscnnvxbojvduwxiag
kirxmsqxgnyewfzoolmptitgspo
ypskjcfltuphy
uvsevjgrjdazagkbkrizxvkocnpwajssegeftymx
usoiu
ppzzhe
jhgtwkstwmqnhugrbivhetxmndommpmjghhbr
ctrvabmwnmhqidlqlqzpscwoxwhiaapbeu
gkhrbteujy
urrcvve
webjqvipbrlykvdtldtzllzuizpj
teabknualvrwbvcwrqtynnnhzfftbasyx
vfjabhszhmcldtchhrgdawm
qisuhb
qqmnzeeneoxlbsfqontzuofpteleaiwfeu
suiopognniyhlyubmtanjarpsiv
howolqtovhrrfojyln
qvhxuvmd
totqocxmotywwxlzrlfehvufnopw
zfniksmjuwzivkvamsbgoddalktkxmfxkzcqtpmthodtala
pewwbakmputapusghxtk
kjymstxpoyi
qgslhl
zfhrzuxthgsohmisgqfamypx
ufvgxwtfptazxegwgatcoyug
fimtaaldjbslpmdcofezyoucieprczqwzjjarrgcnefsjogb
kozbehtlujzwudgehakdidlycqzwtelenxxgjryayakue
crybexpnvlpvjjcuexdacelwejysvrechsxyxhzkcfvngnu
bknpurlfqdvbevmenyyxhuitrwdcqsxtfgvtvwcrkeq
hvovmodkhbldmasd
zilrhknhbefexfmedujynoadnawdbgnagzmkefiizishqjjs
ppmyjmfuznjzsoqnswkrcsainso
cqrpquifgldxhgneqgqmfcjoeyqcxlvfhzrbnrbmryjhfxm
gshqiritifautmq
bbzpgrkjxubokfobglixoivqibxkanuhfmwvpenwjcvfyz
ncrs
aalgvusyrp
vytwskdacgaxqyqaeajidausalid
vvuvfnahgftwnrbsidhshbgvh
ftafdsguhjkwcugueszhdlpwfggmjachsjuljryjvbyj
pqvets
jtntnudrfuthtuecoctepae
gplquengdkkjnehsipzkaxbad
jbzoqyhghcbadm
nssnuqay
afdghcwiqridyzfytjgxtcfvaavduzooausahwkvchwizydf
cnipdittekmhs
gpodpmukwbobnox
wwhzrnruybzvzqmidojjxnhyvyfwhdyairwemnsxbvv
wubbgbwir
ffplvcwnrtfdnjhaugyibhgvpqxpoyabjnqtuqgndooqvfrfj
ewrdqfakxzbtxtiuzmkttsuecsheskwkbzgtbglb
y
pmtewbtunckfqqhdizikpxvioxzgkrataxtxsdnnntos
nyaiaaftulggcrljkb
kvxjfvei
pqryrijskqysronzm
mcqedk
pyvscxnmotiparprlunxgtcafv
wajbvnyjmwhvnzqit
miaifjagagfkzisgyvwkmseab`

func runExe(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(s string) string {
	target := byte('a')
	b := []byte(s)
	for i := 0; i < len(b) && target <= 'z'; i++ {
		if b[i] <= target {
			b[i] = target
			target++
		}
	}
	if target <= 'z' {
		return "-1"
	}
	return string(b)
}

func parseTestcases() ([]string, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("parse t: %w", err)
	}
	if len(lines)-1 < t {
		return nil, fmt.Errorf("expected %d cases, found %d", t, len(lines)-1)
	}
	cases := make([]string, 0, t)
	for i := 0; i < t; i++ {
		cases = append(cases, strings.TrimSpace(lines[i+1]))
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to load embedded testcases:", err)
		os.Exit(1)
	}
	for i, s := range cases {
		input := s + "\n"
		exp := expected(s)
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
