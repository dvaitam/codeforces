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

const testcases = `
100
906691060 413654000
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
357701129 261897307
784130655 349185523
755530427 934661371
67628852 205156724
984641620 609360020
238052748 256211902
862585180 153002189
862407392 583031025
481003666 97942385
86378037 343656009
939617817 545397110
525367681 117099973
323676027 591918702
312556225 758664544
134014483 587810203
357288143 874527132
990258111 580125105
218186328 858378322
647665637 587583850
630949022 308869636
477803330 98389212
640258140 856776201
413284458 340426464
618100573 259960579
311738932 197427542
203357386 882043706
200499310 35403864
657960186 705082654
279233229 511671264
74179728 96448169
728774303 813471023
139827435 941425019
160578447 991472819
41491066 904584779
86165975 964406046
750892139 991152223
890519443 580464751
733900803 420150929
899650944 757292285
563256981 295959888
560267850 871479690
252868060 912128613
231070626 961040770
729580034 633294197
886119717 450352273
622442781 295505361
483788449 528984690
708933077 688479847
751861428 984558133
851826327 383720487
88447325 348241150
657970850 123855848
522315486 630366660
676615554 359993934
907395137 204417732
260957516 17404312
785430573 291024800
125771989 757345667
236717702 399496591
853176975 183053409

`

func referenceSolve(n, k int64) string {
	return strconv.FormatInt((n/k+1)*k, 10)
}

type testCase struct {
	n int64
	k int64
}

func parseCases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcases))
	scan.Split(bufio.ScanWords)
	nextInt := func() (int64, error) {
		if !scan.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.ParseInt(scan.Text(), 10, 64)
		if err != nil {
			return 0, err
		}
		return v, nil
	}

	tInt, err := nextInt()
	if err != nil {
		return nil, fmt.Errorf("read test count: %w", err)
	}
	t := int(tInt)
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read n: %w", i+1, err)
		}
		k, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read k: %w", i+1, err)
		}
		cases = append(cases, testCase{n: n, k: k})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.k)
		expected := referenceSolve(tc.n, tc.k)
		out, stderr, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
