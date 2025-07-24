package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k, r int
   if _, err := fmt.Fscan(reader, &k, &r); err != nil {
       return
   }
   for i := 1; i <= 10; i++ {
       total := i * k
       rem := total % 10
       if rem == 0 || rem == r {
           fmt.Println(i)
           return
       }
   }
}
