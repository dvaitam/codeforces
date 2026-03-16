package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const testcasesBRaw = `
198790859
800786929
47718372
583637907
557229624
525065452
850547752
668248295
540812701
906193520
485359334
724170403
447338150
334965507
969551107
34219200
179967148
549839696
246531488
900252695
663645840
200458724
299052265
755558516
678138960
334914267
925627506
540308423
168532534
786205366
344977492
383278211
650555445
996876576
20875673
556045384
108771244
338500909
558661643
680474834
659728864
97287250
47180057
694322665
182482691
909134100
543047307
848372
649381097
550344379
162903427
568775824
646011784
903177526
627877148
258965783
784184711
370920616
719263116
658529383
651900717
990229140
969003519
473764659
443432681
387271075
455641971
856110679
950911111
754750231
698300824
921342238
932052465
952027539
252790172
549366369
166093414
464256513
708956841
813133591
695207613
819372513
288464212
717535705
788585823
257459789
441227004
278642578
744216714
738648749
467460767
588324918
906638018
200876443
799511090
108991679
264736135
269390170
136644630
461756056
`

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "177B1.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := line + "\n"
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
