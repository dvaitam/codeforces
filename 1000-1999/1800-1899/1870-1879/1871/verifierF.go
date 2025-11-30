package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesB64 = "MTAgLTQgMiAxMyAtMTkgOSAtNSAtMTcgLTEwIC0xMyAzCjggLTUgNCAxNCAtMTQgMTYgLTUgLTIwIC03CjcgLTMgLTkgNCAtMTAgLTE2IC0xMiAxOQoxMCA4IC0xMiAtMTIgLTIwIC0yMCAtNyAtNyAtMTAgLTEwIC0yCjYgLTggMTQgMjAgLTcgLTkgLTgKNyAtMSAtMTkgMyA2IC0xMCAtMTEgLTQKMiAxIC0xCjEwIDE3IC0yMCAxOCAxIC0xNiAtMSAyIC0xIDEwIDAKMyAxMCAxMCAtOQoxIC00CjEgMgo3IC0xOSAxNSA2IDMgNCAxNyAtMjAKOCAtMTggLTkgMTkgLTggLTEzIC01IDkgMgo5IDIgMTMgLTQgOSAtMTQgMTcgMyAtMiAtMTgKNyAtMTUgLTcgMSAxMiAxOSAzIC0xMQo2IC0zIDE0IC0xNSAtMSAwIC0xCjMgLTE1IDIwIC0xMQo1IDEwIC0xMCAtMTcgLTE1IDE4CjkgNSAtMTggLTUgMTggMiAtNCA5IDYgLTExCjEgMjAKMSAxMQo2IC03IC0xMiAxNiAtMTIgMjAgNgoyIC0xMCA3CjYgLTExIC0xNyA2IC0yIC0xMSA5CjEwIC0xMCAxMyA5IDExIDAgMTAgLTMgLTIgMTAgNQozIC0xMyA0IDE0CjMgMjAgMTEgMQozIC0xNSAxMSAtMwo5IDE1IDEyIDMgLTE2IDIgMTcgLTE4IC0xIDMKOSAtMyAxMSAtNCAtMiAxIC05IDE3IC0yMCAxMAo5IC00IDAgLTMgOSAtMiAxMiAyIDIgLTMKNiA2IDIgLTkgOCAzIDEKOSAtMTEgMTMgLTEwIC04IDMgMTAgLTIgLTE1IDYKMyAxOSAxNyAxMwo3IC0xIDE5IDE1IDIwIC0zIC0xOSAtOAozIDE3IDggMTkKMyAtNiAtOSAyMAoxIDEwCjQgLTEwIC0xNyAtMTIgLTEzCjYgLTkgMTAgLTggMTUgLTE4IDYKOCAyIDQgMTkgLTE2IDE3IC03IC01IDMKMSAyCjcgLTMgNiAtMTMgMTUgMyAtMTggMTUKMTAgLTEgLTE0IC0yIDE0IDEyIDEgMTcgLTIgMiAtMTIKNyA2IDE2IDE0IDMgOSAtMTEgLTEwCjEwIDQgMTYgMTAgLTggLTEyIDE4IC0xNSAyIC0yMCA0CjIgMCAxNgoxMCAxNCAtMTEgMCAyMCAxNiA0IDcgNyAtNiAxMQo1IDEwIDQgNCAtMTAgMTgKMTAgLTQgLTEgMTEgLTQgNiAtMTkgMCAtMSAxMSAtMgozIDEwIC0xOSAtMTMKMTAgOCAtNSAtMiAtMTggLTEyIDUgLTIwIDEwIDE0IDE1CjUgLTUgMTAgLTE4IC01IDExCjUgLTExIC0yIC0yIDExIDE4CjggMTMgMTggLTEzIC0xOSAtMTIgLTEgLTIgMTQKNiAxOSAtMiAxMyAtMTkgOSAyCjYgMTcgLTEyIC0xOCAtMjAgLTQgMTUKOCAtMTQgMTQgLTggLTIwIDcgNyAxOCAxNgo4IDQgMTAgNSAtOCAtMiA5IC0xNiAtMQoxIDcKMTAgLTIgMTAgLTEgLTExIC0xMCAxMCAxNSAxMSAxIDE0CjMgNyAxNyAxNAoxIC0xNgo0IC0zIC0xNSAtMTYgLTE5CjYgNyAtMTYgNSAxMSAtMTcgLTEzCjIgLTYgMTkKMiAtMTIgLTIKOCAtMTEgLTkgMTkgLTkgNiAtMTAgLTE2IDE5CjQgLTE4IDE1IC0xNCA0CjIgLTMgLTE3CjEwIDE2IC0xMyA1IDE5IC0xMiAtMjAgNyAtMTUgMCAxOAo4IDExIDIgMyAtMTcgLTEyIC0yIC0xMSAxNgo5IC0yIDE1IDE1IDE5IC02IC00IC0xNiAxNSAtNQo1IC0yIDEzIC0xMiAtNSAzCjggNCAtOSAtMTIgLTE5IDEgLTE1IDE2IC0xOAoyIC0xMyAxMgoxMCA5IC01IDQgOSAxMCAwIC0xNCAxMyAtMTkgMTQKNyAtMTcgLTExIDcgLTYgLTEzIC0xNSAxMQo0IC0xMiAxOSA0IDIKNCAtMiAxIDE5IDIKNyA0IC0xMiAyIC0yIDIwIDcgMwo5IC0xOCAxNyAxNiAtNyAtOSA1IC0xNiAtMTQgLTE4CjEgLTkKNCAtOCAtMTggMTEgMTAKNiAtMjAgNyAxMCAtMSAxOSA3CjYgOSA5IC0xNCAtOCAtMTEgLTEwCjIgMyA0CjggLTExIDE1IC00IC0xMyAtMyAtMTAgLTIgLTUKMSAxMAoxIDIKNiAwIC0xNyAtMTkgOSAxMCAtMTEKMiAwIC0yCjggLTUgLTEwIC0xOCAtOCAtMTkgMTYgLTYgLTE1CjEwIDUgMyAtMSAtOSA5IDMgLTIgLTE2IDkgLTEwCjQgLTkgLTcgLTE4IDE4CjcgLTMgLTIwIDkgLTE3IDggNiAtMTAKMSAtMTgKOSAxMyAxNiAyIC0xNCAtMTYgLTUgMTEgLTE1IDEwCjEgLTUKMSAxMQo="

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

func expected(nums []int) string {
	sort.Ints(nums)
	out := make([]string, len(nums))
	for i, v := range nums {
		out[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(out, " ")
}

func loadCases() ([]string, []string) {
	data, err := base64.StdEncoding.DecodeString(testcasesB64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode embedded testcases: %v\n", err)
		os.Exit(1)
	}
	fields := strings.Fields(string(data))
	var inputs []string
	var exps []string
	for idx := 0; idx < len(fields); {
		if idx >= len(fields) {
			break
		}
		n, errN := strconv.Atoi(fields[idx])
		if errN != nil {
			fmt.Fprintf(os.Stderr, "invalid n at position %d\n", idx)
			os.Exit(1)
		}
		idx++
		if idx+n > len(fields) {
			fmt.Fprintf(os.Stderr, "truncated case starting at position %d\n", idx)
			os.Exit(1)
		}
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[idx+i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid value in case starting at %d\n", idx)
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
		exps = append(exps, expected(append([]int(nil), nums...)))
	}
	return inputs, exps
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
		if strings.TrimSpace(out) != exps[idx] {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, exps[idx], out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
