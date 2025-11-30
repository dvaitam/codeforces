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

type point struct {
	x, y int
}

const testcasesRaw = `100
2
4 2
2 3
1 4
3 4
3
2 3
1 4
2 3
2 3
1 3
4 2
2
2 4
2 3
2 2
4 4
3
4 2
1 4
4 3
3 2
3 4
3 4
2
4 2
2 4
4 1
4 2
3
4 1
3 3
2 2
2 3
2 2
3 2
2
1 1
2 3
3 3
2 1
3
3 1
4 3
3 3
1 1
1 3
4 4
3
2 1
4 4
1 3
3 2
1 1
2 4
3
1 3
3 4
1 1
2 2
3 4
2 1
2
1 3
2 4
1 4
2 3
2
1 2
3 2
2 2
4 4
2
2 1
3 4
2 1
3 3
3
2 4
2 1
1 3
4 1
1 1
3 3
2
3 1
2 3
1 1
4 4
2
3 2
2 2
2 1
1 4
3
2 4
2 2
1 3
2 3
1 4
4 2
3
3 4
2 4
2 2
4 3
4 3
4 2
2
4 1
2 3
3 1
3 4
2
4 2
3 3
1 3
1 4
2
1 1
1 4
3 4
3 1
3
1 4
3 1
3 1
1 3
1 1
1 1
2
4 3
1 1
2 1
2 2
2
4 3
4 2
1 4
1 2
2
4 4
1 2
2 1
2 4
3
4 3
2 1
4 2
4 2
1 3
4 4
3
4 3
2 1
3 2
3 3
2 4
3 3
3
3 3
2 1
3 3
2 1
2 4
2 3
2
4 4
2 3
3 4
1 4
3
1 2
1 3
3 1
4 4
2 1
2 1
3
4 3
4 4
1 3
4 1
1 1
2 2
2
3 4
2 3
4 2
1 3
2
2 2
2 1
1 4
3 3
2
4 4
2 1
2 1
4 3
2
3 1
4 2
3 4
1 4
3
4 2
2 1
1 2
4 3
3 2
1 3
3
1 2
1 4
1 4
3 3
4 2
1 2
2
4 2
3 1
1 2
1 3
3
4 3
3 3
3 2
1 4
1 1
2 2
3
4 2
3 3
3 3
1 1
4 1
4 4
2
1 4
2 4
2 4
1 3
3
3 3
2 1
4 3
2 1
2 4
4 2
2
2 3
2 3
1 3
3 4
2
1 2
3 2
4 3
1 1
3
4 1
2 1
3 4
2 3
4 3
2 1
3
2 1
4 4
3 3
2 2
1 4
2 2
3
2 2
4 2
1 1
2 4
4 4
1 2
3
3 1
2 3
4 4
1 3
1 4
3 2
2
1 3
4 2
3 3
1 1
2
4 2
4 2
4 4
4 2
2
1 4
1 4
3 2
1 3
3
2 4
4 1
2 3
3 3
4 1
4 1
3
3 4
3 4
1 3
2 4
4 3
4 2
3
3 1
3 2
2 3
1 2
1 4
4 3
2
1 4
3 3
4 1
1 1
2
4 1
1 4
4 2
1 3
2
2 4
1 1
1 3
2 4
2
2 2
2 1
1 4
4 3
2
1 1
1 1
2 4
3 1
3
3 1
2 1
3 3
2 2
3 2
2 2
3
1 1
3 4
4 3
2 1
3 4
1 2
2
4 1
3 3
2 1
1 1
2
2 3
1 1
3 3
4 1
2
4 2
3 3
1 2
1 3
3
3 3
1 1
4 3
3 3
1 2
4 3
2
2 1
4 3
1 2
2 2
2
3 4
1 3
2 4
2 4
2
4 2
3 3
3 1
4 1
3
4 4
3 2
2 3
4 4
3 2
2 1
3
2 2
2 1
3 4
1 1
1 3
2 2
3
4 4
3 3
3 3
4 2
4 2
4 2
3
3 3
1 1
2 3
3 2
1 3
3 1
2
4 2
1 4
4 1
2 1
2
2 2
3 3
1 3
3 1
3
1 3
2 4
1 4
2 4
4 3
1 3
2
1 4
4 1
2 1
1 4
2
4 2
4 3
4 1
1 2
2
2 4
2 4
2 1
1 4
3
4 3
4 3
1 2
1 4
4 4
3 4
2
1 1
3 4
4 2
1 2
3
2 2
3 2
4 1
3 4
2 1
3 4
2
4 4
1 3
4 3
3 4
3
4 3
1 2
4 3
2 1
3 1
4 4
2
4 2
4 3
1 1
2 3
2
2 3
1 2
2 2
4 3
3
2 2
4 1
2 4
1 1
3 4
3 2
2
3 3
2 2
3 4
3 3
3
4 4
2 4
2 3
4 4
4 1
3 4
3
1 3
4 3
4 1
1 3
3 4
3 4
2
3 4
4 1
4 3
1 3
2
4 4
2 1
3 2
2 1
3
1 1
3 3
3 1
4 2
1 4
3 4
3
2 3
2 2
2 4
2 2
3 3
2 3
2
1 4
2 2
2 4
2 4
3
1 1
3 2
1 1
3 1
4 3
3 4
3
3 2
3 1
3 4
3 1
2 4
4 3
3
2 2
1 2
1 3
1 2
3 4
3 1
3
1 2
4 1
4 4
2 1
1 3
4 3
2
4 4
2 3
4 3
1 1
3
3 4
2 3
1 1
4 2
4 3
2 2`

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func hasRectangle(points []point) bool {
	colRows := make(map[int]map[int]struct{})
	for _, p := range points {
		if colRows[p.y] == nil {
			colRows[p.y] = make(map[int]struct{})
		}
		colRows[p.y][p.x] = struct{}{}
	}
	type pair struct{ a, b int }
	counts := make(map[pair]int)
	for _, rowSet := range colRows {
		rows := make([]int, 0, len(rowSet))
		for r := range rowSet {
			rows = append(rows, r)
		}
		for i := 0; i < len(rows); i++ {
			for j := i + 1; j < len(rows); j++ {
				a, b := rows[i], rows[j]
				if a > b {
					a, b = b, a
				}
				key := pair{a, b}
				counts[key]++
				if counts[key] >= 2 {
					return true
				}
			}
		}
	}
	return false
}

func validateYes(points []point, out string) bool {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 5 {
		return false
	}
	if strings.ToLower(strings.TrimSpace(lines[0])) != "yes" {
		return false
	}
	inSet := make(map[point]struct{})
	for _, p := range points {
		inSet[p] = struct{}{}
	}
	pts := make([]point, 4)
	for i := 0; i < 4; i++ {
		fields := strings.Fields(lines[i+1])
		if len(fields) != 2 {
			return false
		}
		x, err1 := strconv.Atoi(fields[0])
		y, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return false
		}
		pts[i] = point{x, y}
		if _, ok := inSet[pts[i]]; !ok {
			return false
		}
	}
	xs := make(map[int]struct{})
	ys := make(map[int]struct{})
	setCheck := make(map[point]struct{})
	for _, p := range pts {
		xs[p.x] = struct{}{}
		ys[p.y] = struct{}{}
		setCheck[p] = struct{}{}
	}
	if len(xs) != 2 || len(ys) != 2 || len(setCheck) != 4 {
		return false
	}
	var xv, yv []int
	for x := range xs {
		xv = append(xv, x)
	}
	for y := range ys {
		yv = append(yv, y)
	}
	for _, x := range xv {
		for _, y := range yv {
			if _, ok := inSet[point{x, y}]; !ok {
				return false
			}
		}
	}
	return true
}

func loadCases() ([]string, []bool) {
	tokens := strings.Fields(testcasesRaw)
	if len(tokens) == 0 {
		fmt.Println("no embedded testcases")
		os.Exit(1)
	}
	t, err := strconv.Atoi(tokens[0])
	if err != nil {
		fmt.Println("invalid testcase count")
		os.Exit(1)
	}
	var inputs []string
	var answers []bool
	pos := 1
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos >= len(tokens) {
			fmt.Printf("case %d incomplete\n", caseIdx+1)
			os.Exit(1)
		}
		n, errN := strconv.Atoi(tokens[pos])
		if errN != nil {
			fmt.Printf("case %d invalid n\n", caseIdx+1)
			os.Exit(1)
		}
		pos++
		points := make([]point, 0, 2*n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i := 0; i < 2*n; i++ {
			if pos+1 >= len(tokens) {
				fmt.Printf("case %d missing points\n", caseIdx+1)
				os.Exit(1)
			}
			x, _ := strconv.Atoi(tokens[pos])
			y, _ := strconv.Atoi(tokens[pos+1])
			pos += 2
			points = append(points, point{x, y})
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
		inputs = append(inputs, sb.String())
		answers = append(answers, hasRectangle(points))
	}
	return inputs, answers
}

func parsePoints(input string) []point {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Scan() // t
	sc.Scan() // n
	n, _ := strconv.Atoi(strings.TrimSpace(sc.Text()))
	points := make([]point, 0, 2*n)
	for i := 0; i < 2*n; i++ {
		if !sc.Scan() {
			break
		}
		fields := strings.Fields(sc.Text())
		x, _ := strconv.Atoi(fields[0])
		y, _ := strconv.Atoi(fields[1])
		points = append(points, point{x, y})
	}
	return points
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	inputs, exists := loadCases()
	for idx, input := range inputs {
		got, err := run(cand, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", idx+1, err, got)
			os.Exit(1)
		}
		hasRect := exists[idx]
		first := strings.ToLower(strings.TrimSpace(strings.Split(got, "\n")[0]))
		points := parsePoints(input)
		if hasRect {
			if first != "yes" || !validateYes(points, got) {
				fmt.Printf("case %d failed validation (rectangle exists)\n", idx+1)
				os.Exit(1)
			}
		} else {
			if first != "no" {
				fmt.Printf("case %d expected No but got %s\n", idx+1, got)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
