package main
import (
    "bufio"
    "fmt"
    "os"
)
func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()
    var t int
    fmt.Fscan(reader, &t)
    for tc := 0; tc < t; tc++ {
        var n, m, k int64
        fmt.Fscan(reader, &n, &m, &k)
        stBig := int64(1)
        for i := int64(0); i < k; i++ {
            temp := m
            big := n % m
            stSmall := stBig
            for j := int64(0); j < temp; j++ {
                if big > 0 {
                    size := (n + m - 1) / m
                    fmt.Fprint(writer, size, " ")
                    eles := size
                    for eles > 0 {
                        fmt.Fprint(writer, stBig, " ")
                        stBig++
                        if stBig > n {
                            stBig = 1
                        }
                        stSmall = stBig
                        eles--
                    }
                    fmt.Fprint(writer, "\n")
                    big--
                } else {
                    size := n / m
                    fmt.Fprint(writer, size, " ")
                    eles := size
                    for eles > 0 {
                        fmt.Fprint(writer, stSmall, " ")
                        stSmall++
                        if stSmall > n {
                            stSmall = 1
                        }
                        eles--
                    }
                    fmt.Fprint(writer, "\n")
                }
            }
        }
        fmt.Fprint(writer, "\n")
    }
}
