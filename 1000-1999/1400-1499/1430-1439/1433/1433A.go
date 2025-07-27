package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var x string
       fmt.Fscan(reader, &x)
       // digit value and length
       d := int(x[0] - '0')
       k := len(x)
       // sum of calls for previous digits: each has total presses 1+2+3+4=10
       ans := 10*(d-1) + k*(k+1)/2
       fmt.Println(ans)
   }
}
