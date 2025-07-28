package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

func runBinary(path string, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(path, ".go") {
        cmd = exec.Command("go", "run", path)
    } else {
        cmd = exec.Command(path)
    }
    cmd.Stdin = strings.NewReader(input)
    out, err := cmd.CombinedOutput()
    return strings.TrimSpace(string(out)), err
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func solveOnce(a, b []int) int64 {
    n := len(a)
    minA1, minAn := 1<<60, 1<<60
    minB1, minBn := 1<<60, 1<<60
    for i := 0; i < n; i++ {
        v := abs(a[0] - b[i])
        if v < minA1 {
            minA1 = v
        }
        v = abs(a[n-1] - b[i])
        if v < minAn {
            minAn = v
        }
    }
    for i := 0; i < n; i++ {
        v := abs(b[0] - a[i])
        if v < minB1 {
            minB1 = v
        }
        v = abs(b[n-1] - a[i])
        if v < minBn {
            minBn = v
        }
    }
    ans := int64(minA1 + minAn + minB1 + minBn)
    candidate := int64(abs(a[0]-b[0]) + abs(a[n-1]-b[n-1]))
    if candidate < ans {
        ans = candidate
    }
    candidate = int64(abs(a[0]-b[n-1]) + abs(a[n-1]-b[0]))
    if candidate < ans {
        ans = candidate
    }
    candidate = int64(abs(a[0]-b[0]) + minAn + minBn)
    if candidate < ans {
        ans = candidate
    }
    candidate = int64(abs(a[n-1]-b[n-1]) + minA1 + minB1)
    if candidate < ans {
        ans = candidate
    }
    candidate = int64(abs(a[0]-b[n-1]) + minAn + minB1)
    if candidate < ans {
        ans = candidate
    }
    candidate = int64(abs(a[n-1]-b[0]) + minA1 + minBn)
    if candidate < ans {
        ans = candidate
    }
    return ans
}

func main() {
    args := os.Args[1:]
    if len(args) > 0 && args[0] == "--" {
        args = args[1:]
    }
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "usage: verifierC <binary>")
        os.Exit(1)
    }
    path := args[0]
    rand.Seed(time.Now().UnixNano())

    for test := 0; test < 100; test++ {
        n := rand.Intn(5) + 1
        a := make([]int, n)
        b := make([]int, n)
        for i := 0; i < n; i++ {
            a[i] = rand.Intn(20)
        }
        for i := 0; i < n; i++ {
            b[i] = rand.Intn(20)
        }
        var buf bytes.Buffer
        writer := bufio.NewWriter(&buf)
        fmt.Fprintln(writer, 1)
        fmt.Fprintln(writer, n)
        for i := 0; i < n; i++ {
            if i > 0 {
                fmt.Fprint(writer, " ")
            }
            fmt.Fprint(writer, a[i])
        }
        fmt.Fprintln(writer)
        for i := 0; i < n; i++ {
            if i > 0 {
                fmt.Fprint(writer, " ")
            }
            fmt.Fprint(writer, b[i])
        }
        fmt.Fprintln(writer)
        writer.Flush()
        input := buf.String()
        want := fmt.Sprintf("%d", solveOnce(a, b))
        got, err := runBinary(path, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", test+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != want {
            fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected %s got %s\n", test+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

