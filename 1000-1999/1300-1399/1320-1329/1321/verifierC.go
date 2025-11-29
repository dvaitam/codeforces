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

// Embedded testcases from testcasesC.txt so we do not depend on external files.
const testcasesRaw = `ynbiqpmzjplsgqejeydtzirwztejdxcvkprdlnktugrpoqibzr
cx
wzvuatpkhxkwcgshhzezrocckqpdjrjwdrkrgztrsjoctzmkshjf
fbtvipccvyeebcwrvmwqiqzhg
snsiopvuwzlcktdpsukghaxidwhlzfknbdzewhbsurtvcadugtsdmcldbtagfwdpgxzbvarntdichcujlnfbqob
dwmgilxpsfwvgybzvffkqidtovfapvnsqjulmvierwaoxckxbriehypltjvlsutewjmxnucatgwkf
huomwvsnbmwsnyvwbfociwfoqprtyab
kjobzzngrucxeamvnkagawyavqtdgdtugjiwfdpmucaiozzdieuquuldeiabbgvi
klsbxwtupwuounlrfgmsjaeeikkzlwckytbbifesjlmrejdpxhbjfqxcjmkjnddrppkzzkdp
wpnbjkxvefusmzu
trgmmcbbmcnnnrtztiueeskdqnhtwxqdvdzcumfapnkziwteyrmpxmclhtdmafjhtcguetxiifgluzdgzpvrldri
aqsv
hjcshmedpchlhifkq
wlpzgejxdxzdahpbvmnssrmcdbtyikkbapmtuqkbwubqqkdfkvplnjorivxuqtwnpuxuyeqebsfzcsdxfzarpwfqufdgodqnpxws
glslmwfketnawfoqizfzbdnizqcgbycjevxbzdbvittphjbgyvtrgxmvolygicducoqeudprltguktkmzzqwiocwgypuswzdo
vxhddhymqdxopibkyrqazhyjfseribhuwthdgqsysmjjlcednezboykggweedehudnulkgbqsqqnyoinryqizeqojkiibeheip
ungydk
gndrcwnvbfmmhciizwrznehgjpahjsmbm
hjrfbdjnlhrkitltl
dzjkmdyxmvvmnpfoyxthqfxztlblbkfqmzylhuzyvwehuabmxzq
cfnedkkytgecvkzjimelykzyhkdapahkrlygbwyikwtculqrfcfhoahpluluwswkzrexdkszmvwmofxrevlsli
smbvlhzjrtqgvdwfknmkijtivujxtpelthnwlegkypskmnkptyifevplkfkkbcesxgvoqrn
frxwzxjmiqwsdiexlleisgqvnwjjwotkwcyxajeiwbhrvhxgjtsorkthvjhietgaamvzjajywdzmxfkhpicm
cswycjeywdkl
ktvnlrayxacnihzrfebplbtvffahlzdqmiauyzjkvyhmzfqqphtxfkzumfn
lalrcaclwqgefysnqpidsljyhxxrhlfypoexlwqttoqsnhiieurbrrqvrqodgtdhdcnfsomkxgklabwsmdu
koafqrhllkicmehsktwdhgxqfckhabzanuoqi
eewcsrcsbogbbidx
kctzwxikrnnknvttdjprjrtnfkolkfh
tugmhwbgmmpecnvfznyzrvolistbkyhyxxtyqzvxhuikpasoxjqezmyddu
mijlw
urrnprezoyhwyjttlkuemmvdevskcofmxwlybcuctuuvptnevapbhmxeyuvxwsubskbzmqnshtdfkyqcsjyrvcatbzhqontuu
uh
bsee
uaeuetcxlgtppduxvonnldkskeypahfwszydsvtxqponcmbeejdtuv
elxtkxfikdhxfzmyojhoctwauhvzkqg
fccikcjjoeijpscbludwkeilbgmphgmcayfofkyigugtoehupneveuhimcziwow
foejmagy
igqrumzikbmhaawpqhbqzooiepoylmwvzecfeizhkybwvavixppilwmpmekdbahbumkjlovuwlhqku
iekzynvjrjkcbdjsaavouvpztllejplbfx
ggbqrfcdpxywixughoqyrwujbfnbqfcizkpunzrhqcnhwhqhennuevhfcthrithqyuwkpqredbnbeskpshz
fuxrbbnlnwiovucyyhhlyohxzlqnthhtpjvqtporbzyjbscfywgvzdxaxnur
ybcdzpjetrbbtxykgesgbcgefcqgefenhudzbapvyqt
nqlfjbkgnzkjfjbtd
mpekneoupdfdpfcwxnoesxskgffofpffzlptevfclocpsxwikeznilpocxv
zlynwrktnujszqwkdsnjquazkzkugtrggzpnvzrlquvgsdcmzviywkxttacdbosgidlx
fzfktgzkgbonyzekqjmsxxykwgvxoyykmzmbivhazchyzptbzuxrhcugsdrdvthjkghvwbmactperbdtdcoskec
jbpzwgenmtgbmszkcnbr
udzmuaidstouprahjbyuvztyvfqewamyqbsmrbzvwktyehipbesplfgywzplrhkljtwnlshvjldqrndanzglbpifbfnktn
rjjvv
ypvzorvfshgumwgkdjskxchjykwoybmdjqprnwkpblqhgzkklwctv
qwffyarkmanwyogydqqx
bpmdjcwxgkiqpdnvymonmbbiugzboaqzpe
jozttymysbwkpmrtvzajswycdxqzluhupbv
vcwcnzsrzgnfaq
lhprhizzpjfqufyjzrlvubljvdvjnzytfodgmlnuawrpjxmymwmjvzbevvkhxylvw
tqqkhrsdqpamofxngzgvntbcshgjcjfokygpgpo
dfqlqfykq
ruvjadxezlsxsongmeygxaqazsslinlr
mbvuwrlkaywcfxtgr
mkhyii
dsrldkwbcaenynuzdojpyslqdkqcaypidhqgfvraytkhhpgfrd
ysbjbdeojyqobuusfwfaoywlygmjhcxhzgvamymhbzlnprtnpvduxqnkitveepgqhhqreafqgouiumdjxcgju
ufwxvlbpruwheampyazictbbhdcloefkymogxkdnvtlfhwjzkslnwkes
gvexlkvzgqxftoaskqjcjsrdighvlzfxmnaajuckmutnysdzrzuekxtqpqpuzxifskaqchxekwfbwul
rsshxfvyezftlmfcxopqrktxzvipbsdkk
zmcaruzprzsqohzlrlylrgiujvpbjhsdjlwbtzuauwmdpdhibyiufjruhs
wqvvkaijxvmzjzcjrdrdsqtazjhwwhrziadnurfyuymqrenkbcxube
qltamvlmclhiswcaukatzvkwhrpqalxftpmrecsbxlbpnqqr
tpsbjwjwregeubaajnyodacrqlnhsxmyjksfyqfycnwgivlqddggwgaijmnmtvk
ncusrwil
hwezzhxpznaofrt
duxsjvcazyfcfhjcapthromtmykapnpojztctevqwkkionqgnwzicdnwigqnlwyfqbgpqbovqhowayvqblc
acprlorrnr
uwjhkxxdtldohunulkj
ebeeplq
mkxwcuukwmsgrztxvultbsplczgrezhm
gvgjxpzhkmdllfbalzexeidcnwlsunxpnacqxw
ixp
hnxhzmyqzphzctmdejrju
lsymzmdonhoyaxn
jmelnlgqkjykx
krjekkzdjv
hunnfdxoseeemdjujkvsijawffav
tnyrcpumqjdcsjkqqmyj
xsjtmhbisbgwyjqrsoahht
ufzvfzitzhzigrgydrvsedappdgldbwxtgbvmaqbaqywg
zusaoehohgebeehsqnspkepzszvvqsmftguhwlfdozrnwctjywamvpdpk
tkuwrjaadtrdijmbktyp
jxywzngdey
rfrynjezahmdifczqbhbkfcspcsushkzva
xdzftgkjnnsuddquigmgibehmohlad
eyoglrdpzmoagtxepdxapstuglagohyqusdhsutetsueks
vhfgxqiqghbgskbzmllioalqqscleupkyugoclpxoukfvuxhszyixiqzvjratdpavwjaipykc
kdfbffnklkbjqrngvbhymg
yf
matvlcwbkivdnjvrpzpywjmvksedpkkhiufisqlfewhgxrexuwifjvyktydhtaozbaniztrvzck
wwzpehmnk
fkaxdhhg`

type testCase struct {
	s string
}

// referenceSolution replicates 1321C.go logic: remove chars if adjacent to previous letter.
func referenceSolution(s string) int {
	bytes := []byte(s)
	removed := 0
	for ch := byte('z'); ch > 'a'; ch-- {
		for {
			idx := -1
			for i := 0; i < len(bytes); i++ {
				if bytes[i] == ch {
					if (i > 0 && bytes[i-1] == ch-1) || (i+1 < len(bytes) && bytes[i+1] == ch-1) {
						idx = i
						break
					}
				}
			}
			if idx == -1 {
				break
			}
			bytes = append(bytes[:idx], bytes[idx+1:]...)
			removed++
		}
	}
	return removed
}

func parseTestcases() []testCase {
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scanner.Buffer(make([]byte, 1024), 1<<20)
	tests := make([]testCase, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		tests = append(tests, testCase{s: line})
	}
	return tests
}

func runBin(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := parseTestcases()
	for idx, tc := range tests {
		input := fmt.Sprintf("%d\n%s\n", len(tc.s), tc.s)
		expected := referenceSolution(tc.s)
		gotStr, err := runBin(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(gotStr)
		if err != nil {
			fmt.Printf("test %d: cannot parse output %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
