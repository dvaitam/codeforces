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

func runCandidate(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
    }
    return strings.TrimSpace(out.String()), nil
}

// Feasibility check from problem logic
func ok(a []int, n, m, x int) bool {
    cur := 0
    for i := 0; i < n; i++ {
        ai := a[i]
        if ai+x < m {
            if ai < cur {
                return false
            }
            cur = ai
        } else {
            wrap := (ai + x) % m
            if wrap < cur {
                if ai < cur {
                    return false
                }
                cur = ai
            }
        }
    }
    return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
    bin := os.Args[1]

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "empty test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing test case %d\n", i+1)
			os.Exit(1)
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			i--
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", i+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+n {
			fmt.Fprintf(os.Stderr, "test %d expected %d numbers got %d\n", i+1, 2+n, len(parts))
			os.Exit(1)
		}
		arr := strings.Join(parts[2:], " ")
        input := fmt.Sprintf("%d %d\n%s\n", n, m, arr)
        gotStr, err := runCandidate(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
            os.Exit(1)
        }
        // parse array values
        nums := strings.Fields(arr)
        a := make([]int, n)
        for idx := 0; idx < n; idx++ {
            v, _ := strconv.Atoi(nums[idx])
            a[idx] = v
        }
        // parse candidate answer
        val, err := strconv.Atoi(strings.TrimSpace(gotStr))
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d failed: invalid output %q\n", i+1, gotStr)
            os.Exit(1)
        }
        // validate minimality
        if val < 0 || val > m {
            fmt.Fprintf(os.Stderr, "test %d failed: out of range %d\n", i+1, val)
            os.Exit(1)
        }
        if !ok(a, n, m, val) {
            fmt.Fprintf(os.Stderr, "test %d failed: answer %d not feasible\n", i+1, val)
            os.Exit(1)
        }
        if val > 0 && ok(a, n, m, val-1) {
            fmt.Fprintf(os.Stderr, "test %d failed: answer %d not minimal\n", i+1, val)
            os.Exit(1)
        }
    }
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
