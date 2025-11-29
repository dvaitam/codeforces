package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const inf = int(1e9)

const testcases = `9161510903
173
008691413145
2087091
634579230225
419720769
456428071
084237
945992
466109352337
9606960
714
278789007547
6
8120
6503008
1319344217
1047142
85124000348
590977
65823694022
55515
0042294568
417
304281465461
87
55171760
52296
11133060168
477936153
92635
10
731764303
2137658219
29668757
38930555
8
492
9471180
32
4
7
227586
809189163
89676
9300248945
74
66602
23450076279
125609767017
009
2518536710
7951942641
306753751
0
40899331
88684126961
61
62
076607541511
055203
192303
10450932271
54203557
48529
18
4262235833
456
02906112648
2964513
58066065
35471214
89
276266237582
717340
797034148926
717386401400
68967
45
39101
48518
2164
82983
168644759
22116697284
763778
707427
9
0576
8
16
5
1
0443294316
56256627
285232756676
337390603125
2
9491
845670886977`

func distances(x, y int) [10][10]int {
	var dist [10][10]int
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			dist[i][j] = inf
		}
	}
	for start := 0; start < 10; start++ {
		q := []int{start}
		dist[start][start] = 0
		for head := 0; head < len(q); head++ {
			cur := q[head]
			step := dist[start][cur]
			nxt1 := (cur + x) % 10
			if dist[start][nxt1] > step+1 {
				dist[start][nxt1] = step + 1
				q = append(q, nxt1)
			}
			nxt2 := (cur + y) % 10
			if dist[start][nxt2] > step+1 {
				dist[start][nxt2] = step + 1
				q = append(q, nxt2)
			}
		}
	}
	return dist
}

func referenceSolve(s string) string {
	digits := []byte(s)
	var cnt [10][10]int64
	for i := 1; i < len(digits); i++ {
		a := digits[i-1] - '0'
		b := digits[i] - '0'
		cnt[a][b]++
	}
	var sb strings.Builder
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			dist := distances(x, y)
			total := int64(0)
			ok := true
			for a := 0; a < 10 && ok; a++ {
				for b := 0; b < 10; b++ {
					c := cnt[a][b]
					if c == 0 {
						continue
					}
					d := dist[a][b]
					if d == inf {
						ok = false
						break
					}
					total += int64(d-1) * c
				}
			}
			if !ok {
				sb.WriteString("-1")
			} else {
				sb.WriteString(strconv.FormatInt(total, 10))
			}
			if y < 9 {
				sb.WriteByte(' ')
			} else if x < 9 {
				sb.WriteByte('\n')
			}
		}
	}
	return sb.String()
}

func parseTestcases() ([]string, error) {
	lines := strings.Split(strings.TrimSpace(testcases), "\n")
	cases := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			cases = append(cases, ln)
		}
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := tc + "\n"
		want := referenceSolve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\noutput:\n%s", idx+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
