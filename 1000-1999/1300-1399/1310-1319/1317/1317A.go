package main
import "fmt"
func main() {
    var a, b int
    fmt.Scan(&a, &b)
    arr := []int{}
    // This will panic when a != 0
    fmt.Println(arr[a+b])
}
