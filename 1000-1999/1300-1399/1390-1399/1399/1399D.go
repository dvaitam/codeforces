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
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       solve(reader, writer)
   }
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n int
   var s string
   fmt.Fscan(reader, &n, &s)
   // stacks for subsequences ending with 0 or 1
   var stack0 []int
   var stack1 []int
   res := make([]int, n)
   cnt := 0
   for i := 0; i < n; i++ {
       bit := s[i] - '0'
       var id int
       if bit == 0 {
           if len(stack1) > 0 {
               // reuse subsequence ending with 1
               id = stack1[len(stack1)-1]
               stack1 = stack1[:len(stack1)-1]
           } else {
               cnt++
               id = cnt
           }
           // now ends with 0
           stack0 = append(stack0, id)
       } else {
           if len(stack0) > 0 {
               id = stack0[len(stack0)-1]
               stack0 = stack0[:len(stack0)-1]
           } else {
               cnt++
               id = cnt
           }
           stack1 = append(stack1, id)
       }
       res[i] = id
   }
   // output
   fmt.Fprintln(writer, cnt)
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(v))
   }
   writer.WriteByte('\n')
}
