package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solve(input string) (string, error) {
	s := strings.TrimSpace(input)
	seen := [26]bool{}
	distinct := 0
	for _, ch := range s {
		idx := ch - 'a'
		if idx >= 0 && idx < 26 {
			if !seen[idx] {
				seen[idx] = true
				distinct++
			}
		}
	}
	if distinct%2 == 1 {
		return "IGNORE HIM!", nil
	}
	return "CHAT WITH HER!", nil
}

var testcases = []string{
	"ynbiqpmzjplsgqejeydtzirwztejdxcvkprdlnktugrpoqibzr",
	"cx",
	"wzvuatpkhxkwcgshhzezrocckqpdjrjwdrkrgztrsjoctzmkshjf",
	"fbtvipccvyeebcwrvmwqiqzhg",
	"snsiopvuwzlcktdpsukghaxidwhlzfknbdzewhbsurtvcadugtsdmcldbtagfwdpgxzbvarntdichcujlnfbqob",
	"dwmgilxpsfwvgybzvffkqidtovfapvnsqjulmvierwaoxckxbriehypltjvlsutewjmxnucatgwkf",
	"huomwvsnbmwsnyvwbfociwfoqprtyab",
	"kjobzzngrucxeamvnkagawyavqtdgdtugjiwfdpmucaiozzdieuquuldeiabbgvi",
	"klsbxwtupwuounlrfgmsjaeeikkzlwckytbbifesjlmrejdpxhbjfqxcjmkjnddrppkzzkdp",
	"wpnbjkxvefusmzu",
	"czcgxhbmadmr",
	"jopzswvgnclhisyfngldcwaqoyvgdpmigubzgtedgomlredtpesmuvnqpvkppuvgrth",
	"kw",
	"kkbqeitzemsjwwzpczcqbchebjayokfzeuolqmqqbscvzzqytcxnygjrtnpzmtshzavaxfjqsikcpijynmzmbfuehjxkbbpn",
	"ptwcvwezlnbtomobdpy",
	"abtteukdwrulgmzyypdb",
	"wotukudvwtjzemzjxvzdqzgbzmolygolzucbbpiaqvssghcyuyqwqnqjdensncdncdnyexazonvnapk",
	"iclcdlwallfahlctegagvvxdxajlwathefodplwieaglkpjjrukfscdrsjfmeezhkqhhyfjlnvbetamcwcenjrnxesnju",
	"chouluqbmnanxkogljpcfzdidrtwezwomfynnfhokqelou",
	"cpygjawotoagjdyujrtenwypcvpyhrymiuadivbaimqswmodxiljyvgtcbczijrkdqhyfcnjjqesqugrdn",
	"rmxyzijolsuefdwdmmsoervjlupxngppwqkpubojexpbtgalpmaqcvcvxvmalbdtaiuwjxheysjgdnowkm",
	"knuvneoweqkegfolzmnzpm",
	"hzgogswbmbhucflbxuvfhtjtcwqyjylnobuwqvurxnsopiwpgkibbbflajuaecznvhtmrhogktdtczkkrokiaqbglcgqlg",
	"ivxxjjqmiplwhbjrcaopxobznp",
	"odcchdyengotcnrymbfhpheilkndrjtrzgwjyoqtoruiihadtzwdfxnhgjx",
	"axrqnbdmuidxslhvwwrvjhxhcqjvkhlupjsfezarqklsuazemefqceygzypsywghxehymltseupdtaqtlpoja",
	"rufvpzxprkwcietmwgkzjmbgbkxxh",
	"ovxvvhilvfjalsrbxuelapubahbahukcblvnegoneljf",
	"kxzxnmanirrzxvwoybsdnmfaqetvqxweckhfhazfxzvrfwcntdtuowettbikzxxmaubpcljveohqlfxymkiz",
	"majqjrpbyrsrivbomxdmlpbaixbivvswyjvygyqqkmigdskzhsvxrvlfekxasbse",
	"ljujkpzmtnfazesboekaxpvvyixtgcrnifqfcvufsdquz",
	"tmynijjanywiirqrkkgwznzeayqezvwzsmlobrnutzzyhalqfvguluwpaxxhshifyncsoh",
	"woqzzwdgfocnumiinyyltkcjapayigymmnyuuvmwbsolseswikampqebcsllacgwdvrpbkakmeyuinvetemjqbfeepwuwyxb",
	"qbrwxwumflzsccrfigzikwwiziqoeyorebusfuqbykcguzothzoqfwkuepyrbrcqkaycdntlsokmqludekafxeaktgbnub",
	"jzmbtyzwflcnbolttyivjszonfaozigmclddalafmtwuakozrxwppcbrmzziauqdckldpbequzjbamkfrweffyfuvh",
	"kapuvmbhhujkfhlhfnolsemsyafsavmwzfeaakqabbydsteyevzmzannsvkwhelqgrmcensvldnnhpmh",
	"hupmsciiqlrattyphibtkzmudrbewmaynxmndwotoffkpnfsjyq",
	"lleulypuqybgifx",
	"kjmubjrnbvnimxgleedtlfansmocuwwvcynrxrefegfhaqepltxjwzkvdznyifkuqkrwemxxrjh",
	"lmpqjnndveeastuwqdwyuwgtutqdiwxtfmczbadlxpkdvolsw",
	"vpzyhfrrcyqfaufqntvgowvmiasemfosb",
	"zcusmkhqobptdiqprumpifhrlfjteoccpmsnrciphdjelzdev",
	"vesrgabm",
	"wytpwwdprlykdvahhpjihaplqkccjsnhxlmyehjgxypvzljmteydmlqphuwlulnilmyywjdp",
	"doelhxfkphdvmmoqosthvmqjphkqvacpkmhnbs",
	"bncigxkfdfzwlahbamraedttgcxogaqtncrfhhnmpzangmybtiaslwlkvouetqcidwdiawetyvemgsvkgnqqdrdwpdqopforken",
	"mcsqkhohlpynaoxarmyyohnhippehojlu",
	"teqvcgjquvdbekbktzfyomvgnpgfmbzzykyvqgxkfqrwtsneuwtpvsgtobhptkr",
	"gabbehovthvvdyqnuujvkpvgflrlwomovkocehydyemowrglwbamgcnssznrguavequpjmrrcctxpyagntyljcxayzhtas",
	"leoidsoihmmxqkhsicbekusspcqpnkmxadntveyeadwnhbuiimdljltzugxrrbkqra",
	"pauikeldtrnzycisdqcse",
	"bqzhldrwgfwvhmgtwxolfixlrggokeirvyyacvkjwzkbctajzdfwigqrlocnqrn",
	"csokfjhunvwjbaau",
	"cvkkrujbzwhwqeitbzmnishikzrpowdnlfq",
	"hwzpmah",
	"jdfeoaylqnd",
	"hpjimnpwbfwqlghehzjzefrirafhxnwfhyresvfxhqyjzzwxaxozgbnfmqmgmrqpdecr",
	"ca",
	"hxmovinlvmixcyfshxgbybqzcwohowyyvetuqjatelvlljquwieapzxgtmkpjoepjzwklutmmsmwqjvyshgngiqfqeztqyry",
	"upkinjouqbzxgiowcagtdizagftfvltziavvaatfjhmpxllcsnzdwayfatrilxgmg",
	"fzfpvetsqnsuhabntbuuueqcnjkuutgwbkekvxipp",
	"kiucozyfdqirwmtonvftnyevvyzdsydfxpjvxdvyniizdulqrzcabtayednmzpakuecwmusxhbgnab",
	"zmerxwjnmrbxorereugsfcspjiigmvaubqyycaujeyeyqedkaxkeksycclgiokvpprqcckomymrpr",
	"qcheumqkeceebmusdvcqisjlvxpmskxvtoqlhyqbxmzlkzqahscpdln",
	"thraqqdbslblfkippynazuyhuupxxlvzzxfbyxpqvfmwjvspgawottvftoksknsyueqnmhlmjhgihrslohtktuekheyirfrwwtd",
	"pgwftlthwzcrpgjmiqglnkxabggtfrociqjwpnkvxwwmbdbkfyitxrojrn",
	"lhvcxczucqavsqkxtdshabzdxjcsmweystqgijffsheehipzjextigcevnfgsa",
	"rvpveyxicxwtwiztxhaegxsvaneeprbqkhvfokvohuuxi",
	"wgdlpilrsjgjrmcfmdsudqewoenixrgyzfxciokmqblmxxmreywwrseubjtoqysngogmcdrmdqvoxufbghhdtmn",
	"rgqggpcgd",
	"etgrsqutntmuivzpuutzallyzgwubgogqvgaehopctiujmkxxvwjvctdplttuchenh",
	"pje",
	"ulxvedqkkvxclzggjqwnuxxmmbeocjuveqgtplesesnixgdfeornfbjypscvluddphq",
	"awvzsc",
	"ixsiwrzmfhudvsxyqfpuuhahewtvyeueyzstfcdiditosyisnbwceoblkiun",
	"ajuhgycsxehqgokayinsqrewkwsurmgjdouuihzuiyjcleoamu",
	"buyyjqtynotsxqgfaidmwigppiqmetpifwagkiwwdvanjktxgknyeudhly",
	"pofvlnnnlsdcrzoedazqwmuphmuxdnbvjtfjiqdrvvkaeackykoghlvsoirfomadadkxgjrdqyjcdrtjxbhdvfoffoiwf",
	"wsarmaddjcuymgixmhoutqrlfnuzasnthjlqkpuyxoxwiliyffhokualp",
	"zwjvhadoxqfgmwtkrebpnozolagsgcpevcxqwvrmtigcnarpwxkxagtjyyznowqwjcbclmunmzuwyuyh",
	"pmraunujvpemvnmwysxtjprljnsuyflbnoojc",
	"hkhmatkefpcydftbhufauqvnvnfjmojrbfpfdvdpxjxzjjmbsikyocbyuwjhoujrgmcnbbzclajjbabtbetfpucs",
	"axrezdcwctyvdgdpnc",
	"pwoqxfrdelxnqec",
	"qbxrhdxzmlqkrwsrmjhswcetmylxjlqlaaaeixvwguefdpxpveiiebfvlxkgudqzpeuikgudtvwhlzmyd",
	"ymrlbukhmbynriyzfibxixckunayvihuywfi",
	"qudfpgstscnjwvbiaaztfkfmhsrgburygorvostklldyzanl",
	"phxvxuzxxumytavgtxpnymfe",
	"nmoxxuxgicwfsajxfdeifmiqunkrxnpoditkiixjmtqapuqqgnixuynxbanrcziwzaapvzmrlrpqorxxwpsfzsy",
	"tlgj",
	"chdfztouepkowepzipkwqxzhyblqimrutdjeqeeckgbappcsqbeuoses",
	"cjleolbkddtvdjvjslfyavypmmugcnipnv",
	"menzgynuygiczeokljzzokhfoewkzhpedrstoxkfcvkfilvwddmkykjkyagvujuf",
	"rfphcwqiwdbmxcmopveweb",
	"eudmtdjcssdgxjm",
	"qowiykufuzmiffxcguqylvx",
	"vckqnzzjozbelfphhceyondlyhdkeghtrkvpksupqviyjatvtfigqrwdjof",
	"boqqrtsrxvmnktz",
	"hkpeqxxaxpzeacraeydbkjvgycvbyhxfzlmuqtvlngmx",
	"txhsao",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range testcases {
		input := strings.TrimSpace(tc) + "\n"

		expected, err := solve(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}

		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(testcases))
}
