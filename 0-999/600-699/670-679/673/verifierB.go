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

func expected(n int, a, b []int) string {
    L := 1
    R := n - 1
    for i := 0; i < len(a); i++ {
        u := a[i]
        v := b[i]
        if u > v {
            u, v = v, u
        }
        if u > L {
            L = u
        }
        if v-1 < R {
            R = v - 1
        }
    }
    if L > R {
        return "0"
    }
    return fmt.Sprintf("%d", R-L+1)
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
        fmt.Println("Usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    exe := os.Args[1]
    data, err := os.ReadFile("testcasesB.txt")
    if err != nil {
        fmt.Println("could not read testcasesB.txt:", err)
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
        scan.Scan()
        m, _ := strconv.Atoi(scan.Text())
        a := make([]int, m)
        b := make([]int, m)
        for i := 0; i < m; i++ {
            scan.Scan()
            a[i], _ = strconv.Atoi(scan.Text())
            scan.Scan()
            b[i], _ = strconv.Atoi(scan.Text())
        }
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
        for i := 0; i < m; i++ {
            sb.WriteString(fmt.Sprintf("%d %d\n", a[i], b[i]))
        }
        in := sb.String()
        exp := expected(n, a, b)
        if err := runCase(exe, in, exp); err != nil {
            fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}
