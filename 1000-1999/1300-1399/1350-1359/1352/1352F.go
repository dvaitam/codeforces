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

   var tc int
   fmt.Fscan(reader, &tc)
   for tc > 0 {
       tc--
       var n0, n1p, n2p int
       fmt.Fscan(reader, &n0, &n1p, &n2p)
       // Map inputs to C++ code vars: n1=n2p, n2=n1p, n3=n0
       n1 := n2p
       n2 := n1p
       n3 := n0

       var sBuilder []byte
       var tBuilder []byte
       // initialize tBuilder with '1'
       tBuilder = append(tBuilder, '1')
       cFlag := false

       if n2 == 0 {
           if n1 > 0 {
               sBuilder = append(sBuilder, '1')
               for i := 0; i < n1; i++ {
                   sBuilder = append(sBuilder, '1')
               }
           } else {
               sBuilder = append(sBuilder, '0')
               for i := 0; i < n3; i++ {
                   sBuilder = append(sBuilder, '0')
               }
           }
       } else if n2%2 == 1 {
           cnt := n2
           for cnt > 0 {
               if cFlag {
                   tBuilder = append(tBuilder, '1')
                   cFlag = false
               } else {
                   tBuilder = append(tBuilder, '0')
                   cFlag = true
               }
               cnt--
           }
           for i := 0; i < n1; i++ {
               sBuilder = append(sBuilder, '1')
           }
           // append tBuilder to sBuilder
           sBuilder = append(sBuilder, tBuilder...)
           for i := 0; i < n3; i++ {
               sBuilder = append(sBuilder, '0')
           }
       } else {
           // even n2 > 0
           cnt := n2 - 1
           for cnt > 0 {
               if cFlag {
                   tBuilder = append(tBuilder, '1')
                   cFlag = false
               } else {
                   tBuilder = append(tBuilder, '0')
                   cFlag = true
               }
               cnt--
           }
           for i := 0; i < n3; i++ {
               tBuilder = append(tBuilder, '0')
           }
           tBuilder = append(tBuilder, '1')
           for i := 0; i < n1; i++ {
               sBuilder = append(sBuilder, '1')
           }
           sBuilder = append(sBuilder, tBuilder...)
       }

       // write result
       writer.Write(sBuilder)
       writer.WriteByte('\n')
   }
}
