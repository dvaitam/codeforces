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
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // freqRightTotal holds counts of elements in suffix > current j
       freqRightTotal := make([]int, n+1)
       for _, v := range a {
           freqRightTotal[v]++
       }
       freqRight := make([]int, n+1)
       freqLeft := make([]int, n+1)
       var ans int64
       for j := 0; j < n; j++ {
           // remove a[j] from suffix counts
           freqRightTotal[a[j]]--
           // initialize freqRight for this j
           copy(freqRight, freqRightTotal)
           // for each k > j, count pairs (i < j, l > k)
           for k := j + 1; k < n; k++ {
               // remove a[k] from suffix beyond k
               freqRight[a[k]]--
               // add pairs: i with a[i]=a[k], l with a[l]=a[j]
               ans += int64(freqLeft[a[k]] * freqRight[a[j]])
           }
           // include a[j] in prefix counts
           freqLeft[a[j]]++
       }
       fmt.Fprintln(writer, ans)
   }
}
