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

   var n, k, w int
   fmt.Fscan(reader, &n, &k, &w)
   var str string
   fmt.Fscan(reader, &str)
   
   // Prefix sums for each residue mod k
   PS := make([][]int, k)
   for i := range PS {
       PS[i] = make([]int, n+1)
   }
   // Total ones prefix sum
   totPS := make([]int, n+1)
   for i := 1; i <= n; i++ {
       c := str[i-1]
       totPS[i] = totPS[i-1]
       // copy previous counts
       for r := 0; r < k; r++ {
           PS[r][i] = PS[r][i-1]
       }
       if c == '1' {
           totPS[i]++
           r := (i - 1) % k
           PS[r][i]++
       }
   }

   for qi := 0; qi < w; qi++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       // total ones in [l, r]
       total1 := totPS[r] - totPS[l-1]
       // desired residue class = (l-1) mod k
       t := (l - 1) % k
       // number of ones at desired positions
       good1 := PS[t][r] - PS[t][l-1]
       // number of desired positions
       m := (r - l + 1) / k
       // operations = remove unwanted ones + add missing ones
       ans := total1 + m - 2*good1
       fmt.Fprintln(writer, ans)
   }
}
