package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var T int
    fmt.Fscan(reader, &T)
    for ; T > 0; T-- {
        var n, k int
        var x int64
        fmt.Fscan(reader, &n, &k, &x)
        var s string
        fmt.Fscan(reader, &s)

        letters := make([]byte, 0, len(s))
        bases := make([]int64, 0)
        for i := 0; i < len(s); {
            if s[i] == 'a' {
                letters = append(letters, 'a')
                i++
            } else {
                j := i
                for j < len(s) && s[j] == '*' {
                    j++
                }
                base := int64(j-i)*int64(k) + 1
                bases = append(bases, base)
                letters = append(letters, '*')
                i = j
            }
        }

        x--
        cntB := make([]int64, len(bases))
        for i := len(bases) - 1; i >= 0; i-- {
            cntB[i] = x % bases[i]
            x /= bases[i]
        }

        var sb strings.Builder
        idx := 0
        for _, ch := range letters {
            if ch == 'a' {
                sb.WriteByte('a')
            } else {
                if cntB[idx] > 0 {
                    sb.WriteString(strings.Repeat("b", int(cntB[idx])))
                }
                idx++
            }
        }
        fmt.Fprintln(writer, sb.String())
    }
}

