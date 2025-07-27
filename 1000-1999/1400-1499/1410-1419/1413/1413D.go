package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   total := 2 * n
   opsType := make([]bool, total)
   opsVal := make([]int, total)
   plusIndices := make([]int, 0, n)
   for i := 0; i < total; i++ {
       var s string
       fmt.Fscan(in, &s)
       if s == "+" {
           opsType[i] = true
           plusIndices = append(plusIndices, i)
       } else {
           opsType[i] = false
           var x int
           fmt.Fscan(in, &x)
           opsVal[i] = x
       }
   }
   res := make([]int, n)
   stack := make([]int, 0, n)
   pi := n - 1
   for i := total - 1; i >= 0; i-- {
       if !opsType[i] {
           x := opsVal[i]
           if len(stack) > 0 && x > stack[len(stack)-1] {
               fmt.Fprintln(out, "NO")
               return
           }
           stack = append(stack, x)
       } else {
           if len(stack) == 0 {
               fmt.Fprintln(out, "NO")
               return
           }
           x := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           res[pi] = x
           pi--
       }
   }
   if pi != -1 {
       fmt.Fprintln(out, "NO")
       return
   }
   fmt.Fprintln(out, "YES")
   for i := 0; i < n; i++ {
       if i > 0 {
           out.WriteByte(' ')
       }
       out.WriteString(strconv.Itoa(res[i]))
   }
   out.WriteByte('\n')
}
