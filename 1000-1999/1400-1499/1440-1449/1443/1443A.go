package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var tests int
   if _, err := fmt.Fscan(in, &tests); err != nil {
       return
   }
   for t := 0; t < tests; t++ {
       var n int
       fmt.Fscan(in, &n)
       for i := 0; i < n; i++ {
           // print even numbers in descending order: 4*n - 2*(i+1)
           fmt.Fprintln(out, 4*n-2*(i+1))
       }
       // blank line after each test case
       fmt.Fprintln(out)
   }
}
