package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test cases from testcasesA.txt so the verifier is self-contained.
const testcasesA = `906691060 413654000
813847340 955892129
451585302 43469774
278009743 548977049
521760890 434794719
985946605 841597327
891047769 325679555
511742082 384452588
626401696 957413343
975078789 234551095
541903390 149544007
302621085 150050892
811538591 101823754
663968656 858351977
268979134 976832603
571835845 757172937
869964136 646287426
968693315 157798603
333018423 106046332
783650879 79180333
965120264 913189318
734422155 354546568
506959382 601095368
108127102 379880546
466188457 339513622
655934895 687649392
980338160 219556307
593267778 512185346
475338373 929119464
559799207 279701489
66872193 864392047
986194170 589161386
983541587 15077163
100149904 772777020
902041077 428233517
762628806 885670548
842938613 717424033
671374074 1227090
657019496 529975200
889126175 931581387
34491750 386820475
401867228 102891201
587077106 47746496
104685406 978763388
193023735 470796995
972572539 618479891
531482835 682626329
600020187 72040750
95854964 800565216
679316178 871327852
718951094 782510795
828424444 775234180
701117697 4104918
196033582 966593061
495842095 21942293
722049372 778100160
802771221 466894553
876499452 258470064
942022849 666057940
35578620 201996856
83414012 48874221
992636443 432216920
111196190 118678630
266118304 61444477
380100344 906502564
30181042 305883793
858250966 122180542
648909229 64034765
821155098 38471336
787908172 941149454
363172538 844504305
827272144 177207796
796372459 440742265
212035350 729540821
902901335 450379132
179806677 123780794
701449765 101290178
293004262 149270370
555038788 690806782
97881285 14428517
813645548 183383354
679378563 72344282
11292702 331308228
284472308 155825666
828716417 505047819
56446165 872037836
531996139 760007923
648776755 885984575
688070335 748766966
854109977 315443086
279048042 850780821
203092264 444205039
809033952 975171301
870946995 127998102
9899161 976658678
361844724 784567614
613519526 497531023
706988060 826472841
199923756 77841213
448303733 493843093
246285456 414788462
105550069 873254131
725984019 680671478
680275933 524401183
248156797 192555430`

// Embedded solution from 226A.go.
func modPow(base, exp, mod int64) int64 {
	result := int64(1)
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			result = (result * base) % mod
		}
		base = (base * base) % mod
		exp >>= 1
	}
	return result
}

func expected(n, m int64) int64 {
	if m == 1 {
		return 0
	}
	ans := modPow(3, n, m) - 1
	if ans < 0 {
		ans += m
	}
	return ans
}

type testCase struct {
	n int64
	m int64
}

func parseCases() ([]testCase, error) {
	data := strings.TrimSpace(testcasesA)
	if data == "" {
		return nil, fmt.Errorf("no testcases found")
	}
	lines := strings.Split(data, "\n")
	cases := make([]testCase, 0, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d malformed", i+1)
		}
		n, err1 := strconv.ParseInt(fields[0], 10, 64)
		m, err2 := strconv.ParseInt(fields[1], 10, 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d parse error", i+1)
		}
		cases = append(cases, testCase{n: n, m: m})
	}
	return cases, nil
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.n, tc.m)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	exp := expected(tc.n, tc.m)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d m=%d)\n", idx+1, err, tc.n, tc.m)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
