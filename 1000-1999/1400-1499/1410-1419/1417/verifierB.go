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

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
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

func lineToInput(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	tval, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", err
	}
	if len(fields) != 2+n {
		return "", fmt.Errorf("expected %d numbers got %d", 2+n, len(fields))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, tval)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fields[2+i])
	}
	sb.WriteByte('\n')
	return sb.String(), nil
}

func loadTests() ([]string, error) {
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var tests []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		tests = append(tests, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func parseLine(line string) (n int, T int64, a []int64, err error) {
    fields := strings.Fields(line)
    if len(fields) < 2 {
        return 0, 0, nil, fmt.Errorf("invalid line")
    }
    n64, err := strconv.ParseInt(fields[0], 10, 64)
    if err != nil {
        return 0, 0, nil, err
    }
    n = int(n64)
    t64, err := strconv.ParseInt(fields[1], 10, 64)
    if err != nil {
        return 0, 0, nil, err
    }
    T = t64
    if len(fields) != 2+n {
        return 0, 0, nil, fmt.Errorf("expected %d numbers got %d", 2+n, len(fields))
    }
    a = make([]int64, n)
    for i := 0; i < n; i++ {
        v, e := strconv.ParseInt(fields[2+i], 10, 64)
        if e != nil {
            return 0, 0, nil, e
        }
        a[i] = v
    }
    return
}

func parseLabels(out string, n int) ([]int, error) {
    fields := strings.Fields(out)
    if len(fields) != n {
        return nil, fmt.Errorf("expected %d labels got %d", n, len(fields))
    }
    res := make([]int, n)
    for i := 0; i < n; i++ {
        if fields[i] == "0" || fields[i] == "1" {
            if fields[i] == "1" {
                res[i] = 1
            } else {
                res[i] = 0
            }
        } else {
            return nil, fmt.Errorf("invalid label %q", fields[i])
        }
    }
    return res, nil
}

func checkValid(a []int64, T int64, lab []int) bool {
    n := len(a)
    if len(lab) != n {
        return false
    }
    // Determine base label for values < T/2 if any
    baseSet := false
    base := 0
    for i := 0; i < n; i++ {
        if 2*a[i] < T {
            base = lab[i]
            baseSet = true
            break
        }
    }
    // If no <T/2, derive from >T/2 if any
    if !baseSet {
        for i := 0; i < n; i++ {
            if 2*a[i] > T {
                base = 1 - lab[i]
                baseSet = true
                break
            }
        }
    }
    // Validate assignments
    prevEq := -1
    for i := 0; i < n; i++ {
        v := a[i]
        if 2*v < T {
            if lab[i] != base {
                return false
            }
        } else if 2*v > T {
            if lab[i] != 1-base {
                return false
            }
        } else { // equal case: must alternate
            if prevEq != -1 && lab[i] == prevEq {
                return false
            }
            prevEq = lab[i]
        }
    }
    return true
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    candidate := os.Args[1]
    tests, err := loadTests()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    for idx, line := range tests {
        n, T, arr, err := parseLine(line)
        if err != nil {
            fmt.Printf("bad test %d: %v\n", idx+1, err)
            os.Exit(1)
        }
        input, _ := lineToInput(line)
        got, err := runProg(candidate, input)
        if err != nil {
            fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
            os.Exit(1)
        }
        labels, err := parseLabels(got, n)
        if err != nil {
            fmt.Printf("Test %d failed\nInput:\n%sError: %v\nOutput:\n%s\n", idx+1, input, err, got)
            os.Exit(1)
        }
        if !checkValid(arr, T, labels) {
            fmt.Printf("Test %d failed\nInput:\n%sGot:\n%s\n", idx+1, input, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(tests))
}
