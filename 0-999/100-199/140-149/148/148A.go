package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k, l, m, n, d int
   fmt.Fscan(reader, &k, &l, &m, &n, &d)
   count := 0
   for i := 1; i <= d; i++ {
       if i%k == 0 || i%l == 0 || i%m == 0 || i%n == 0 {
           count++
       }
   }
   fmt.Println(count)
}
