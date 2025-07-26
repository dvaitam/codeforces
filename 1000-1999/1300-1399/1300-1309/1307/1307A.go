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
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, d int
       fmt.Fscan(reader, &n, &d)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // initial haybales in pile 1
       result := a[0]
       days := d
       // move haybales from other piles, starting from closest
       for i := 1; i < n && days > 0; i++ {
           cost := i // number of days to move one bale from pile i+1 to pile 1
           maxMove := days / cost
           if maxMove <= 0 {
               continue
           }
           // can move at most a[i] bales
           moves := a[i]
           if moves > maxMove {
               moves = maxMove
           }
           result += moves
           days -= moves * cost
       }
       fmt.Fprintln(writer, result)
   }
}
