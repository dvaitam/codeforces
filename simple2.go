package main
import (
 "bufio"
 "fmt"
 "os"
)
func main() {
 reader := bufio.NewReader(os.Stdin)
 var x int
 if _, err := fmt.Fscan(reader, &x); err != nil {
     fmt.Println("err", err)
     return
 }
 fmt.Println("x", x)
}
