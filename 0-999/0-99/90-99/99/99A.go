package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   fmt.Fscan(reader, &s)
   // split integer and fractional parts
   parts := strings.SplitN(s, ".", 2)
   intPart := parts[0]
   fracPart := parts[1]
   // if last digit of integer part is 9, go to Vasilisa
   if intPart[len(intPart)-1] == '9' {
       fmt.Println("GOTO Vasilisa.")
       return
   }
   // decide rounding
   if fracPart[0] >= '5' {
       // increment last digit (no carry, since not '9')
       b := []byte(intPart)
       b[len(b)-1]++
       intPart = string(b)
   }
   // remove leading zeros, leave single zero if needed
   i := 0
   for i < len(intPart)-1 && intPart[i] == '0' {
       i++
   }
   intPart = intPart[i:]
   fmt.Println(intPart)
}
