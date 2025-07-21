package main

import (
   "fmt"
   "io"
   "os"
)

func main() {
   // Read entire input
   b, err := io.ReadAll(os.Stdin)
   if err != nil {
       return
   }
   // Trim trailing newline or carriage return
   if len(b) > 0 {
       if b[len(b)-1] == '\n' || b[len(b)-1] == '\r' {
           b = b[:len(b)-1]
       }
   }
   n := len(b)
   // Stack to store indices, initialized with -1
   stack := make([]int, 0, n+1)
   stack = append(stack, -1)
   best := 0
   count := 1
   // Iterate over characters
   for i := 0; i < n; i++ {
       if b[i] == '(' {
           stack = append(stack, i)
       } else if b[i] == ')' {
           // Pop
           stack = stack[:len(stack)-1]
           if len(stack) == 0 {
               // Reset base index
               stack = append(stack, i)
           } else {
               length := i - stack[len(stack)-1]
               if length == best {
                   count++
               } else if length > best {
                   best = length
                   count = 1
               }
           }
       }
   }
   if best == 0 {
       // No regular substring
       fmt.Println(0, 1)
   } else {
       fmt.Println(best, count)
   }
}
