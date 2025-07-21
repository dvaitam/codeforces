package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var path string
   if _, err := fmt.Fscan(reader, &path); err != nil {
       return
   }
   parts := strings.Split(path, "/")
   var stack []string
   for _, p := range parts {
       if p == "" || p == "." {
           continue
       }
       if p == ".." {
           if len(stack) == 0 {
               fmt.Println(-1)
               return
           }
           stack = stack[:len(stack)-1]
       } else {
           stack = append(stack, p)
       }
   }
   if len(stack) == 0 {
       fmt.Println("/")
   } else {
       fmt.Println("/" + strings.Join(stack, "/"))
   }
}
