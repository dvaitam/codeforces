package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, p int
   var s string
   fmt.Fscan(in, &n, &p)
   fmt.Fscan(in, &s)
   p-- // zero-based index

   flag1, flag2 := 0, 0
   if p == 0 {
       flag1, flag2 = 1, 0
   } else if p == n-1 {
       flag1, flag2 = 0, 0
   } else if p < n-1-p {
       flag1, flag2 = 0, 1
   } else {
       flag1, flag2 = 1, 1
   }

   var cmds []string
   // initial print
   cmds = append(cmds, fmt.Sprintf("PRINT %c", s[p]))

   // first phase: go right or left
   if flag1 == 1 {
       cnt := p + 1
       // print from p+1 to n-2
       for cnt < n-1 {
           cmds = append(cmds, "RIGHT")
           cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
           cnt++
       }
       // last at n-1
       if cnt < n {
           cmds = append(cmds, "RIGHT")
           cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
       }
       // set cnt to last printed index
       cnt = n - 1
       // second phase if needed
       if flag2 == 1 {
           // move back to before p
           for cnt >= p {
               cmds = append(cmds, "LEFT")
               cnt--
           }
           // print remaining left of p
           for cnt > 0 {
               cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
               cmds = append(cmds, "LEFT")
               cnt--
           }
           // finally at 0
           if cnt >= 0 {
               cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
           }
       }
   } else {
       cnt := p - 1
       // print from p-1 down to 1
       for cnt > 0 {
           cmds = append(cmds, "LEFT")
           cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
           cnt--
       }
       // last at 0
       if cnt >= 0 {
           cmds = append(cmds, "LEFT")
           cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
       }
       // set cnt to last printed index (0)
       cnt = 0
       // second phase if needed
       if flag2 == 1 {
           // move to right of p
           for cnt <= p {
               cmds = append(cmds, "RIGHT")
               cnt++
           }
           // print remaining right of p
           for cnt < n-1 {
               cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
               cmds = append(cmds, "RIGHT")
               cnt++
           }
           // last at n-1
           if cnt < n {
               cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
           }
       }
   }

   // output commands
   out := bufio.NewWriter(os.Stdout)
   for _, c := range cmds {
       fmt.Fprintln(out, c)
   }
   out.Flush()
}
