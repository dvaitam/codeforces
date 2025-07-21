package main

import "fmt"

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   // use a stack to process operations
   stack := make([]int, 0, len(s))
   for _, c := range s {
       switch c {
       case '+', '*':
           n := len(stack)
           // pop two values
           a := stack[n-1]
           b := stack[n-2]
           stack = stack[:n-2]
           // apply operation and push result
           if c == '+' {
               stack = append(stack, b+a)
           } else {
               stack = append(stack, b*a)
           }
       default:
           // digit: push its value
           if c >= '0' && c <= '9' {
               stack = append(stack, int(c-'0'))
           }
       }
   }
   // output top of stack
   if len(stack) > 0 {
       fmt.Println(stack[len(stack)-1])
   }
}
