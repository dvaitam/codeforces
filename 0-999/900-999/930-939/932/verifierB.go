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

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "932B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func buildInput(fields []string) (string, error) {
	if len(fields) != 3 {
		return "", fmt.Errorf("expected 3 fields")
	}
	return fmt.Sprintf("1\n%s %s %s\n", fields[0], fields[1], fields[2]), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	const testcasesRaw = `140892 141474 2
267460 267580 8
797927 798387 8
683245 683633 4
98419 98918 1
936711 937566 7
453790 454412 1
729634 730090 5
756590 757411 4
619870 620837 2
945216 945541 1
23407 23433 9
9653 10614 7
719831 720052 7
761112 761141 9
232461 233243 8
984788 985295 9
244407 244760 4
709728 709952 8
998501 998797 1
436397 437254 9
966985 967642 2
194937 195581 5
126763 127523 6
939079 939817 9
981930 982920 7
532381 533230 4
318105 318395 8
887303 888266 9
412462 413065 1
503555 503803 7
434440 435120 3
384958 385519 6
90668 91117 9
113175 113972 3
546244 547104 7
388522 389023 1
492118 492162 5
737550 738418 7
678593 678767 3
526636 526868 1
807953 808157 9
964781 965661 9
243455 243869 9
360528 361503 6
481435 482366 5
691237 691798 1
402328 403130 9
848445 848577 9
815161 815735 4
446789 447761 1
504472 505362 6
597688 598255 4
986725 987241 7
508481 509313 6
434556 434910 1
564636 565189 6
480402 481016 1
843653 843888 3
577510 578108 3
902834 902927 9
835818 836689 5
34036 34897 2
87278 88166 1
475004 475018 5
261682 261957 2
836017 836656 3
361154 361451 2
175606 175769 5
552999 553973 3
688555 688834 5
476790 477509 6
520612 521097 2
24783 25102 7
360021 360452 4
270974 271085 5
943529 944276 9
219248 220236 7
856729 857727 1
236322 236340 7
153577 153613 3
467318 468039 9
711119 711555 9
872673 872898 9
472746 472974 9
680009 680040 7
707687 708276 6
691876 692522 7
61641 62396 5
131789 132780 4
918065 918113 5
74163 75042 2
325440 326379 5
779975 780137 7
592384 592642 3
8893 9467 1
619273 620112 4
944571 945154 8
179849 180696 9
39242 39629 4`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		input, err := buildInput(fields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad testcase %d: %v\n", idx, err)
			os.Exit(1)
		}
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
