package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "sort"
    "strings"
)

func solveCase(n int, arr []int) int {
    sort.Sort(sort.Reverse(sort.IntSlice(arr)))
    ans := 0
    for i := 0; i < n; i++ {
        if arr[i] >= i+1 {
            ans = i + 1
        }
    }
    return ans
}

func runProg(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) == 3 && os.Args[1] == "--" {
        os.Args = append([]string{os.Args[0]}, os.Args[2])
    }
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]

    f, err := os.Open("testcasesA.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "could not open testcasesA.txt: %v\n", err)
        os.Exit(1)
    }
    defer f.Close()
    reader := bufio.NewReader(f)

    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        fmt.Fprintf(os.Stderr, "failed to read test count: %v\n", err)
        os.Exit(1)
    }

    for caseNum := 1; caseNum <= t; caseNum++ {
        var n int
        if _, err := fmt.Fscan(reader, &n); err != nil {
            fmt.Fprintf(os.Stderr, "case %d: read n: %v\n", caseNum, err)
            os.Exit(1)
        }
        arr := make([]int, n)
        for i := 0; i < n; i++ {
            if _, err := fmt.Fscan(reader, &arr[i]); err != nil {
                fmt.Fprintf(os.Stderr, "case %d: read arr: %v\n", caseNum, err)
                os.Exit(1)
            }
        }
        var input strings.Builder
        fmt.Fprintf(&input, "1\n%d\n", n)
        for i, v := range arr {
            if i > 0 {
                input.WriteByte(' ')
            }
            fmt.Fprintf(&input, "%d", v)
        }
        input.WriteByte('\n')

        want := fmt.Sprintf("%d", solveCase(n, append([]int(nil), arr...)))
        got, err := runProg(bin, input.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum, err, input.String())
            os.Exit(1)
        }
        if strings.TrimSpace(got) != strings.TrimSpace(want) {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum, want, got, input.String())
            os.Exit(1)
        }
    }

    fmt.Printf("All %d tests passed\n", t)
}

