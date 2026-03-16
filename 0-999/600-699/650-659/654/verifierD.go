package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "654D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("oracle build failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func lineToInput(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return "", fmt.Errorf("expected two integers")
	}
	if _, err := strconv.Atoi(fields[0]); err != nil {
		return "", err
	}
	if _, err := strconv.Atoi(fields[1]); err != nil {
		return "", err
	}
	return fields[0] + " " + fields[1] + "\n", nil
}

func run(bin string, input string) (string, string, error) {
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
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesDRaw = `243 606
557 133
378 937
618 485
640 594
67 620
13 930
857 480
265 564
239 196
734 481
553 856
562 487
406 654
881 154
237 650
155 888
948 535
399 759
15 687
795 65
163 776
980 605
43 308
798 31
843 886
275 484
609 736
942 899
396 731
807 943
437 404
745 820
590 455
987 958
137 899
374 99
36 139
506 222
264 988
688 446
797 641
875 308
431 519
853 395
587 359
546 599
417 598
237 925
344 698
937 951
29 876
286 620
687 712
167 715
881 334
987 554
926 585
582 106
730 671
216 648
851 587
273 291
127 64
493 874
654 495
90 352
819 68
420 918
154 20
300 437
787 425
893 121
45 619
629 779
46 386
735 600
338 564
902 944
285 517
241 36
317 7
78 110
614 548
32 971
202 994
417 298
625 269
159 706
43 888
347 321
368 981
141 918
882 386
385 471
890 532
395 659
887 609
697 572
105 635`

	scanner := bufio.NewScanner(strings.NewReader(testcasesDRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input, err := lineToInput(line)
		if err != nil {
			fmt.Printf("test %d: parse error: %v\n", idx, err)
			os.Exit(1)
		}
		expOut, expErr, err := run(oracle, input)
		if err != nil {
			fmt.Printf("oracle run failed on test %d: %v\nstderr: %s\n", idx, err, expErr)
			os.Exit(1)
		}
		gotOut, gotErr, err := run(target, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, gotErr)
			os.Exit(1)
		}
		if strings.TrimSpace(gotOut) != strings.TrimSpace(expOut) {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, strings.TrimSpace(expOut), strings.TrimSpace(gotOut))
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
