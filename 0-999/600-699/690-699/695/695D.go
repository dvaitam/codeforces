package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func nextTime(start int64, sec, min, hour, dow, dom, month int) int64 {
	t := start + 1
	for {
		tm := time.Unix(t, 0).UTC()
		y, m, d := tm.Date()
		h, mi, s := tm.Clock()
		wd := (int(tm.Weekday())+6)%7 + 1

		if month != -1 && int(m) != month {
			if int(m) < month {
				tm = time.Date(y, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			} else {
				tm = time.Date(y+1, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			}
			t = tm.Unix()
			continue
		}

		domOk := dom == -1 || d == dom
		dowOk := dow == -1 || wd == dow
		validDay := true
		if dom != -1 && dow != -1 {
			validDay = domOk || dowOk
		} else {
			if dom != -1 && !domOk {
				validDay = false
			}
			if dow != -1 && !dowOk {
				validDay = false
			}
		}
		if !validDay {
			tm = time.Date(y, m, d, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
			t = tm.Unix()
			continue
		}

		if hour != -1 && h != hour {
			if h < hour {
				tm = time.Date(y, m, d, hour, 0, 0, 0, time.UTC)
			} else {
				tm = time.Date(y, m, d, hour, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
			}
			t = tm.Unix()
			continue
		}

		if min != -1 && mi != min {
			if mi < min {
				tm = time.Date(y, m, d, h, min, 0, 0, time.UTC)
			} else {
				tm = time.Date(y, m, d, h, min, 0, 0, time.UTC).Add(time.Hour)
			}
			t = tm.Unix()
			continue
		}

		if sec != -1 && s != sec {
			if s < sec {
				tm = time.Date(y, m, d, h, mi, sec, 0, time.UTC)
			} else {
				tm = time.Date(y, m, d, h, mi, sec, 0, time.UTC).Add(time.Minute)
			}
			t = tm.Unix()
			continue
		}

		return t
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s, m, h, dow, dom, mo int
	if _, err := fmt.Fscan(reader, &s, &m, &h, &dow, &dom, &mo); err != nil {
		return
	}
	var n int
	fmt.Fscan(reader, &n)
	for i := 0; i < n; i++ {
		var t int64
		fmt.Fscan(reader, &t)
		ans := nextTime(t, s, m, h, dow, dom, mo)
		fmt.Fprintln(writer, ans)
	}
}
