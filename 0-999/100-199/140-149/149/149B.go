package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
   "unicode"
)

func charValue(r rune) int {
   if r >= '0' && r <= '9' {
       return int(r - '0')
   }
   if r >= 'A' && r <= 'Z' {
       return int(r - 'A' + 10)
   }
   return 0
}

// eval evaluates digit values in base r with Horner's method, returns value and if value <= limit
func eval(digits []int, r, limit int64) (int64, bool) {
   var res int64
   for _, d := range digits {
       res = res*r + int64(d)
       if res > limit {
           return res, false
       }
   }
   return res, true
}

func maxRadixForLimit(digits []int, minR int64, limit int64) int64 {
   // if constant polynomial (all leading zeros), unlimited
   allZero := true
   for i := 0; i < len(digits)-1; i++ {
       if digits[i] != 0 {
           allZero = false
           break
       }
   }
   if allZero {
       return -1 // indicates unlimited
   }
   // find high bound
   high := minR
   // increase until exceeds limit
   for {
       if _, ok := eval(digits, high, limit); !ok {
           break
       }
       high *= 2
       // safety cap
       if high > limit+1 {
           high = limit + 1
           break
       }
   }
   // binary search max r in [minR, high) such that eval <= limit
   lo, hi := minR, high
   for lo < hi {
       mid := lo + (hi-lo)/2
       if _, ok := eval(digits, mid, limit); ok {
           lo = mid + 1
       } else {
           hi = mid
       }
   }
   // lo is first invalid, so max valid is lo-1
   return lo - 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, err := reader.ReadString('\n')
   if err != nil && len(line) == 0 {
       return
   }
   line = strings.TrimSpace(line)
   parts := strings.Split(line, ":")
   if len(parts) != 2 {
       return
   }
   aStr, bStr := parts[0], parts[1]
   var aDigits, bDigits []int
   maxDig := 0
   for _, r := range aStr {
       if unicode.IsLower(r) {
           r = unicode.ToUpper(r)
       }
       d := charValue(r)
       aDigits = append(aDigits, d)
       if d > maxDig {
           maxDig = d
       }
   }
   for _, r := range bStr {
       if unicode.IsLower(r) {
           r = unicode.ToUpper(r)
       }
       d := charValue(r)
       bDigits = append(bDigits, d)
       if d > maxDig {
           maxDig = d
       }
   }
   minRadix := int64(max(maxDig+1, 2))
   // check constant for infinite
   aConst := true
   for i := 0; i < len(aDigits)-1; i++ {
       if aDigits[i] != 0 {
           aConst = false
           break
       }
   }
   bConst := true
   for i := 0; i < len(bDigits)-1; i++ {
       if bDigits[i] != 0 {
           bConst = false
           break
       }
   }
   // constant values
   aVal := int64(aDigits[len(aDigits)-1])
   bVal := int64(bDigits[len(bDigits)-1])
   if aConst && bConst {
       if aVal <= 23 && bVal <= 59 {
           fmt.Println(-1)
           return
       }
   }
   // find maximum radix
   maxR := int64(-1)
   // start with unlimited
   // limit A<=23
   rA := maxRadixForLimit(aDigits, minRadix, 23)
   if rA >= 0 {
       maxR = rA
   }
   // limit B<=59
   rB := maxRadixForLimit(bDigits, minRadix, 59)
   if rB >= 0 {
       if maxR < 0 || rB < maxR {
           maxR = rB
       }
   }
   // if both unlimited
   if rA < 0 && rB < 0 {
       // both a and b are constant polynomials but not valid infinite case, no bases
       fmt.Println(0)
       return
   }
   // if one unlimited, use the other
   if maxR < minRadix {
       fmt.Println(0)
       return
   }
   var res []int64
   for r := minRadix; r <= maxR; r++ {
       if av, okA := eval(aDigits, r, 23); okA {
           if bv, okB := eval(bDigits, r, 59); okB {
               if av <= 23 && bv <= 59 {
                   res = append(res, r)
               }
           }
       }
   }
   if len(res) == 0 {
       fmt.Println(0)
       return
   }
   for i, r := range res {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Print(r)
   }
   fmt.Println()
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}
