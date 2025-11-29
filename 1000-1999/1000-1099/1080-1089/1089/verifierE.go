package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesRaw = `697227012 9
982959855 9
403177920 9
878737245 9
881079065 9
460416173 9
978363084 9
346681174 9
408961000 9
299005664 9
738813869 9
222124952 9
724013126 9
919241793 9
665610161 9
33807794 8
706310842 9
346469750 9
642859667 9
650205371 9
316230288 9
895411994 9
628529510 9
589844461 9
543729272 9
999226194 9
255361447 9
651563129 9
193426091 9
600627831 9
263429130 9
824535979 9
230319767 9
436386572 9
986703366 9
404614466 9
354935560 9
665649671 9
437774560 9
327367270 9
598716121 9
10432410 8
643516779 9
946423095 9
150174608 9
993934468 9
94056335 8
898572954 9
619038123 9
697362138 9
824513616 9
639793139 9
106797913 9
720058116 9
722293077 9
51599289 8
622501035 9
700519128 9
227331733 9
665977071 9
107836370 9
289711637 9
676123637 9
494991357 9
128800528 9
117716233 9
156373449 9
549801260 9
397305252 9
348527799 9
936983421 9
776375732 9
535177241 9
51979753 8
660781339 9
331472949 9
940369759 9
129284083 9
407410349 9
383985637 9
922313472 9
880335068 9
515431917 9
712397790 9
225452305 9
537522492 9
777932676 9
123620866 9
237996995 9
723894732 9
682798934 9
775180093 9
953929263 9
53300460 8
774787165 9
815188609 9
795792054 9
326132763 9
155808098 9
233981289 9`

// referenceSolve mirrors 1089E.go.
func referenceSolve(num string) string {
	if len(num) > 0 && (num[0] == '+' || num[0] == '-') {
		num = num[1:]
	}
	return fmt.Sprintf("%d", len(num))
}

func runCase(bin, num, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(num + "\n")
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	caseNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Printf("Case %d: invalid test line\n", caseNum+1)
			os.Exit(1)
		}
		num := parts[0]
		expected := parts[1]
		caseNum++
		if err := runCase(bin, num, expected); err != nil {
			fmt.Printf("Case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d cases passed\n", caseNum)
}
