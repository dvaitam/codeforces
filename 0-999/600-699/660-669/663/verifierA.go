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

func parseRebus(line string) ([]string, int, error) {
    tokens := strings.Fields(line)
    signs := []string{"+"}
    i := 0
    for i < len(tokens) {
        t := tokens[i]
        if t == "=" {
            if i+1 >= len(tokens) {
                return nil, 0, fmt.Errorf("invalid input")
            }
            n, err := strconv.Atoi(tokens[i+1])
            if err != nil {
                return nil, 0, err
            }
            return signs, n, nil
        }
        if t == "+" || t == "-" {
            signs = append(signs, t)
        }
        i++
    }
    return nil, 0, fmt.Errorf("no equal sign")
}

func solve(signs []string, n int) (bool, []int) {
    st := 1
    for j := 1; j < len(signs); j++ {
        if signs[j] == "+" {
            st++
        } else {
            st--
        }
    }
    sol := make([]int, len(signs))
    for j := range sol {
        sol[j] = 1
    }
    for j := 0; j < len(signs); j++ {
        if st > n && signs[j] == "-" {
            diff := st - n
            if diff > n-1 {
                sol[j] = n
                st -= n - 1
            } else {
                sol[j] += diff
                st = n
            }
        }
        if st < n && signs[j] == "+" {
            diff := n - st
            if diff > n-1 {
                sol[j] = n
                st += n - 1
            } else {
                sol[j] += diff
                st = n
            }
        }
    }
    if st == n {
        return true, sol
    }
    return false, nil
}

func run(bin string, input string) (string, error) {
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

func checkOutput(signs []string, n int, output string, expectPossible bool) error {
    lines := strings.Split(strings.TrimSpace(output), "\n")
    if len(lines) == 0 {
        return fmt.Errorf("no output")
    }
    first := strings.ToLower(strings.TrimSpace(lines[0]))
    if expectPossible {
        if first != "possible" {
            return fmt.Errorf("expected Possible got %s", lines[0])
        }
        if len(lines) < 2 {
            return fmt.Errorf("missing solution line")
        }
        tokens := strings.Fields(lines[1])
        nums := []int{}
        sIdx := 0
        i := 0
        for i < len(tokens) {
            val, err := strconv.Atoi(tokens[i])
            if err != nil {
                return fmt.Errorf("invalid number")
            }
            nums = append(nums, val)
            i++
            if i >= len(tokens) {
                return fmt.Errorf("incomplete equation")
            }
            if tokens[i] == "=" {
                i++
                if i >= len(tokens) {
                    return fmt.Errorf("missing n")
                }
                nn, err := strconv.Atoi(tokens[i])
                if err != nil || nn != n {
                    return fmt.Errorf("wrong n")
                }
                i++
                break
            }
            if tokens[i] != "+" && tokens[i] != "-" {
                return fmt.Errorf("missing operator")
            }
            if sIdx+1 >= len(signs) || tokens[i] != signs[sIdx+1] {
                return fmt.Errorf("operator mismatch")
            }
            sIdx++
            i++
        }
        if sIdx != len(signs)-1 || len(nums) != len(signs) {
            return fmt.Errorf("wrong number of terms")
        }
        if i != len(tokens) {
            return fmt.Errorf("extra tokens")
        }
        sum := 0
        for j, v := range nums {
            if v < 1 || v > n {
                return fmt.Errorf("number out of range")
            }
            if j == 0 {
                sum += v
            } else {
                if signs[j] == "+" {
                    sum += v
                } else {
                    sum -= v
                }
            }
        }
        if sum != n {
            return fmt.Errorf("equation does not hold")
        }
    } else {
        if first != "impossible" {
            return fmt.Errorf("expected Impossible got %s", lines[0])
        }
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    file, err := os.Open("testcasesA.txt")
    if err != nil {
        fmt.Println("failed to open testcases:", err)
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
        signs, n, err := parseRebus(line)
        if err != nil {
            fmt.Fprintf(os.Stderr, "bad testcase %d: %v\n", idx, err)
            os.Exit(1)
        }
        possible, _ := solve(signs, n)
        output, err := run(bin, line+"\n")
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
            os.Exit(1)
        }
        if err := checkOutput(signs, n, output, possible); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Println("scanner error:", err)
        os.Exit(1)
    }
    fmt.Println("All tests passed")
}

