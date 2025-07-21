package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   a, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       fmt.Fprintln(os.Stderr, err)
       return
   }
   b, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       fmt.Fprintln(os.Stderr, err)
       return
   }
   a = strings.TrimRight(a, "\r\n")
   b = strings.TrimRight(b, "\r\n")
   // lengths must be equal
   if len(a) != len(b) {
       fmt.Println("NO")
       return
   }
   // count ones
   cntA := strings.Count(a, "1")
   cntB := strings.Count(b, "1")
   // if no ones in a, can only reach b if b also has no ones
   if cntA == 0 {
       if cntB == 0 {
           fmt.Println("YES")
       } else {
           fmt.Println("NO")
       }
       return
   }
   // a has at least one '1'; can reach any b with at least one '1'
   if cntB > 0 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
