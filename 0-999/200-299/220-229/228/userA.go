package main
import "fmt"
func main() {
    var s1, s2, s3, s4 int
    fmt.Scan(&s1, &s2, &s3, &s4)
    m := make(map[int]struct{})
    m[s1] = struct{}{}
    m[s2] = struct{}{}
    m[s3] = struct{}{}
    m[s4] = struct{}{}
    fmt.Println(4 - len(m))
}