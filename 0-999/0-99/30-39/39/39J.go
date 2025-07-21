package main

import (
   "bufio"
   "bytes"
   "fmt"
   "io"
   "os"
   "strconv"
)

func main() {
   data, err := io.ReadAll(os.Stdin)
   if err != nil {
       fmt.Fprintln(os.Stderr, err)
       return
   }
   fields := bytes.Fields(data)
   if len(fields) < 2 {
       fmt.Println(0)
       return
   }
   s := fields[0]
   t := fields[1]
   n := len(s)
   if len(t) != n-1 {
       fmt.Println(0)
       return
   }
   // prefix matches: pre[i] = s[0:i] == t[0:i]
   pre := make([]bool, n)
   pre[0] = true
   for i := 1; i < n; i++ {
       if pre[i-1] && s[i-1] == t[i-1] {
           pre[i] = true
       } else {
           pre[i] = false
       }
   }
   // suffix matches: suf[i] = s[i+1:] == t[i:]
   suf := make([]bool, n)
   suf[n-1] = true
   // t indexes 0..n-2
   for i := n - 2; i >= 0; i-- {
       if s[i+1] == t[i] && suf[i+1] {
           suf[i] = true
       } else {
           suf[i] = false
       }
   }
   // collect positions
   var pos []int
   for i := 0; i < n; i++ {
       if pre[i] && suf[i] {
           pos = append(pos, i+1)
       }
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   if len(pos) == 0 {
       w.WriteString("0\n")
       return
   }
   w.WriteString(strconv.Itoa(len(pos)))
   w.WriteByte('\n')
   for i, v := range pos {
       if i > 0 {
           w.WriteByte(' ')
       }
       w.WriteString(strconv.Itoa(v))
   }
}
