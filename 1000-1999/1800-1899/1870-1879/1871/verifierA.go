package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB64 = "NzI5IC0yMTIKNTUyIDgyMwotMTM5IC05MTgKLTQ3MCA5NzcKNDcgLTUKLTE3MSA4ODAKNjA1IDY5OQotMzc5IDk4MgotMjQgLTI2NwoxOTQgODI2Cjg1OSAtNTUzCjMzIC03MTUKLTQyMyAtNzE0CjU0NyAtODA2CjI2NiA2MzcKLTQ4NyA4NjMKOTAgNDQ0CjY1OSAyMzIKODQ3IC03MDAKLTM2NSAtNzk4CjQ5NCAtODQ5Cjg0MCA3NDEKNDAwIC0zMjQKLTM0IDE0NgotNzk0IC0yNzYKLTExMSAtMzUzCjI1MSAzMTEKODY5IC01ODIKOTc5IDEzMQotMjQgLTk0Cjc3MiA2NwotNDY3IC04NzMKNjQ4IDg4MQoxMjMgODc1Ci05NzIgLTgwOQo0NzMgNzIwCi0xODQgNDU0CjY4OSA2MDcKMzY4IDI4MAotOTk4IDI1MwoxMCA2OTUKNzc2IC0zMTgKLTUwMSA0OTUKLTMzNCA0NDEKNzgyIC04NzIKLTYwOSA4NzgKMTYyIC01NDYKLTUxMiA2NDUKOTgxIC03MDkKNjQ0IDExMgotODMgLTgxNAotODM2IC0zNDUKNzkyIDQwCjkxMCAyCi03NzcgLTM4MwoxMjggLTQwNAo0NDcgLTc0NQoxMjEgLTMxOQo2NjggODg4CjEwNiAtNTg0Cjk3MyA2MzcKMjM1IDEyMAoyMDMgLTQxMQotODkgLTgxMwoyMjEgNjM0Ci0yMTIgLTM1MQoxNzggLTUwNQotNDA2IC02MjQKLTYxMyA2ODIKLTYxOCAtOTMzCjI1NCAzNDQKLTQ2OCAtMjUKLTg1OSAtODE3CjM5MCA1NTEKLTczNCA3OTUKLTY5NCA4OTEKLTkyMSA3MjUKLTgzNiA4MzkKNDMyIDg5MAo2OTggMTA3CjM5OSAtMTk5CjcxNSA0NDQKNzQgLTQzNgo2OCA2NjIKLTUxOCA3MzkKLTU2MCA4MzMKMzkxIDIwNwo2OTAgOTQ1Ci0xNDIgMTg3Ci00MzcgLTc4CjggMzUyCjMxMyA0MzQKODc3IDYyNAotMjY5IC04MzIKLTMzNiAyNTQKLTc2NCAtNAoyMDIgMjkwCi0zMTQgNzMwCi02MTEgLTUwMwotOTY3IDQ5OAo="

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func loadCases() ([]string, []int) {
	data, err := base64.StdEncoding.DecodeString(testcasesB64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode embedded testcases: %v\n", err)
		os.Exit(1)
	}
	fields := strings.Fields(string(data))
	if len(fields)%2 != 0 || len(fields) == 0 {
		fmt.Fprintln(os.Stderr, "invalid embedded testcases")
		os.Exit(1)
	}
	var inputs []string
	var exps []int
	for i := 0; i < len(fields); i += 2 {
		a, errA := strconv.Atoi(fields[i])
		b, errB := strconv.Atoi(fields[i+1])
		if errA != nil || errB != nil {
			fmt.Fprintln(os.Stderr, "invalid number in embedded tests")
			os.Exit(1)
		}
		inputs = append(inputs, fmt.Sprintf("%d %d\n", a, b))
		exps = append(exps, a+b)
	}
	return inputs, exps
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	inputs, exps := loadCases()
	for idx, input := range inputs {
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != fmt.Sprintf("%d", exps[idx]) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, exps[idx], out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
