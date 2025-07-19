package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewScanner(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   for reader.Scan() {
       s := reader.Text()
       // prepend a 'b' to simplify calculation
       s = "b" + s
       var ans int64 = 1
       var cnt int64 = 0
       // iterate from end to start
       for i := len(s) - 1; i >= 0; i-- {
           switch s[i] {
           case 'a':
               cnt++
           case 'b':
               cnt++
               ans = ans * cnt % MOD
               cnt = 0
           }
       }
       // subtract the empty sequence
       res := (ans - 1) % MOD
       if res < 0 {
           res += MOD
       }
       fmt.Fprintln(writer, res)
   }
}
