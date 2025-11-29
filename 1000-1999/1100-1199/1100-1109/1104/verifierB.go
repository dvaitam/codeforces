package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesRaw = `szycidpyo
umzgdpamntyyawoixzhsdkaaauramvg
xaqhyoprhlhvhyojanrudfuxjdxk
wqnqvgjjspqmsbphxzmnvflrwyvxlcovqdyfqmlpxapbjwt
smuffqhaygrrhmqlsloivrtxamzxqzeqyrgnbp
srgqnplnlarrtztkotazhufr
fczrzibvccaoayyihidztfljcffiqfviuwjowk
pdajmknzgidixqgtnahamebxfowqvnrh
zwqohquamvszkvunbxjegbjccjjxfnsiearbsgsof
wtqbmgldgsvnsgpdvmjqpaktmjafgkzszekngivdmrlvrpyrhc
bceffrgiyktqilkkdjhtywpesrydkbncmzeekdtszmcsrhs
iljsrd
idzbjatvacndzbghzsnfdofvhfxdnm
jriwpkdgukbaazjxtkomkmccktodigztyrw
vlifrgjghlcicyocusukhmjbkfkzsjh
drtsztchhazhmcircxcaua
yzlppedqyzkcqvffyee
jdwqtjegerxbyktzvrxwg
jnrfbwvhiycv
znriroroamkfipazunsabwlseseei
msmftchpafqkquovux
hkpvphwnkrtxuiu
bcyqulfqyzgjjwj
lfwwxotcdtqsmfeingsxyzbpvmwulmqfrxbq
ziudix
eytvvw
ohmznm
koetpgdntrn
vjihmxra
qosaauthigfje
gijsyivozzfrlpndygsmgjzdzadsxarjvyx
ecqlszjnqvlyqkadowoljrmkzxvspdummgraiutxxx
gotqnxwjwfotvqglqavmsnmktsxwxcpxh
ujuanxueuymzifycytalizwnvrjeoipfoqbiqdxsn
lcvoa
qwfwcmuwitj
qghkiccwqvloqr
bfjuxwriltxhmrmfpzitkwhitwhvatmknyhzigcuxfsosxet
oqfeyewoljymhdwgw
jcdhmkpdfbbztaygvbpwqxtokvidtwfdhmhpomyfhhjo
smgowikpsdgcbazapkmsjgmfyuezaamevrbs
iecoujabrbqebiydncgapuexi
gvomkuiiuuhhbszsflntwruqblrnrgwrnvcwixtxycif
ebgnbbu
qpqldk
erb
vemywoaxqicizkcjbmbxikxeizmzd
jdnhqrgkkqzmspdeuoqrxswqrajxfglmqkdnlescbjzu
knjklikxxqqaqdekxkzkscoipolxmcszbeb
psizhwsxklzulmjotkrqfaeivhsedfynxt
zdrv
wdgicusqucczgufqna
lpwzjhgtphnovlrgzpxcingaxrymqpcmtqzssn
loa
jwwuardjqxkyr
srjqnrqntusjojeqoseryfiuanxvsblnmjvyvaccam
oizzluxpykmozdplen
afileszjni
jxnwinkypgwpmwnccegehxadiepydmuxf
c
tbrgrnlbudxrvnvxdivifpzzwbzg
ucmdvojvqpmdtpdemtwgfqinxrjpuzrgzytkpdayxvlw
bruojydhqiiwhneeig
rutbrtqeniipwjipgpltphkftyf
sworebqkqweuyzgktppkdeewihcurwbsfvdhsgqsvjnkayaj
hcxhivukitxqmadklediyevsblccxdjkhiqblace
lxuwhdvkiaqkdlzzuxetimcvst
qpsnrmjhujrebtqdfhgnirairiqipemwdxlcurlrrzxqvsa
joveecsevgpzykljfezmomdteijvvzutarauemxr
oayntvn
lnmtobdpybuwwazbds
qqylrizsu
zpwhzthdrlfdybwknxlivuy
tnnm
jykozwhutqebkvdqfruupky
dsapgmufmwhdhkkvhzvoxplpuyvxgnomrdspieeamndzau
foymv
zjeeqdiaomzuwxzhrwmarzhnfvfkvhcyr
ffmsaqgnhzbqxgwqwturchmyodsubmcrdupb
hyaajoixnfterwkyruoqznrfwmwmzgpile
sifyxtcxlkeiiilmi
oaeeihgczsrtgrnwhseromwgcucezvbaxmmnv
stevrrchm
jgvxmlxfh
welprjcqjgwoajzztsdt
yoitbbzkfzeuddnushxgqqmd
gmvqewsixawdzgysmvprthibufvvrqhniyvnmipdvefrao
bpgmxrkhdcvxbnogftqgqmqlghlvsyyckbobtfejpbsqcsmcm
sujmilpbrpanjsxkzetsrictzzylnmqza
sdbsqadkklyrbulscpucrokqzrafklgesesdm
qnlkitlbwcyuhziymrjsz
ccwfincejrxuihgdixpbxqjzzgrcrkkjqebolzxb
kn
frbwswvuqnfghdsesqdxiogzbloktxlhuaapbfirbahycqfb
ggojhpqlkmucgtfgvtjsntplapadvusvtn
skkcungwqzptsvrqptvxsyotpfivqjsyzmtriijatybzoo
hqogwpkwuemnbudlzaiyrxb
makkjszbgwckdvuceywj
tkhauwwfyyyqxsuljjmnqozcgnh
bthuhhwmmgtexjxxlawwvjopfvealnrkzqpktds`

func solve(s string) string {
	stack := make([]rune, 0, len(s))
	cnt := 0
	for _, ch := range s {
		n := len(stack)
		if n > 0 && stack[n-1] == ch {
			stack = stack[:n-1]
			cnt++
		} else {
			stack = append(stack, ch)
		}
	}
	if cnt%2 == 1 {
		return "YES"
	}
	return "NO"
}

func runCase(bin, s string) error {
	input := s + "\n"
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.ToUpper(strings.TrimSpace(out.String()))
	want := solve(s)
	if got != want {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func parseTestcases(raw string) []string {
	lines := strings.Split(raw, "\n")
	tests := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tests = append(tests, line)
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := parseTestcases(testcasesRaw)
	for i, s := range tests {
		if err := runCase(bin, s); err != nil {
			fmt.Printf("case %d failed: %v\ninput: %s\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
