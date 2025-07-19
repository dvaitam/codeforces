package main

import (
    "bufio"
    "fmt"
    "os"
)

func isRegular(s string) bool {
    acc := 0
    for i := 0; i < len(s); i++ {
        if s[i] == '(' {
            acc++
        } else {
            acc--
        }
        if acc < 0 {
            return false
        }
    }
    return acc == 0
}

func isBalanced(s string) bool {
    dif := 0
    for i := 0; i < len(s); i++ {
        if s[i] == '(' {
            dif++
        } else {
            dif--
        }
    }
    return dif == 0
}

func reverseString(s string) string {
    bs := []byte(s)
    i, j := 0, len(bs)-1
    for i < j {
        bs[i], bs[j] = bs[j], bs[i]
        i++
        j--
    }
    return string(bs)
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var tc int
    if _, err := fmt.Fscan(reader, &tc); err != nil {
        return
    }
    for t := 0; t < tc; t++ {
        var n int
        var s string
        fmt.Fscan(reader, &n, &s)
        r := reverseString(s)

        ans := make([]int, n)
        for i := 0; i < n; i++ {
            ans[i] = 1
        }

        switch {
        case isRegular(s) || isRegular(r):
            fmt.Fprintln(writer, 1)
        case isBalanced(s):
            half := n / 2
            for i := 0; i < half; i++ {
                if s[i] == ')' {
                    ans[i] = 2
                }
            }
            for i := half; i < n; i++ {
                if s[i] == '(' {
                    ans[i] = 2
                }
            }
            fmt.Fprintln(writer, 2)
        default:
            fmt.Fprintln(writer, -1)
        }
        // output ans if valid
        if isRegular(s) || isRegular(r) || isBalanced(s) {
            for i := 0; i < n; i++ {
                if i > 0 {
                    writer.WriteByte(' ')
                }
                fmt.Fprint(writer, ans[i])
            }
            writer.WriteByte('\n')
        }
    }
}
