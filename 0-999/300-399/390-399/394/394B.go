package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

// powMod computes (base^exp) mod m.
func powMod(base, exp, m int) int {
   result := 1 % m
   b := base % m
   for exp > 0 {
       if exp&1 != 0 {
           result = (result * b) % m
       }
       b = (b * b) % m
       exp >>= 1
   }
   return result
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var p, x int
   if _, err := fmt.Fscan(reader, &p, &x); err != nil {
       return
   }
   D := 10*x - 1
   // compute 10^p mod D
   r := powMod(10, p, D)
   // try digits dp from 1 to 9
   for dp := 1; dp <= 9; dp++ {
       if (dp*((r-1+D)%D))%D != 0 {
           continue
       }
       g := gcd(dp, D)
       D1 := D / g
       // check that 10^p mod D1 == 1
       if powMod(10, p, D1) != 1 {
           continue
       }
       // find order L1 of 10 mod D1
       L1 := 1
       for ; L1 <= D1; L1++ {
           if powMod(10, L1, D1) == 1 {
               break
           }
       }
       // compute cycle C = (10^L1 - 1) / D1 as string of length L1
       C := make([]byte, L1)
       cur := 0
       for i := 0; i < L1; i++ {
           cur = cur*10 + 9
           C[i] = byte(cur/D1) + '0'
           cur %= D1
       }
       // repeat C to length p
       reps := p / L1
       E := make([]byte, p)
       for i := 0; i < reps; i++ {
           copy(E[i*L1:(i+1)*L1], C)
       }
       // multiply E by m = dp/g
       m := dp / g
       carry := 0
       for i := p - 1; i >= 0; i-- {
           prod := int(E[i]-'0')*m + carry
           E[i] = byte(prod%10) + '0'
           carry = prod / 10
       }
       // prepend carry if any
       if carry > 0 {
           carryStr := fmt.Sprint(carry)
           E = append([]byte(carryStr), E...)
       }
       // result should have exactly p digits
       if len(E) != p {
           // skip invalid
           continue
       }
       // leading digit must not be '0'
       if E[0] == '0' {
           continue
       }
       // last digit must equal dp
       if int(E[len(E)-1]-'0') != dp {
           continue
       }
       // found minimal
       writer.Write(E)
       writer.WriteByte('\n')
       return
   }
   // no solution
   fmt.Fprintln(writer, "Impossible")
}
