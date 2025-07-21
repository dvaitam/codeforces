package main

import "fmt"

func main() {
    var s string
    if _, err := fmt.Scan(&s); err != nil {
        return
    }
    // thresholds for each type
    limits := []struct{
        lim  string
        name string
    }{
        {"127", "byte"},
        {"32767", "short"},
        {"2147483647", "int"},
        {"9223372036854775807", "long"},
    }
    for _, lt := range limits {
        if fits(s, lt.lim) {
            fmt.Println(lt.name)
            return
        }
    }
    fmt.Println("BigInteger")
}

// fits reports whether the decimal string s represents a number <= lim.
func fits(s, lim string) bool {
    if len(s) != len(lim) {
        return len(s) < len(lim)
    }
    return s <= lim
}
