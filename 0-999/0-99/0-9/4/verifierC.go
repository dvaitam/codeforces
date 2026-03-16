package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesCRaw = `1 c
6 xz jitgtb vfnum qlroqi a okm
7 frfhh f feq lqvrf znxq zsl ofy
7 xouqhp pqq voo sxr opvhkw ti jjzw
9 qquts jxgp lvtcz xag dbsubi hvdyq ih bn ybbllf
4 acdcab aliefx qw m
10 z eb l uxxdj paj rytx i twep cv kdaoze
9 ympqk ki tnu w evbib ff o hqwbhh ocicsh
10 wlivn qya bm fdqx h d f dg q oojrum
4 ygxznn assbn sfdvz laqd
10 jwl avn d gyv z bnup gstc j l xch
8 ds mwo yl didd t umg datvpy x
8 loe ipq xxzn vjmh pt rnw wcsxsd l
3 ezncz v bejmhw
6 fqjd ry dkqh qiffow ml xeoox
1 mxfmq
1 imiw
7 uplrkw vcyxhr gmvmu k qwou da gxst
7 dm yzgix sgpzt at npiqsf wgyc apr
2 pvkoi oacty
6 yy ivuz ebfpmo jeajro l r
7 ogvjp epwrwj i jku uum cqugm qezqu
2 bho hqibd
2 zmlgkl k
8 fpo oew ugik dh gyvl le hi ummzx
6 xtq wxkxm wuywjr uvclj pfil pcfk
7 ad flc yunark tm jpuel gpd zg
6 enl ckg wh tbkluy befcn yiek
9 dkuyw mhbmy ptkr tcsqr pmwofn qobd sedv fcmjo aidvlh
3 e vcku bphc
8 ra wq bbgra kvqhe pae dhdog t um
6 umwqq fqdeug mg kne emkz zdr
2 ijqy ihnw
3 rvdatr gg sbue
1 iwxprb
4 tk w de wrfycv
8 jgfkwi scnnv bojvdu xiagnk rxm qxgny wf olmp
10 tgs ogyp kjcfl uphyt vsevjg jdaza kb riz vkocnp ajsseg
3 ty xcus iucp
8 es hgt kstwmq hugr i hetxmn o mpmj
4 hb qctrv b wnmh
9 dlq qzp cwoxw ia p e egkhrb eujyd rrcvve
7 ebjqvi brly vdt dtz lzu zpj teabk
7 alvrwb cwrqty nnhz ft a yxlvf abh
10 mc dtc hr da mcqisu bq qmnze ne xlbs qo
7 zuofp eleai feunsu opo nn yhl bmtanj
1 psivi
4 wolq ovhrr oj ndq
4 uvmdnt tqoc motyww lzrlfe
4 ufnopw zfniks juwz vkv
1 sbgo
2 a ktk
7 xk q pmtho t l j ewwb
1 mpu
10 p sghxtk kj stxp yicq sl ll hr xthgso mi
10 qf m xluf gxwtfp azxeg gatcoy gxfimt a djb lpmdc
8 ez ucie rczq zjjarr cn fs ogb kozbeh
10 ujz udgeha did ycq telenx gjryay k excryb xp vlpv
5 cue dacelw jy vrech xyxhz
6 f ngnuvb npu lfqdv e menyyx
4 itrwdc sxtfg tvwcrk qh
4 ovmodk bl m s
2 zilrhk hbef
3 fmeduj oadn w
2 g agzm
6 fi zis qj snp myjm uz
7 zso nswkr s i soxc rpqui gl
2 hgneqg mfcjo
3 cxlvf zr n
9 m yjhfx zhgs qi itifa tmqwbb grkj ubokfo g
6 xoi qibxka uhfm vpenwj v yz
1 crse
1 l
4 usyrpn ytwskd c ax
9 aeaji a salidm vuvfna gf wnrbs dhs bg hvftaf
2 guhjk cugues
4 l wfgg jach juljr
5 byjcpq etsljt tnud futht ecocte
8 e gplq engdkk neh ipzka badgjb qyhg cb
1 m
2 ssnu ayxaf
2 hc iqridy
3 jgxtc va v
2 zooaus h
6 chwizy f cn pdi tekmh hgpod
8 ukwb bnox wwhzrn uybzv midoj xnh yfwhdy i
9 emnsxb vewubb bw ryf pl cwnrtf n hau yi
1 gv
8 xpoya j qtuq nd oqvf fjtew dqfak zbtxti
7 tts ecshes wkb tb lb y pmtewb
10 nckfqq di kpx ioxzgk ataxt sdnnnt siny i a tu
6 gc ljkbd vxj ve ipq yrijs
6 ysron cmcq dk pyvs x moti
8 r rlun gtcafv waj v yjmw vn itmmi
1 fja
4 g kz sgy wkmsea
1 h
7 dmjeig bvnn cgo wr ompwit no up
5 vp aiqi kgs nbcdff jzr
4 gyucdq qok yeghd qpy
`

func expected(names []string) []string {
	cnt := make(map[string]int)
	res := make([]string, len(names))
	for i, n := range names {
		if cnt[n] == 0 {
			res[i] = "OK"
			cnt[n] = 1
		} else {
			res[i] = fmt.Sprintf("%s%d", n, cnt[n])
			cnt[n]++
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesCRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 1 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		var n int
		fmt.Sscan(parts[0], &n)
		if len(parts)-1 < n {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		names := parts[1 : 1+n]
		exp := expected(names)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d\n", n)
		for _, name := range names {
			fmt.Fprintf(&buf, "%s\n", name)
		}
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(outLines) != len(exp) {
			fmt.Printf("Test %d failed: expected %d lines got %d\n", idx, len(exp), len(outLines))
			os.Exit(1)
		}
		for i := range exp {
			if strings.TrimSpace(outLines[i]) != exp[i] {
				fmt.Printf("Test %d failed line %d: expected %s got %s\n", idx, i+1, exp[i], outLines[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
