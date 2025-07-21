package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
   "strconv"
   "strings"
)

// compute s1(n) = n*(n+1)/2
func s1(n int64) *big.Int {
   if n <= 0 {
       return big.NewInt(0)
   }
   a := big.NewInt(n)
   b := big.NewInt(n + 1)
   a.Mul(a, b)
   a.Div(a, big.NewInt(2))
   return a
}

// compute s2(n) = n*(n+1)*(2n+1)/6
func s2(n int64) *big.Int {
   if n <= 0 {
       return big.NewInt(0)
   }
   a := big.NewInt(n)
   b := big.NewInt(n + 1)
   c := big.NewInt(2*n + 1)
   a.Mul(a, b)
   a.Mul(a, c)
   a.Div(a, big.NewInt(6))
   return a
}

// compute s3(n) = (n*(n+1)/2)^2
func s3(n int64) *big.Int {
   t := s1(n)
   t.Mul(t, t)
   return t
}

// compute sum of a_{i,j} for 1<=i<=x,1<=j<=y
func sumXY(x, y int64) *big.Int {
   zero := big.NewInt(0)
   if x <= 0 || y <= 0 {
       return zero
   }
   // T1: sum over i>=j
   A := x
   if y < A {
       A = y
   }
   // Part1: i=1..A
   S1A := s1(A)
   S2A := s2(A)
   S3A := s3(A)
   // P1a = S3A - 2*S2A + S1A
   P1a := new(big.Int).Set(S3A)
   tmp := new(big.Int).Mul(big.NewInt(2), S2A)
   P1a.Sub(P1a, tmp)
   P1a.Add(P1a, S1A)
   // P1b = sum_{i=1..A} i*(i+1)/2 = A*(A+1)*(A+2)/6
   P1b := big.NewInt(A)
   P1b.Mul(P1b, big.NewInt(A+1))
   P1b.Mul(P1b, big.NewInt(A+2))
   P1b.Div(P1b, big.NewInt(6))
   T1 := new(big.Int).Add(P1a, P1b)
   // Part2: i=A+1..x, k=y
   B := x - A
   if B > 0 {
       // sum_{u=A..x-1} u^2 = s2(x-1) - s2(A-1)
       S2x1 := s2(x - 1)
       S2A1 := s2(A - 1)
       sumU2 := new(big.Int).Sub(S2x1, S2A1)
       // Part2a = y * sumU2
       part2a := new(big.Int).Mul(sumU2, big.NewInt(y))
       T1.Add(T1, part2a)
       // Part2b = B * (y*(y+1)/2)
       S1y := s1(y)
       part2b := new(big.Int).Mul(big.NewInt(B), S1y)
       T1.Add(T1, part2b)
   }
   // T2: sum over i<j
   // Part1: j=1..min(y,x+1)
   C := y
   if x+1 < C {
       C = x + 1
   }
   // N = C-1 for k from 0..N
   N := C - 1
   S1N := s1(N)
   S2N := s2(N)
   S3N := s3(N)
   // numerator = 2*S3N + 3*S2N +3*S1N
   num := new(big.Int).Mul(big.NewInt(2), S3N)
   tmp2 := new(big.Int).Mul(big.NewInt(3), S2N)
   num.Add(num, tmp2)
   tmp2.Mul(big.NewInt(3), S1N)
   num.Add(num, tmp2)
   // T2_part1 = num /2
   T2 := new(big.Int).Div(num, big.NewInt(2))
   // Part2: j=C+1..y, m=x
   D := y - C
   if D > 0 {
       // sum2 = sum_{u=C..y-1} u^2 = s2(y-1) - s2(C-1)
       S2y1 := s2(y - 1)
       S2C1 := s2(C - 1)
       sum2 := new(big.Int).Sub(S2y1, S2C1)
       // sum1 = sum_{j=C+1..y} j = s1(y) - s1(C)
       S1y := s1(y)
       S1C := s1(C)
       sum1 := new(big.Int).Sub(S1y, S1C)
       // sumX = x * (sum2 + 2*sum1)
       tmp3 := new(big.Int).Mul(big.NewInt(2), sum1)
       tmp3.Add(tmp3, sum2)
       sumX := new(big.Int).Mul(tmp3, big.NewInt(x))
       T2.Add(T2, sumX)
       // subtract D * x*(x+1)/2
       sx := new(big.Int).Mul(big.NewInt(x), big.NewInt(x+1))
       sx.Div(sx, big.NewInt(2))
       sub := new(big.Int).Mul(big.NewInt(D), sx)
       T2.Sub(T2, sub)
   }
   // total = T1 + T2
   total := new(big.Int).Add(T1, T2)
   return total
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // read t
   line, _ := reader.ReadString('\n')
   t, _ := strconv.Atoi(line[:len(line)-1])
   mod := big.NewInt(10000000000)
   for i := 0; i < t; i++ {
       s, _ := reader.ReadString('\n')
       tokens := strings.Fields(s)
       x1, _ := strconv.ParseInt(tokens[0], 10, 64)
       y1, _ := strconv.ParseInt(tokens[1], 10, 64)
       x2, _ := strconv.ParseInt(tokens[2], 10, 64)
       y2, _ := strconv.ParseInt(tokens[3], 10, 64)
       // compute raw sum
       s22 := sumXY(x2, y2)
       s12 := sumXY(x1-1, y2)
       s21 := sumXY(x2, y1-1)
       s11 := sumXY(x1-1, y1-1)
       raw := new(big.Int).Sub(s22, s12)
       raw.Sub(raw, s21)
       raw.Add(raw, s11)
       // raw < mod?
       if raw.Cmp(mod) < 0 {
           fmt.Fprintln(writer, raw.String())
       } else {
           // print ... + last 10 digits
           m1 := new(big.Int).Mod(s22, mod)
           m2 := new(big.Int).Mod(s12, mod)
           m3 := new(big.Int).Mod(s21, mod)
           m4 := new(big.Int).Mod(s11, mod)
           temp := new(big.Int).Sub(m1, m2)
           temp.Sub(temp, m3)
           temp.Add(temp, m4)
           temp.Mod(temp, mod)
           // format with leading zeros to 10 digits
           txt := temp.String()
           if len(txt) < 10 {
               txt = strings.Repeat("0", 10-len(txt)) + txt
           }
           // last 10 characters
           if len(txt) > 10 {
               txt = txt[len(txt)-10:]
           }
           fmt.Fprintln(writer, "..."+txt)
       }
   }
}
