package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB64 = "NCAyNSAxOSAtMzQgLTMKMTAgMTAgMzAgMjQgLTQyIDI3IC00OSAxMCAtMTcgMjAgLTIxCjQgNDEgMTAgMTkgMjAKOCAwIDMxIC0zMSAtMjEgMzEgLTMxIDE2IC0xCjEgMzUKMiAtMzAgNDcKMTAgLTQ1IC0xMiA0OSAtNDcgLTE2IDEwIDI2IDQyIC0xIDQxCjcgMCA0MyAyMyA2IC0zMyAtNCAtMzgKMSAtMzMKOCAtMjMgLTE3IDM2IDUgNDkgMzAgLTEyIDMKOSAtMSAyMyAtNiAxOCAyNCAyIDI0IC0yMSAtNwoxIC0xNQoxMCAzNSAzOSAtMzAgMzkgLTkgMTkgMjMgMjIgLTM3IDQxCjQgMzEgMjMgLTE2IC0xNAoyIC00MiAxMQo4IC0zOSAtNiAtNDIgMiAtMzEgLTQ4IC0xMyA0CjcgLTM1IC00NSAyNyAyOCA0NyAtNDUgLTIKMTAgLTggMjAgLTE1IDE0IC0yMCAtNDYgLTExIC01MCAtNDEgLTM3CjEwIDE4IC00NiAtMjUgMiAtMTMgMjggLTE3IC0zMSAzOCAtNDUKNiAtMTAgLTQgLTMzIC0yIC0yIDgKOSAtMSAzMiAyNiAzNyAyMSAtMzcgMjkgMTQgLTE2CjcgMzEgNDIgNDEgLTIwIC0xMiA1IC0xNwo5IC0xMiAyMCAtNyAtNDkgNTAgMyAyNCAtMTAgLTQ4CjcgMjggMjUgMzAgLTMzIC00MyAzMSAzMAo2IDkgLTUgMzYgLTUgMjcgNDAKNSA0NCAxMiAtNDggMjUgLTQzCjEgLTMKNSAzMCA4IC0xMiAyNSAyNgo2IC0yOCAtNCAtMjcgLTEwIDQ3IC0zCjEwIC0xNyAtMTIgNTAgLTIgLTM3IDQ4IC00NyAyMiAzNyA0NAozIC0xMSAxNCAtMjIKNSAtMjAgLTkgLTI3IDM2IDUKMiAtMzcgMjYKNiAtOCAzNiAtMjIgNiAtMjkgLTQwCjYgNDQgMzMgLTIzIDIyIDcgLTE2CjQgNTAgLTM1IC00NiAxNwo0IC0xMCAyMyAtMjcgLTE1CjYgMzIgLTQwIDI5IC02IDI1IC0zNAo3IC0xMyAxNiAtMTYgOSAtNiAzMSAzCjUgMyAyMiAyIC00NiAyCjMgLTI1IC01MCAxMQoxMCAxNSA1IDIxIDQxIC0yMiAtNDYgNDUgOCA0NiAzNAo5IC0xNCAxOSAtNyAtMjEgLTQyIDI1IC0xNCAtMzUgLTE5CjEgLTQ2CjkgLTI1IDUgMjMgLTQ0IC00OSAxMSA0NSAtMzUgLTI5CjkgLTEyIC0yMCAzNCAtNDggMTcgMTggMiAtNDQgMjgKMiAtNyAtMzQKNSAxOSAxMSA1MCAtNDMgLTUKNCAtMjUgLTM1IDE4IC0zNQozIC0yMCAtMTUgLTM0CjEgMTIKMTAgMSAtNDQgNDYgLTE2IC0xOSAtMTYgMjkgMTcgMTYgNAoxIDEwCjYgNDkgLTUwIC00MyA0OSAtMzQgLTQ1CjIgLTQ0IC00Mgo4IC00NiA0MSAtMzkgMTUgMTQgMTIgLTEwIC0zMAo2IC00MSAtNiAtMSAzMiAtMSAyNQo1IC00IC0xNyAtMjYgLTggNAoyIC0zNCAyMQoxIDQxCjcgLTQwIDIyIC0yOCAtNDUgLTMgOCAyNwo5IC0yIDMxIC00NSAyOSA1IC00NCAtMyAzMCAxMwo2IDMgMzggMyA4IC00OCAtMTkKNCAxOCAtMTYgMzggMjUKMiA0IC0yMgo3IC0zNCAtNDcgLTkgLTMgMjEgLTE3IC0zNQo4IDM4IC0zNSA0MyAzNCAxNyAtMiAzNSAtMzcKNiAyMiAxOCAtMzcgMjUgNDEgLTUwCjggLTMyIC0yMCA0OSAtMSAtNDUgMTcgLTM5IDIyCjIgMzQgLTIKMyAtNDcgLTcgLTM1CjEgLTM2CjggMzkgLTE0IDI0IC0xMiAtMzkgLTQ2IDQ4IDIyCjkgMTcgNDEgLTIwIC0zNyAyMCA0NSAtMzggMjAgLTQzCjkgLTkgMjIgLTI3IC00MSAtMjAgLTI3IDMyIC0xOSA4CjEwIDM5IDQ2IDAgLTE4IC0zIDI2IDAgLTYgMjEgMwoyIC0yIDE0CjQgMiA0NSAtMzAgMwoxMCA0NiAyNCAzNiAxNiAzNyAxMSAtMzEgMzIgMSAtMzEKMyAtMzggMTMgNDUKOCAzOSAxNiA2IDI1IDQyIC0yNyAtMzMgLTE2CjQgLTMyIDI0IDE1IC0xMAo0IDM4IDE4IDQ5IC0xMwo3IDI2IDI0IDI0IC0xNiAtMjMgLTExIC00OAo1IDExIC0yIC0yNSAtMjggMjIKNiAtMjAgLTkgMTEgNDkgLTMyIDMKOCAzOSAyNiAtMjQgOSAyNCAzMyAyMSAtNDcKOCA0MiAtNDEgMSA1MCA0MyAtNDUgOSAtMjEKNCAzMiA0MSA0OSAzNgoyIC0yMyAtMTgKNCAtMjYgNDkgLTE3IC0zMwozIDI5IDQwIDM2CjEgLTE4CjMgLTQ1IC0xMCAtMjcKNyAtMzkgNDMgLTQwIC0zNSAtMzkgLTE3IC0xMwoxIC01CjggMjQgNDMgMzYgLTcgLTUwIC00NyAtOCAtOAo3IC0yIDEyIC00MSAtMjQgMzIgMjQgNDUKOCAwIC0zNCAxOSAtMTAgLTM1IC0xNSAtNDEgMzUKNyAtMzYgNiAxNyAtMTggLTM4IDE3IDM5Cg=="

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

func expected(nums []int) int {
	max := nums[0]
	for _, v := range nums[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func loadCases() ([]string, []int) {
	data, err := base64.StdEncoding.DecodeString(testcasesB64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode embedded testcases: %v\n", err)
		os.Exit(1)
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		fmt.Fprintln(os.Stderr, "no embedded testcases")
		os.Exit(1)
	}
	var inputs []string
	var exps []int
	for idx := 0; idx < len(fields); {
		n, errN := strconv.Atoi(fields[idx])
		if errN != nil {
			fmt.Fprintf(os.Stderr, "invalid n at position %d\n", idx)
			os.Exit(1)
		}
		idx++
		if idx+n > len(fields) {
			fmt.Fprintf(os.Stderr, "truncated case starting before position %d\n", idx)
			os.Exit(1)
		}
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[idx+i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid value in case starting before position %d\n", idx)
				os.Exit(1)
			}
			nums[i] = val
		}
		idx += n
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		inputs = append(inputs, sb.String())
		exps = append(exps, expected(nums))
	}
	return inputs, exps
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
