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
   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for t := 0; t < T; t++ {
       var A, B int
       var s string
       fmt.Fscan(reader, &A, &B)
       fmt.Fscan(reader, &s)
       n := len(s)
       arr := []byte(s)
       wenhao := 0
       for i := 0; i < n; i++ {
           switch arr[i] {
           case '?':
               wenhao++
           case '0':
               A--
           case '1':
               B--
           }
       }
       // Mirror known characters
       for i := 0; i < n; i++ {
           j := n - 1 - i
           if arr[i] != '?' && arr[j] == '?' {
               if arr[i] == '0' {
                   arr[j] = '0'
                   A--
               } else {
                   arr[j] = '1'
                   B--
               }
               wenhao--
           }
       }
       // Fill remaining placeholders
       for i := n / 2; i < n; i++ {
           if arr[i] != '?' {
               continue
           }
           j := n - 1 - i
           if i == j {
               if A%2 == 1 {
                   arr[i] = '0'
                   A--
               } else if B%2 == 1 {
                   arr[i] = '1'
                   B--
               }
               wenhao--
               continue
           }
           if A >= 2 {
               arr[i], arr[j] = '0', '0'
               A -= 2
               wenhao -= 2
           } else if B >= 2 {
               arr[i], arr[j] = '1', '1'
               B -= 2
               wenhao -= 2
           }
       }
       // Validate result
       ok := true
       for i := 0; i < n; i++ {
           if arr[i] != arr[n-1-i] {
               ok = false
               break
           }
       }
       if !ok || A != 0 || B != 0 || wenhao != 0 {
           writer.WriteString("-1\n")
       } else {
           writer.Write(arr)
           writer.WriteByte('\n')
       }
   }
}
