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

   var T, n int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       fmt.Fscan(reader, &n)
       if n%2 == 0 {
           // n is even: split into n/2-1, n/2, and 1
           fmt.Fprintf(writer, "%d %d 1\n", n/2-1, n/2)
       } else if (n-1)%4 == 0 {
           // n ≡ 1 mod 4: k = (n-1)/2, split into k-1, k+1, 1
           k := (n - 1) / 2
           fmt.Fprintf(writer, "%d %d 1\n", k-1, k+1)
       } else {
           // n ≡ 3 mod 4: k = (n-1)/2, split into k-2, k+2, 1
           k := (n - 1) / 2
           fmt.Fprintf(writer, "%d %d 1\n", k-2, k+2)
       }
   }
}
