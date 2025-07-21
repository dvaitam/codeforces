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

   var N, D int
   if _, err := fmt.Fscan(reader, &N, &D); err != nil {
       return
   }
   // We construct N-2 ones, one 2, and one (D+N)
   // Then product = 2*(D+N), sum = (N-2)+2+(D+N) = 2N+D, so product - sum = D
   for i := 0; i < N-2; i++ {
       fmt.Fprint(writer, 1, " ")
   }
   if N >= 2 {
       fmt.Fprint(writer, 2)
       fmt.Fprint(writer, " ")
       fmt.Fprint(writer, N+D)
   }
   // If N == 1 (though constraints say N >= 2), just output D+1
   if N == 1 {
       fmt.Fprint(writer, D+1)
   }
}
