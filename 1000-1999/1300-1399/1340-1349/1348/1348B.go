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

   var TC int
   if _, err := fmt.Fscan(in, &TC); err != nil {
       return
   }
   for TC > 0 {
       TC--
       var n, k int
       fmt.Fscan(in, &n, &k)
       vals := make([]int, n)
       distinct := make(map[int]bool, n)
       beauty := make([]int, 0, k)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &vals[i])
           if !distinct[vals[i]] {
               distinct[vals[i]] = true
               beauty = append(beauty, vals[i])
           }
       }
       if len(beauty) > k {
           fmt.Fprintln(out, -1)
           continue
       }
       // add filler values
       for i := 1; i <= n && len(beauty) < k; i++ {
           if !distinct[i] {
               distinct[i] = true
               beauty = append(beauty, i)
           }
       }
       // output total length and sequence
       total := n * len(beauty)
       fmt.Fprintln(out, total)
       for i := 0; i < n; i++ {
           for _, v := range beauty {
               fmt.Fprint(out, v, " ")
           }
       }
       fmt.Fprintln(out)
   }
}
