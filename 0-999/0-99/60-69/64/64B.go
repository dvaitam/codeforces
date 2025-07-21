package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   input, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       // unexpected error
       return
   }
   input = strings.TrimSpace(input)
   if len(input) != 3 {
       return
   }
   a := int(input[0] - '0')
   op := input[1]
   b := int(input[2] - '0')
   var res int
   if op == '+' {
       res = a + b
   } else if op == '-' {
       res = a - b
   } else {
       return
   }
   fmt.Println(res)
}
