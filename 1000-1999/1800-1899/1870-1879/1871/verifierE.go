package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB64 = "LTQwIC0yMwotNzQgODQKMSAyMgotNjEgLTc3Ci04MyAtOTUKMiA0MAotMjYgOTUKLTg1IC00NAozMyAzNwotOCAtMzAKOTkgLTU2Ci03MyAtMzMKLTQ2IC05NAo2NCAtMzQKLTMxIC01MQotNTggLTIxCi0yNiA2MAo4NyAtNQotNzggNTUKLTE0IDcxCi0xIDI5Ci0zNyAtNTUKLTM3IDIxCi0yOSAtNzgKNDAgLTI0Ci05OSAtMjYKNDYgODAKLTIxIDk1CjMwIC01MQo1IDgKNTMgLTI3CjEwIDE1Ci01OSAtNDEKLTIyIC0zNAotODkgLTgwCi04OSAxOAo2MCAtMjkKMzIgMzYKNjUgMjAKNzkgLTEzCi02MyA3MgotNTAgLTgzCjUgLTQ5CjYyIDYxCjEyIC0zMAotNTMgLTkKMTEgOTEKNTAgLTE4CjYyIDQyCi01MCAtMTgKLTc1IC04NQo4MSAtNDIKLTI5IDk1CjQ5IDU3Ci00MCAtNjkKLTE2IC01NQotMjYgMTcKLTk0IC05MAotOSA3OAotNzkgLTI3Cjg4IDcyCi0xNyAtOTYKLTE4IC0yNwotMTggLTYxCjk4IDY2CjUgNTgKNzQgLTgxCi0yNSA1OAotNTEgMTMKLTI2IC02NgotMzYgLTMKNTMgLTYwCi0xNiA0NgotOTggLTcKLTg5IDE2Ci01NyAtNwoxMDAgLTgKLTI2IDQ2Ci03NiAxMgotNDcgOAotNDcgLTcxCi04NSAtODUKLTg2IDg4Ci01NyA1Mgo3MyAtNjIKNTUgLTkwCjM5IDI1CjQ5IC0zNwotMTggLTkxCi02OSAzNQotMjYgOTgKNCA2NgotNDkgMjIKLTQ5IC0zOQoxMiA1CjI1IC05MQotNDQgNwoxMyAtMzcKNjUgOQotNDUgMjcK"

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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
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
			fmt.Fprintln(os.Stderr, "invalid numbers in embedded tests")
			os.Exit(1)
		}
		inputs = append(inputs, fmt.Sprintf("%d %d\n", a, b))
		exps = append(exps, gcd(a, b))
	}
	return inputs, exps
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
