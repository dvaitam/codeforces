package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesB = `xxyxxxyx
yyyyxyyyyyxyyxxxx
xyxxxyxyxyxy
xxyyyyyxxxyxxyxxyyyy
yyxxyxxyy
xyyyyxxxxxxyxyx
yxxxyxxyyyyxyx
yyyyxxxxxy
xyxyxxyy
yxyx
xyyxyyyx
xyyxx
xxxxxyxxxxyyx
xxyyyyxxxyyyxyxy
yyyxxxy
yyxyyxyx
xxyxyxxxxxxy
xyy
xyyxy
xyyyxyyyxx
yxxxxxyxxxxxxx
xyxxyyyxxyxyyx
xyxyyyyyxyyyxyx
yxyyyxxyxyyxyyy
xyxyxxy
x
xyyyyyyyyxxxxyxyxyy
xxyyxxxxxyyyyyxyyyyy
yxxxyyy
yx
y
xxyyxxxyxyxyxyxxx
xyyyyyyyxyxyxyx
yyyxyyxyyxyyx
xy
xxxxyxyxyyxy
xxyyxxxxyxxyxy
xxxyxyyyyyxyxxxxyxx
xyxyxyxxyxxxyyy
yyxxyxxxxyx
yx
yxxyxyxyyy
yyxxyyxyxy
xxyx
yxyyyy
xyy
x
yyx
xxyxxyyxxxxxxyxy
xxxyyxxyxxyyy
xyxxxyxxxxyxxxxxxyx
yyxxyxyxxyyxyy
yyxyxyyx
yyxyxyxyy
xxyx
yyxyyxxyyyxyxy
xyyxxyxyyyxy
yxyyyxyxxxyxxyyyxx
xxxxxxxxyyyyxy
xyyxyyxyy
yyyyxxxxxxxyx
xxxxyyyxyyyyx
xxxxyx
yyxxyxyxyxyyxyx
xyyyyyxxyxyxyxxy
yxyyxxyxyxxxx
xxxyyxyyyxxxxyyyxxy
yxxyyxyyxxyx
yyyxxyyyxyyyyyyy
yyxy
yxxyy
xyyyyxyy
xyxxyyyxyxxyxyxy
xxyxxx
y
xxxxxyxxyxxxxyyyxy
yxyxyyy
xxyxxyxxxyxyxyx
xyxyyxyyxyyyxyxyxy
yxyxyyxyxyyyxxyy
yxxxyxxxyxyx
yxyyxxyxyxyxyxy
xyxxyyxxxxxxyxxyyyy
xyyyyyy
xxxxxyyxyxxyxxy
yxxxyyyxxyxxyyyxxxx
yxyyxyyxyyyxyxxxyyx
yyyyxxyxxxxyxxyyy
yyyx
xyxxyxyyxyxyxyxxxy
xy
xxxy
yyxyyxx
yxyxxyyyxxxyxxy
yyxyyxx
yyxyxyyxxxyxyyyxxxy
xyyyxyyyyxyyxxxyyxxx
xxxyxyyxxy
yyyyxxyxxx
yxxyxxxy`

func solve(s string) string {
	countX := strings.Count(s, "x")
	countY := len(s) - countX
	if countX > countY {
		return strings.Repeat("x", countX-countY)
	}
	return strings.Repeat("y", countY-countX)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesB))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := line + "\n"
		want := solve(line)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
