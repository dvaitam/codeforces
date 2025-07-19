package main

import (
   "bytes"
   "io"
   "os"
)

func main() {
   data, err := io.ReadAll(os.Stdin)
   if err != nil {
       return
   }
   tokens := bytes.Fields(data)
   if len(tokens) < 2 {
       return
   }
   a := tokens[0]
   b := tokens[1]
   n := len(a)
   // Prepare result slice
   res := make([]byte, n)
   // flag indicates next differing char selection: true => pick from b, false => from a
   flag := true
   for i := 0; i < n; i++ {
       if a[i] == b[i] {
           res[i] = a[i]
       } else {
           if flag {
               res[i] = b[i]
           } else {
               res[i] = a[i]
           }
           flag = !flag
       }
   }
   if !flag {
       os.Stdout.WriteString("impossible")
   } else {
       os.Stdout.Write(res)
   }
}
