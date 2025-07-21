package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read base a and base b (or R)
   var a int
   var bstr string
   if _, err := fmt.Fscan(reader, &a, &bstr); err != nil {
       return
   }
   // read number c in base a
   var cstr string
   if _, err := fmt.Fscan(reader, &cstr); err != nil {
       return
   }
   cstr = strings.TrimSpace(cstr)
   // parse c into big.Int
   num := new(big.Int)
   if _, ok := num.SetString(cstr, a); !ok {
       // invalid input
       return
   }
   // convert
   if bstr == "R" {
       // Roman numeral
       // get integer value
       val := num.Int64()
       // generate roman
       fmt.Print(toRoman(int(val)))
   } else {
       // target base
       b, err := strconv.Atoi(bstr)
       if err != nil {
           return
       }
       // output in base b
       out := num.Text(b)
       out = strings.ToUpper(out)
       // trim leading zeros
       out = strings.TrimLeft(out, "0")
       if out == "" {
           out = "0"
       }
       fmt.Print(out)
   }
}

// toRoman converts positive integers to Roman numerals
func toRoman(n int) string {
   // mapping of values to symbols
   vals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
   syms := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
   var sb strings.Builder
   for i := 0; i < len(vals); i++ {
       for n >= vals[i] {
           sb.WriteString(syms[i])
           n -= vals[i]
       }
   }
   return sb.String()
}
