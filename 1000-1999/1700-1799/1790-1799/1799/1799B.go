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
       solve(reader, writer)
   }
}

func solve(r *bufio.Reader, w *bufio.Writer) {
   var n int
   fmt.Fscan(r, &n)
   arr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(r, &arr[i])
   }
   uniq := make(map[int]struct{})
   for _, v := range arr {
       uniq[v] = struct{}{}
   }
   if len(uniq) == 1 {
       fmt.Fprintln(w, 0)
       return
   }
   for _, v := range arr {
       if v == 1 {
           fmt.Fprintln(w, -1)
           return
       }
   }
   var ans [][2]int
   for {
       mn := arr[0]
       ind := 0
       for i, v := range arr {
           if v < mn {
               mn = v
               ind = i
           }
       }
       for i := 0; i < n; i++ {
           for arr[i] > mn {
               ans = append(ans, [2]int{i + 1, ind + 1})
               // ceil division
               arr[i] = (arr[i] + mn - 1) / mn
           }
       }
       uniq = make(map[int]struct{})
       for _, v := range arr {
           uniq[v] = struct{}{}
       }
       if len(uniq) == 1 {
           break
       }
   }
   fmt.Fprintln(w, len(ans))
   for _, p := range ans {
       fmt.Fprintln(w, p[0], p[1])
   }
}
