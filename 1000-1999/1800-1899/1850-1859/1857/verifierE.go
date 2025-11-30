package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesB64 = "MTAwCjIKMSAxMAo1CjkgNCAxMCAzIDEKMwo0IDEgMwoxCjUKNAo2IDcgNyAxMAoyCjkgMQo2CjIgMiA5IDUgOSAxMAozCjYgNSA0CjUKNCA4IDcgMTAgMQozCjMgNyA5CjEKOQoxCjkKMwo2IDIgOQozCjkgMTAgNQo2CjggNyAxMCA4IDQgNgo0CjggMTAgMSA4CjEKNQoxCjcKNgo0IDQgNiA3IDYgOQo1CjEgNSA0IDUgMgo0CjYgMiA5IDMKNAo2IDYgNCAxMAozCjMgNCA4CjYKNCA0IDQgNSA0IDQKMgo0IDEwCjMKMyAyIDIKMQo4CjUKMSA5IDEgOCA2CjUKNiA4IDEwIDUgOAoyCjkgNQo0CjQgMTAgMyAxCjUKNyA4IDkgMSA5CjMKNyAzIDgKMQo5CjQKNiAxIDEgNAo0CjEwIDIgNiA2CjUKNSAzIDkgNSA0CjMKOCAxMCA2CjMKMyA4IDcKMgo3IDYKMwo3IDUgOAoyCjEwIDMKNQo0IDIgNSA0IDkKNAo3IDQgMyAzCjQKMSAzIDcgOAoyCjkgOQo2CjYgOCA1IDUgOSAxMAoxCjcKMQozCjEKNAozCjggMiA4CjYKNiA1IDYgNiA5IDkKMwo2IDEgNwoxCjYKNgoxMCA0IDEwIDEwIDkgMwozCjYgNiA5CjQKOCAyIDcgMTAKMwo4IDMgMwo2CjYgNyAxMCAxMCA3IDIKMQozCjMKNCA2IDYKNgo1IDUgMTAgNyA1IDgKMwoxMCA5IDgKNAozIDMgMSAzCjYKNCA1IDIgNCAzIDcKMQoyCjEKOAo1CjEgMSA3IDkgMgoxCjUKNQozIDIgNyA1IDQKMgo1IDgKNgozIDMgOSAzIDEgMQo2CjYgNiA4IDggOSA1CjUKOCAzIDkgMyAxMAo1CjggMTAgNyA5IDkKMwo2IDggOQoyCjIgOAozCjYgNyA1CjIKNSAxCjEKOQoyCjYgNQo0CjkgMSA4IDgKMwoxIDcgNQo1CjQgMSA5IDYgNAo0CjkgNyA3IDQKMgo3IDcKMgo1IDIKNAo0IDkgMTAgNQozCjkgNSA3CjIKMiAzCjIKNSAyCjYKOCAxMCA4IDUgMTAgOAoyCjYgMQoyCjEgOQoxCjMKMQo1CjMKMyA4IDEKMgoyIDMKMQo4CjIKNiA3Cg=="

func expected(xs []int64) string {
	n := len(xs)
	type pair struct {
		val int64
		idx int
	}
	arr := make([]pair, n)
	for i, v := range xs {
		arr[i] = pair{v, i}
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + arr[i].val
	}
	ans := make([]int64, n)
	total := int64(n)
	for i := 0; i < n; i++ {
		s := arr[i].val
		left := int64(i)*s - prefix[i]
		right := (prefix[n] - prefix[i+1]) - s*int64(n-i-1)
		ans[arr[i].idx] = total + left + right
	}
	out := make([]string, n)
	for i, v := range ans {
		out[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(out, " ")
}

func runCase(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	inputs, exps := loadCases()
	for idx, input := range inputs {
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", idx+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exps[idx] {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, exps[idx], got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}

func loadCases() ([]string, []string) {
	data, err := base64.StdEncoding.DecodeString(testcasesB64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode embedded testcases: %v\n", err)
		os.Exit(1)
	}
	tokens := strings.Fields(string(data))
	if len(tokens) == 0 {
		fmt.Fprintln(os.Stderr, "no testcases found")
		os.Exit(1)
	}
	t, err := strconv.Atoi(tokens[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid test count header\n")
		os.Exit(1)
	}
	pos := 1
	var inputs []string
	var exps []string
	for caseNum := 1; caseNum <= t; caseNum++ {
		if pos >= len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing n\n", caseNum)
			os.Exit(1)
		}
		n, errN := strconv.Atoi(tokens[pos])
		if errN != nil {
			fmt.Fprintf(os.Stderr, "invalid n on case %d\n", caseNum)
			os.Exit(1)
		}
		pos++
		if pos+n > len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing array\n", caseNum)
			os.Exit(1)
		}
		fields := tokens[pos : pos+n]
		pos += n
		arr := make([]int64, n)
		for i, f := range fields {
			val, err := strconv.ParseInt(f, 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid value on case %d\n", caseNum)
				os.Exit(1)
			}
			arr[i] = val
		}
		inputs = append(inputs, fmt.Sprintf("1\n%d\n%s\n", n, strings.Join(fields, " ")))
		exps = append(exps, expected(arr))
	}
	if pos != len(tokens) {
		fmt.Fprintf(os.Stderr, "unused tokens remaining in embedded tests\n")
		os.Exit(1)
	}
	return inputs, exps
}
