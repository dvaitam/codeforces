package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Cave struct {
	need int
	k    int
}

func expectedC(caves [][]int) int {
	arr := make([]Cave, len(caves))
	for i, monsters := range caves {
		maxNeed := 0
		for j, v := range monsters {
			if val := v - j; val > maxNeed {
				maxNeed = val
			}
		}
		arr[i] = Cave{need: maxNeed + 1, k: len(monsters)}
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].need == arr[j].need {
			return arr[i].k < arr[j].k
		}
		return arr[i].need < arr[j].need
	})
	start := arr[0].need
	cur := start + arr[0].k
	for i := 1; i < len(arr); i++ {
		if cur < arr[i].need {
			start += arr[i].need - cur
			cur = arr[i].need
		}
		cur += arr[i].k
	}
	return start
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseCase(line string) ([][]int, error) {
	parts := strings.Split(strings.TrimSpace(line), ";")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid format")
	}
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	if n != len(parts)-1 {
		return nil, fmt.Errorf("expected %d caves got %d", n, len(parts)-1)
	}
	caves := make([][]int, n)
	for i := 0; i < n; i++ {
		f := strings.Fields(strings.TrimSpace(parts[i+1]))
		if len(f) == 0 {
			return nil, fmt.Errorf("empty cave")
		}
		k, err := strconv.Atoi(f[0])
		if err != nil {
			return nil, err
		}
		if len(f)-1 != k {
			return nil, fmt.Errorf("cave %d expected %d monsters got %d", i+1, k, len(f)-1)
		}
		monsters := make([]int, k)
		for j := 0; j < k; j++ {
			v, _ := strconv.Atoi(f[j+1])
			monsters[j] = v
		}
		caves[i] = monsters
	}
	return caves, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		caves, err := parseCase(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", idx, err)
			os.Exit(1)
		}
		expected := expectedC(caves)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", len(caves)))
		for _, m := range caves {
			sb.WriteString(fmt.Sprintf("%d", len(m)))
			for _, v := range m {
				sb.WriteString(" ")
				sb.WriteString(fmt.Sprintf("%d", v))
			}
			sb.WriteByte('\n')
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if out != fmt.Sprintf("%d", expected) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx, expected, out)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
