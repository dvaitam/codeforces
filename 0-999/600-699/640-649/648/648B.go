package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   size := 2 * n
   a := make([]int, size)
   for i := 0; i < size; i++ {
       fmt.Fscan(in, &a[i])
   }
   var l int
   // find sum l that works for all elements
   for k := 1; k < size; k++ {
       lCand := a[0] + a[k]
       ok := true
       for i := 0; i < size && ok; i++ {
           found := false
           for j := 0; j < size; j++ {
               if i != j && a[i]+a[j] == lCand {
                   found = true
                   break
               }
           }
           if !found {
               ok = false
           }
       }
       if ok {
           l = lCand
           break
       }
   }
   removed := make([]bool, size)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i := 0; i < size; i++ {
       if removed[i] {
           continue
       }
       for j := i + 1; j < size; j++ {
           if !removed[j] && a[i]+a[j] == l {
               fmt.Fprintln(out, a[i], a[j])
               removed[j] = true
               break
           }
       }
   }
}
