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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(names []string) []string {
	seen := make(map[string]bool)
	res := make([]string, 0, len(names))
	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]
		if !seen[name] {
			seen[name] = true
			res = append(res, name)
		}
	}
	return res
}

const testcasesBRaw = `100
1 bob
2 frank carol
5 eve judy dave judy alice
10 carol grace grace ivan frank ivan heidi ivan eve alice
1 frank
8 frank grace grace ivan carol ivan carol dave
4 alice carol frank carol
3 ivan ivan frank
9 ivan carol heidi grace ivan frank judy frank frank
8 carol grace heidi ivan dave heidi eve heidi
9 ivan frank heidi heidi frank judy ivan heidi heidi
4 frank carol judy eve
8 eve eve ivan ivan ivan ivan judy judy
7 eve dave heidi ivan frank judy bob
6 alice dave bob alice judy alice
5 judy dave bob ivan carol
5 dave dave alice grace alice
1 frank
6 carol dave alice bob bob bob
1 alice
1 frank
5 carol carol carol ivan alice
7 judy alice dave carol alice alice frank
10 bob eve frank heidi alice eve heidi ivan judy alice
5 grace judy carol heidi dave
2 frank bob
1 heidi
3 ivan judy grace
8 ivan frank carol frank eve eve judy grace
1 ivan
3 alice eve alice
3 carol carol bob
8 dave ivan alice dave dave heidi bob eve
2 judy dave
10 judy frank eve grace eve ivan alice carol alice grace
7 carol bob ivan bob dave bob bob
1 carol
4 bob dave alice ivan
8 heidi eve ivan grace dave dave grace grace
9 alice judy judy alice grace ivan judy carol bob
8 frank alice ivan bob judy frank eve frank
5 alice grace bob bob eve
4 alice heidi alice grace
8 heidi dave judy judy bob alice eve alice
6 eve bob dave heidi dave bob
10 frank grace heidi carol frank grace bob eve bob bob
2 judy frank
7 dave bob alice judy heidi alice heidi
5 frank heidi carol frank eve
8 ivan heidi grace heidi eve grace dave carol
8 judy eve ivan grace bob judy judy bob
2 frank carol
9 carol grace bob bob alice carol eve grace dave
6 heidi carol ivan eve bob carol
9 grace bob frank ivan dave ivan eve carol carol
8 dave grace frank judy carol heidi heidi alice
10 grace carol grace ivan alice heidi eve grace eve grace
8 frank ivan frank bob dave ivan judy dave
7 grace alice frank heidi ivan heidi carol
2 alice grace
4 judy judy grace dave
2 grace ivan
4 eve judy judy dave
8 judy carol alice judy grace heidi eve ivan
10 carol heidi dave bob frank alice heidi ivan bob judy
8 frank heidi eve ivan heidi alice bob judy
6 carol grace eve carol alice carol
8 grace heidi eve carol alice eve ivan heidi
1 frank
1 ivan
7 judy heidi dave eve heidi carol heidi
9 eve bob eve frank eve frank eve grace ivan
2 ivan dave
7 judy ivan carol ivan bob eve alice
4 heidi ivan dave ivan
5 alice bob bob grace frank
4 frank frank bob frank
8 frank carol heidi heidi eve heidi carol heidi
4 eve frank carol bob
4 heidi dave frank carol
6 carol carol dave eve ivan grace
7 frank eve judy ivan judy frank grace
5 ivan judy bob frank eve
7 heidi carol eve frank heidi heidi bob
3 frank grace carol
1 bob
6 carol frank bob grace alice ivan
6 dave judy grace ivan eve heidi
3 frank frank dave
8 bob carol dave frank eve carol grace frank
5 bob frank dave dave dave
10 alice frank frank judy alice carol carol bob grace heidi
5 carol frank ivan judy bob
6 judy grace dave alice grace heidi
8 judy frank ivan judy judy bob judy ivan
9 heidi grace heidi carol grace grace ivan heidi alice
2 heidi judy
3 bob ivan carol
2 grace eve
8 alice eve bob frank dave carol alice carol`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "empty testcases file")
		os.Exit(1)
	}
	t, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test count")
		os.Exit(1)
	}
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing test case %d\n", caseIdx)
			os.Exit(1)
		}
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseIdx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil || n != len(fields)-1 {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseIdx)
			os.Exit(1)
		}
		names := fields[1:]
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(names, " "))
		want := strings.Join(expected(names), "\n")
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseIdx, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\n", caseIdx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
