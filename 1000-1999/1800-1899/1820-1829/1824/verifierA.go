package main

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "fmt"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"
)

// Base64-encoded contents of testcasesA.txt.
const testcasesA64 = "MTAwCjEgMQoxCjMgMgotMiAtMiAtMQo1IDEKMSAtMiAtMSAtMSAtMQoyIDIKLTEgLTIKNSA0Ci0yIDIgMyAxIC0xCjIgNAoxIDEKMiAyCjEgLTEKMSA0Ci0xCjMgMwozIDEgLTEKNCA1Ci0xIDIgLTEgMwozIDUKNSAxIDMKMiAxCi0xIDEKMiAzCi0xIC0yCjEgNAoxCjQgMwotMiAxIDEgLTIKMyAxCjEgMSAtMQo1IDIKLTIgMiAyIDEgLTIKMyAxCi0xIC0xIC0yCjQgMwotMSAtMiAzIDEKMiA0CjIgNAo0IDIKMSAtMiAtMiAxCjUgNAozIDIgMSAtMiAtMgo1IDQKLTEgLTEgLTEgLTIgNAoyIDQKMyAtMQo0IDQKMyAyIDMgMQo1IDEKMSAxIC0xIDEgMQozIDMKLTEgMSAyCjIgNAotMSA0CjMgNQo1IC0yIDMKMSAzCjMKNSAyCi0yIDEgLTIgLTEgMQo0IDEKLTEgLTIgLTIgLTEKMSAyCi0xCjEgNAotMQo1IDIKLTIgMiAtMiAxIDIKMiA1CjUgNAoyIDMKMiAxCjQgNQoyIC0xIC0yIC0yCjEgMwotMQo1IDUKLTIgMyAtMiAtMSAtMQoxIDIKLTEKMSAzCi0xCjUgMgoxIDIgLTIgLTIgMgoyIDQKMiAtMgoxIDEKMQo0IDMKMiAyIDIgLTEKMSAxCi0yCjMgMQotMSAtMSAtMQo1IDQKLTIgMiAtMiAxIDIKMiAxCi0yIDEKMSAxCjEKNSAxCi0xIC0xIC0xIC0yIC0yCjQgMgoyIC0xIC0yIDIKMSA0CjEKNCAzCjIgMyAyIC0yCjIgMwotMiAtMQo1IDUKLTEgNCAxIC0xIC0xCjUgNAozIDMgLTIgLTEgMwoxIDIKLTEKNSAxCjEgLTEgLTIgLTEgMQoyIDUKMyAtMQo1IDEKLTIgMSAxIDEgMQozIDMKLTIgMSAtMgozIDQKLTIgNCA0CjMgNAoxIC0xIC0xCjQgNQozIC0xIC0xIDMKMiA1CjEgLTIKMyAxCi0xIC0yIC0yCjIgNAozIDQKMyA1CjUgNCAzCjEgNQoxCjEgMgoxCjEgMQoxCjUgMgoxIDEgLTIgMSAtMgozIDUKMiAxIDUKMSAxCjEKNCAzCi0xIC0xIDEgLTIKMyAyCjIgMiAtMQoxIDEKMQoyIDUKLTEgNQozIDUKMyAtMiAyCjIgMQotMiAtMgoxIDMKLTIKMiAxCi0yIDEKNCA1CjQgLTIgLTIgNQoyIDIKMiAtMQoyIDMKMiAtMgozIDIKLTEgMiAtMQo0IDIKLTIgMiAxIDEKMiAyCi0xIC0yCjQgMwoxIC0xIDEgMQo1IDQKNCAzIDEgLTEgLTEKMyAyCjEgLTEgLTEKNSA0CjEgNCAxIDIgMwo1IDEKLTIgMSAtMSAtMiAtMQo0IDEKMSAxIDEgMQoyIDMKMiAtMQozIDUKMSA0IDUKMSAzCjMKMyA0CjEgMiA0Cg=="

type testCase struct {
    n, m int
    vals []int
}

// Embedded solution logic from 1824A.go.
func solve(tc testCase) string {
    lcnt, rcnt := 0, 0
    posMap := make(map[int]bool)
    for _, v := range tc.vals {
        switch v {
        case -1:
            lcnt++
        case -2:
            rcnt++
        default:
            if v >= 1 && v <= tc.m {
                posMap[v] = true
            }
        }
    }
    pos := make([]int, 0, len(posMap))
    for k := range posMap {
        pos = append(pos, k)
    }
    sort.Ints(pos)
    k := len(pos)
    ans := 0
    if tmp := lcnt + k; tmp > ans {
        if tmp > tc.m {
            tmp = tc.m
        }
        if tmp > ans {
            ans = tmp
        }
    }
    if tmp := rcnt + k; tmp > ans {
        if tmp > tc.m {
            tmp = tc.m
        }
        if tmp > ans {
            ans = tmp
        }
    }
    for i, p := range pos {
        left := p - 1
        if left > lcnt+i {
            left = lcnt + i
        }
        right := tc.m - p
        if right > rcnt+(k-i-1) {
            right = rcnt + (k - i - 1)
        }
        cur := 1 + left + right
        if cur > ans {
            ans = cur
        }
    }
    if ans > tc.m {
        ans = tc.m
    }
    return strconv.Itoa(ans)
}

func parseTestcases() ([]testCase, error) {
    raw, err := base64.StdEncoding.DecodeString(testcasesA64)
    if err != nil {
        return nil, err
    }
    sc := bufio.NewScanner(bytes.NewReader(raw))
    sc.Split(bufio.ScanWords)
    if !sc.Scan() {
        return nil, fmt.Errorf("invalid test data")
    }
    t, err := strconv.Atoi(sc.Text())
    if err != nil {
        return nil, fmt.Errorf("parse t: %v", err)
    }
    cases := make([]testCase, 0, t)
    for i := 0; i < t; i++ {
        if !sc.Scan() {
            return nil, fmt.Errorf("case %d missing n", i+1)
        }
        n, err := strconv.Atoi(sc.Text())
        if err != nil {
            return nil, fmt.Errorf("case %d n: %v", i+1, err)
        }
        if !sc.Scan() {
            return nil, fmt.Errorf("case %d missing m", i+1)
        }
        m, err := strconv.Atoi(sc.Text())
        if err != nil {
            return nil, fmt.Errorf("case %d m: %v", i+1, err)
        }
        vals := make([]int, n)
        for j := 0; j < n; j++ {
            if !sc.Scan() {
                return nil, fmt.Errorf("case %d truncated vals", i+1)
            }
            v, err := strconv.Atoi(sc.Text())
            if err != nil {
                return nil, fmt.Errorf("case %d vals[%d]: %v", i+1, j, err)
            }
            vals[j] = v
        }
        cases = append(cases, testCase{n: n, m: m, vals: vals})
    }
    return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
    if strings.HasSuffix(path, ".go") {
        tmp, err := os.CreateTemp("", "bin*")
        if err != nil {
            return "", nil, err
        }
        tmp.Close()
        if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
            os.Remove(tmp.Name())
            return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
        }
        return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
    }
    return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
    cmd := exec.Command(bin)
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

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }

    cases, err := parseTestcases()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    bin, cleanup, err := buildIfGo(os.Args[1])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer cleanup()

    for idx, tc := range cases {
        var input strings.Builder
        input.WriteString("1\n")
        fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
        for i, v := range tc.vals {
            if i > 0 {
                input.WriteByte(' ')
            }
            input.WriteString(strconv.Itoa(v))
        }
        input.WriteByte('\n')

        want := solve(tc)
        got, err := runCandidate(bin, input.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input.String(), want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}
