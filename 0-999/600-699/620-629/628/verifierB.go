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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(s string) string {
	var ans int64
	for i := 0; i < len(s); i++ {
		d := int(s[i] - '0')
		if d%4 == 0 {
			ans++
		}
		if i > 0 {
			prev := int(s[i-1]-'0')*10 + d
			if prev%4 == 0 {
				ans += int64(i)
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

const testcasesBRaw = `100
584197
076984
642807150842
75945992
6610935233
696069602714278
8900754706381206
50300891319344
176104
142851240003485
909776582369
0224555159
04
294568
41730
2814654611
775517176045229611
330
0168847793615
49263511
8
317643039213765
219729668757738930
55082492694
1180132040752275
68809189163489676
93002489451744666022
45007627
12560976701720099251
536710979519426418
0675375
007
089933188
84126961161162
7
60754151150552
3
192303104509
22717542
3
574485291189
2622358332
5602906112
48629645137580
6065735471214
8927
2662375825717
4079703
1489267173
640140046896714543
10144851832164482983
168644759221166972
457637785707427093
5
6081160501904432
4316756256627285232
566763373906031
502394
18456708869774735400
25
4
2
639
8373591915587
40803513853
4448647
304810777067
711131263791686
23
3245561849834789
844301912633408
60164195384331845
57492208
55902622812379332
65992719
8
57740382778514296354
0305357345240
980250700300
51056237
25425660648870
1662082982153230282
691
72990456007154273852
5476048487089
407616570040494388
64319539852
50909
55445
6962029072507493
864
821291
89664440644888553
20829657086930
823570393426
973
813271664465951
0704366660
7592945067820195501
1870505`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data := []byte(testcasesBRaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 1; i <= t; i++ {
		if !scan.Scan() {
			fmt.Printf("missing string for case %d\n", i)
			os.Exit(1)
		}
		s := scan.Text()
		input := fmt.Sprintf("%s\n", s)
		expect := expected(s)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", i, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
