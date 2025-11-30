package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB64 = "LTY2IDQ1Cjk1IC04NAotMzUgLTcwCjI2IDk0CjE1IDIwCjY2IC0zCi00NyAtNzYKMjQgLTkzCi0xIDEwCjU1IDk1Cjk2IC0xMDAKNzggMTQKLTMyIDg0Ci00MiA1MQotNzQgLTE5Ci05MyAtOTUKLTk0IDY2CjM4IC05OAotMyA3NQotNDUgOAo4NSAtOTMKMzUgLTQ0Cjk1IDEyCjI2IDQxCi00MSAtMTIKLTQxIDczCi00NCA5NAoxNyAtMjYKLTk1IDYKNDIgNjQKLTc1IC01Mwo2MSA4NQotMjUgLTcwCjkwIC0xNQo4NCA4MgoyOCA4CjI5IDcxCi01MiAtMjMKLTI4IDUwCjI3IDI5CjAgNTAKLTkyIDIyCi0zOCA5MAozIDYKNzAgLTU2Ci03IDQwCjc5IDk4CjcyIDg4Ci01IC03OAoxMiA2OQozMCAtNzMKOTkgLTU5CjMzIDAKLTYgMjUKODcgLTkzCjIwIC04OQotMjIgODAKNTcgNTEKNDggMAo2NSAtNTcKLTU3IDI4Ci00MiAtOTcKOTcgLTQ5CjM4IDQwCi00MSAzCjMxIC0xMgo0NyAtMTAKMTcgLTMyCjY4IDQwCjU1IDg2Ci05OSAtMgoxMDAgODkKMzEgLTY3CjMyIDk5CjQzIC00OAo5IC04NgoyMyAtNwo0NSA0MQotNDkgMjkKNSAyNAotOSA2Ci0xMiAtMTAwCjM3IDM4CjU5IDU2Ci0xNiAxNwo1MyAtOTMKLTQyIDYyCi01NSA0MAo0OSAtNTQKLTc3IDQxCi0zNSAtOTIKNzIgLTgyCi03OSAtOTYKMTUgLTk3CjkzIDkzCi0yOSAtMzcKLTMyIC03Mgo1OSAtNTMKLTEyIC0yNgotODMgLTU4Cg=="

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
		exps = append(exps, a*b)
	}
	return inputs, exps
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
