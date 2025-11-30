package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcaseData = `eszycidp yopumzgd pa.mntyy awoixzhs dkaaaura mvgnxaqh yoprhlhv hyojan.r
udfuxjdx kxwqnq.v gjjspqms bphxzmnv flrwyvxl covqdyfq .mlpxapb jwtssmuf
fqhaygrr hmqlsloi vrtxamz. xqzeqyrg nbplsrgq np.lnlar rtztkota zhufrsfc
zrz.ib.v ccaoayyi hidztflj cffiqfvi uwjowkpp dajmknzg idixqgtn .ahamebx
fowqvnr. huzwqohq uamvszkv unbxjegb jccjjxfn siearbs. gsof.ywt qbmgldgs
vnsgpdvm jqpaktmj afgkzsze kngivd.m rl.vrpyr hcxbceff rgiyktq. ilkkdjht
ywpesryd kbncmze. ekdtszmc srhscilj srdoidzb .jatvacn d.zbghzs nfdofvhf
xdnmzr.j riwpkdgu kbaazjxt komkmcck todigzty rwpvlifr gjghlc.i cyocusuk
hmjbkfkz sjhkdrts ztchhazh mcircxca uajyzlpp edqyzkcq vffyee.k jdwq.tje
gerxbyk. tzvr.xwg fjnrfbwv hiycvozn riroroam .kfipazu nsabwlse seei.ims
mftchpaf qkquovux hhkpvphw nkrtxuiu hbcyqulf qyzgjjwj rlfwwxot cdtqsmfe
ingsxyzb pvmwulmq frxbqczi udixceyt .vvwcohm znmfkoet pgdntrnd vjihmxra
gqosaaut h.igfjer gijsyi.v ozzfrlpn dygsmgjz dzadsxar jvyxuecq lszjnqvl
yqkadowo ljrmkzxv spdummgr aiutxx.x qgot.qnx wjwfotvq glqavmsn mktsxwxc
pxhuujua nxueuymz ifyc.yta lizwnvrj eo.ipfoq biqdxsnc lcvoafqw fwcmuwit
jgqghkic cw.qvloq rxbfjuxw .riltxhm rmfpzitk whitwhva tmknyhzi gcuxfsos
xetioqfe yewoljym hdwgwvjc dhmkpdfb bztaygvb pwq.xtok v.idtwfd hmhpomyf
hh.jorsm gowikpsd gcbazapk msjgmf.y uezaamev rbsmieco u.jabrbq .ebiydnc
gapuexiv .gvomkui iuuhhbsz sflntwru qblrnrgw rnvcwixt xycifdeb gnbbuc.q
pqldkber bovemywo axqicizk cjbmbxik xeizmzdv jdn.hqrg kkqzmspd eu.oqrx.
swqra.jx fglmqkdn lescbj.z urknjkli kxxqqaqd ekxkzksc oipolxm. cszbebqp
sizhwsxk lzulmjot krqfaeiv hsedfynx tbzdrviw dgicusqu cczgu.fq naslpwzj
hgtphnov lrgzpxc. .ingaxry mqpcmtqz ssnbloag jwwuard. jqxkyrus rjqnr.qn
tusjojeq oseryfiu anxvsbln mjvyvacc amioizzl uxpykmoz dpleneaf .ileszjn
iqjxnwin kypgw.pm wnccegeh xadiepyd muxf.acn tbrgrnlb udxrvn.v xdivifpz
zwbzgvuc mdvojvqp mdtpdemt wgfqinxr jpuzrgzy tkpdayxv lwibruoj y.dhqiiw
hneeignr utbr.tqe niipwjip gpltphkf tyfxswor ebqkqweu yzgktppk deewihcu
r.wbsfvd hsgqsvjn kaya.j.t hcxhivuk itxqmadk lediyevs blccxdjk hiqblace
mlxuwhdv kiaqkdlz zux.etim cvstxqps nrmjhujr ebtqdfhg nirairiq ipemwdxl
curlrrzx qvsatjov eecsev.g p.zykljf ezm.omdt eijvvzut ara.uemx rdoayntv
nilnmtob dpybuwwa zb.dseqq ylrizsul zp.whzth drlfdybw knxlivuy btnnmljy
.kozwhut qebkvdqf ruupkywd sapgmu.f mwhdhkkv hzvoxplp uyvxgnom rdspi.ee
amndzauc foymvqz. jeeq.dia omzuwxzh rwmarzhn fvfkvhcy rrffmsaq gnhzbqxg
wqwturch myodsubm crdu.pbq hyaajoix nfterw.k yruoqznr fwmwmzgp .ileisif
yxtcxlke iiilmiso aeeihgcz srtgrnwh seromwgc ucezvbax mmnveste vrrchmej
gvxmlx.f hjwelprj cq.jgwoa jzztsdtl yoitbb.z kfzeu.dd nushxgqq mdwg.mvq
e.wsixaw dzgysmvp rthibufv vrqhniyv nmipdv.. efraoybp gm.xr.kh dcvxbn.o
gftqgqmq lghlvsyy ckbob.tf ejpbsqc. smcmzq.s ujmilpbr panjsxkz etsrictz
zylnmqza ssdbsqad kklyrbul scpucrok qzrafklg esesdmkq n.lkitlb wcyuh.zi
ymrjsztc cwfincej rxuihgdi xpbxqjzz g.rcrkkj qebo.lzx baknxfrb wswvuqnf
ghdsesqd xiogzblo ktxlhuaa pbfirbah ycq.fbqg gojhpqlk mucgtfgv tjsntpla
padvusvt n.wskkcu ngwqzp.. tsvrqptv xsyotpf. iv.qjsyz mtriijat ybzoolhq
ogwpkwue mnbudlza iyrxbjma kkjsz.bg wckdv.uc eywjntkh auwwfyyy qxsuljjm
nqozcgnh tbthuhhw mmgtexjx xlawwvjo pfvealnr kzqpktds ujzrvina .jycupdq
htxuxinl zhbdtqqq fejbcgav bnxwacba brkkzata rgpgijsr qihfgmbh rwobkknd
asfqucyf ghfjzdbz kxec.oeh bxjlbsco gzhvfdbg bxxdczzx hjwiqnhx bxiygkll
oyvtmvmc nh..pkft udhcyzni rjky.lno llkmpqal ejfjserw xefouuee fc.tihlu
kfipjcne rlodevkc vfprbbxg ulxlqlzq uzvlkudf mbitwzgb h.jksmhl ybhjwsag
dehlqief hcjsqqrt rznosqpf qlgnzcig hyeeygaf plfbzlct hvwgcouu gtkfsw..
vwagkprb blprlepc qkvxsvjt kzscpknc icvukafk hki..ijp najfujbd nntgilyu
xspsjtiv fkeldmlq xswgmoe. pwhbxuhc xcbqqpsp wkqzfswp mamrxr.x ofsslb.x
l.lohwuv rjcoylgf eo.blskz fsppasht boufqgmo dkiefkef zxtqjhrw nooqrjfq
tqjs.zgj veva.kdn mwuqxfto o.rol.gb cxddrmeo mfpoqsbs gsopmjly yf.tifya
rbzvcrho kokxdmbx oinokqdf mrntxpqe keletghz zgouedwd nboelrki mampwojx
wjusmkyj fdpfoeod rdrkk.pv rukxskrs zokpwm.r gfhrgthb yktybkna lllttvng
zjhkmwmv yfamultz yt.hhc.t kmgwjdna zlcznedr zx.fykem nkruwqig gffrfedo
sqenektz xwvkteal yfhhwpsp bucerpse glweixlc mpaqogxh gwzaxwjb i.qgczd.
zydmkdow sqwupvie nlulymnn lrggcehh ahvmozto sdbf...q abnzineh wyvlnyks
xbqoewql sbld.huu dnezalee japuapcy zsncprtq dervwmut rnhqmp.x kodcgstw
lddldgdw uscaqnhc jptbsnrj mubvtait pohikypo rbiqfxwo ojssfkqv myvwnvrt
mpyuhjac ep.d.li. jzrjedqe obopxskr lewargyi tzczojau ixqwasmd dvkttuww
soctpqks vbgfbtdz bdrqjyzg frehgcql w.snitej szhctiba ntj.pnn. zcfgyvbu
.ynnllqe fzhhzblc ok.ghiew wqmdpvxz tapjiyzw jgzewumv bzymorae hpudjwtn
gqkdhhps dfplwutu tnmrn.ya umenebjm tnudgtip tniq.ydk zerwrziv varvxdyl
oiydjezc nwmapsxe yyrmpz.y hqamzbnt chvbocjt blybcc.b sjljcrpt lkyfulqh
kthhuywg jj.rkwjs avpivzhe hfcimgef .rz.tckm wgfbogmz dwjyhxuj qzuokccc
.hdqow.r oatfonrd gahjgqtj jilijbaa uyobgcko vjdhvdga guettvva oxarh.pf
.rahecae kscqrigm arilirmm qqroicfy psmetgqa qbkehkmb nx.sp.qz czwbernr
mrisbggj wmjqasig rqxrfhcg pfbv.mja ezd.bwsn pf.gsozv dvmhcekq ppqvln.s
hoimlmzs hmtdfvtu zlcanspb yoduuyho lqckvbi. sqytkesf nvjwoxhp xmaqidji
asckuqvf hjxcfolm u.yozvpv v.dspscv bbaibijf rptwvkao khhlxwba oqgmefhc
mbfkaor. tqfb.nh. ivqog.bt .wmnmqni oksacp.x xnfnf.rq yqxqtfin xpjlwomr
mjhlrrzw wqhiavci wmfiyzsi paf.pdhe dmbfcdor xuozabib qpxuglto dkkmumjc
h.xorlnn wxxnsife bklmcusz ksfeyx.u drgpwhlt quwfygjf xeumnplw ybrcalhe
gmoqsint kpkcstby erxpfcac afigxomw rqiwvirm dwmohcxx kevtauwm ubjlyvaw
toksayrk xzmwyxbs .vovwudn mxdsaayr tnylfmxb ezjqwtn. ufzspxjs ztixvbym
rsnekfom srvequct stmimpxb zuxjf.ui midiad.v z.doeohz hbhcd.dx bsvdbine
.ldb.zmz tt.hfrsp fl.tm.qz zsvfkqcu zybazsjd ocavbxir jstyziom dzuhjuvy
eqqxalwo dnzveidl iy.gkter htahw.pl uenvknto dibqjwqk ggh.xhml iapyqeny
pcqz.idh dnmedvo. qvzgfgil .wklxisy eahipytr aka.zfwg iuhcnwlw ly.ygdam
kskvznks zwimytiy ltcznhtp lyjwadtq bftyhyro jnmtacme xsg.pvmp dnzufwvp
gujrybjj zei.uqjp enkqkgib jqsjpjif jizkeimv ov.pxfmb csgkbqxj bnzdtuwk
ealhtlqn whqcbkau oafzxizv t.gznjuf bbpmrvvd mjnbhkns sptgsqvc kzumufyh
q.pczunv mgiyajbi ycftiyow njdjzbpf .irgebvm rasqjaxm kdiftwgc fywvsumq
sxahmvav aqnytzy. fbxymung fhctorrk vvigqti. mhvjtiwe wuilsxiq uhygrvad
gifxkhfu ubthmiig uimbybex pnjlmylt gjizipte zslembci cypgojbi k.a.wvtp
zn.nnylt xpygxnmj dcxfwkls nwmd.mbn tgdhwpmf vehtdwlk qozyfmup sfbqgvhy
ediraalj ..gbjvve ecwfsniy zecgyftn gzmrpftc zphgcveh gttweric yyslclqi
fvspnzrs wrhersdz nxnlhovr mkvsfbyb ltofvxos lleogrpr jtgethjd .vcuhnqg
v.jpwbmg uwbjxjtg nayoknhu dfqbxmf. aqpplrna utvnmhqa rub.ugu. kbeo.wfo
eeiqmwcq dcwm.plg bqnhpgfh gqvkuwjp ..zytszq hygjdysa ybtkczqu forcnxe.
rbzwekvx loxgmo.z cmlajgsl l.uadoz. ymjwfjxu hxkkgbba wftokwbi zsqicwuh
aeynlizp bmdjwnyh hqsnsiaa zeypwyel cthruree metwxkge etwdeybt iilaew.a
czownwmj vemnylzo ljpfiajh wbpbzact oazuhewm zxtmhsti fxgfcklc dzrhgkzr
rozywcns lfftdlfv ptcongcc vuikm.ls knqttucg nlqxplwd okahjnvv egiwqtbf
.sjbvdiu sdfwuros honbepzl qjxmcsne vwqhnpcs xlrqfbwg galhhwqq unrnfzha
hqrbuerc zaerihlk xedinltb rbvobxuk jjvxmjzm pjvdsvua wdncgdua hpcgkgjj
oowrswqg womcxacj xtogxj.n futvmmxo hhpajipp lydsxwvy dgwomgnb x.fywvmn
lqezcqwf bsygyqrx oujjpeaw onvvslnw lgyigvos p.ivtnji ocdkovvj q..hqkhe
fiw.hn.a vnmheccf otxmzzhj tmiajedx njxvjwnr bvxvedyf qp.ndblu kjbjzobl
jrugixif jqkr.cwb eemkk.zp fjayiaru tsanr.o. yaztq.md dssazmcp glsb.npr
kgaeppin uqdnorqj zcbnelzy .gcoldst kzvdgtkf fkycgxzj rscpso.q omlqyuqp
zf.v.eaf jfuuegye whozecpq rmmtuxnr uqvvnpip evvgmbit zeogemou vblwhvej
syvsypke ztuxcv.t mv.ccaav c.ceribg nkvivlg. vfncldno .kqdavbe nytwggcv
foqakwzt jvjzeobb jfyatkae idhizusp pgcejau. hfyvfhst oda.gslu ufi.dcjh
m.zjreje j.rdjqdy ygomuyda zmpuajwp plfgprgq xzruhenw xgwlphbh .dlcwbgn
kxnoo.ov txtuuolb gvie.qyx dzyn.gkd .axuhgwz .mgj.jvy lhav.hti jfzvwdal
xewsym.p oy.vdhtl bchfgnex mmslcbro z.sljwlx zkladmji bv.yqpib pk.wnxos
r.qhfqbm igkgxdtf nnsitxec ihwisfo. naejqeen boqxpqxb mdwzxjno durnznai
tbjikqae rbgkxdfj xnqyecqp qcnlu.ez omtkvjpr cqeaucgt twlpqeyf wenybclu
xwzjvixl ljnmpolk twnezewp iuwinspt bjqzplpe oeuwpehk ycvrsslf znunjiht
auplc.ip xmo.bnip qekegxmd dzkepuqz oqep.ebg wmkuxipj bncgaskm zjhj.nlo
k.vxidvl b.jgxsdt .xiaikdm ojrilepb dmoka.ou hekegpfj lhuywvgt btozivvo
oplnnchl bkdvmrww pkegczip pokorcpt ixr.udws z.encxe. mdnxrgea rvtd.dcu
kdvsruln vmcn.uoj xynomtrt eotpmdki u.eivbdy fxvbakkb fnotanrs mdvmbaeh
qpumkejt essereye shxga.gp rozlpunq anhwm.it ayqkdzq. .y..hnie qosyrlid
xxxcwlmo lnvqopmr a.b.pkyh dalcwfwp dheoxcfu iypwlfhb juorrotj .hajpegg
fjcifmuz ysz.jxyf omvc.zlt lhaz.rtf otejmekk i.aegigu abubojhq yudxcteg
fanffokb otqlwjse rtqbhdop hqkdvj.e iunidaal qfbkfbba ihaxjpxx tupkscgf
lwxczlel gomoltsc hhdhcxjq lcds.tlh rjyujzdf wngoeyg. cncmeshj itqxspmd`

type testCase struct {
	input    string
	expected string
}

func solve(grid []string) string {
	var sb strings.Builder
	for _, row := range grid {
		for i := 0; i < len(row); i++ {
			if row[i] != '.' {
				sb.WriteByte(row[i])
			}
		}
	}
	return sb.String()
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcaseData, "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 8 {
			return nil, fmt.Errorf("case %d: expected 8 rows got %d", idx+1, len(parts))
		}
		exp := solve(parts)
		var sb strings.Builder
		sb.WriteString("1\n")
		for _, row := range parts {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: exp,
		})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
