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

func expected(times []int) string {
    last := 0
    for _, t := range times {
        if t-last > 15 {
            return fmt.Sprintf("%d", last+15)
        }
        last = t
    }
    if last+15 <= 90 {
        return fmt.Sprintf("%d", last+15)
    }
    return "90"
}

func runCase(exe, input, exp string) error {
    var cmd *exec.Cmd
    if strings.HasSuffix(exe, ".go") {
        cmd = exec.Command("go", "run", exe)
    } else {
        cmd = exec.Command(exe)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    got := strings.TrimSpace(out.String())
    exp = strings.TrimSpace(exp)
    if got != exp {
        return fmt.Errorf("expected %q got %q", exp, got)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    exe := os.Args[1]
    data, err := os.ReadFile("testcasesA.txt")
    if err != nil {
        fmt.Println("could not read testcasesA.txt:", err)
        os.Exit(1)
    }
    scan := bufio.NewScanner(bytes.NewReader(data))
    scan.Split(bufio.ScanWords)
    if !scan.Scan() {
        fmt.Println("invalid test file")
        os.Exit(1)
    }
    t, _ := strconv.Atoi(scan.Text())
    for caseIdx := 0; caseIdx < t; caseIdx++ {
        if !scan.Scan() {
            fmt.Println("bad test file")
            os.Exit(1)
        }
        n, _ := strconv.Atoi(scan.Text())
        times := make([]int, n)
        for i := 0; i < n; i++ {
            scan.Scan()
            times[i], _ = strconv.Atoi(scan.Text())
        }
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for i, v := range times {
            if i > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(strconv.Itoa(v))
        }
        sb.WriteByte('\n')
        in := sb.String()
        exp := expected(times)
        if err := runCase(exe, in, exp); err != nil {
            fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}
