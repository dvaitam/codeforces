use std::io::{self, BufRead};
use std::cmp;

fn main() {
    let stdin = io::stdin();
    let mut lines = stdin.lock().lines();

    let first_line = lines.next().unwrap().unwrap();
    let mut iter = first_line.split_whitespace();
    let n: usize = iter.next().unwrap().parse().unwrap();
    let d: i64 = iter.next().unwrap().parse().unwrap();
    let s: i64 = iter.next().unwrap().parse().unwrap();

    let mut high_fi: Vec<i64> = Vec::new();
    let mut low_fi: Vec<i64> = Vec::new();
    let mut num_h0: i64 = 0;
    let mut num_l0: i64 = 0;
    let mut sum_c_easy: i64 = 0;
    let mut total_easy: i64 = 0;
    let mut total_nc: i64 = 0;

    for _ in 0..n {
        let line = lines.next().unwrap().unwrap();
        let mut iter = line.split_whitespace();
        let c: i64 = iter.next().unwrap().parse().unwrap();
        let f: i64 = iter.next().unwrap().parse().unwrap();
        let l: i64 = iter.next().unwrap().parse().unwrap();

        if l >= d {
            if c >= 1 {
                sum_c_easy += c;
                total_easy += 1;
                if f == 0 {
                    num_h0 += 1;
                } else {
                    high_fi.push(f);
                }
            } else {
                total_nc += 1;
                if f == 0 {
                    num_l0 += 1;
                } else {
                    low_fi.push(f);
                }
            }
        } else {
            if c >= 1 {
                sum_c_easy += c;
                total_easy += 1;
            } else {
                total_nc += 1;
            }
        }
    }

    let l_val: i64 = sum_c_easy - total_easy;

    high_fi.sort();
    low_fi.sort();

    let m_plus = high_fi.len();
    let nn_plus = low_fi.len();

    let mut ps_h: Vec<i64> = vec![0; m_plus + 1];
    for j in 1..=m_plus {
        ps_h[j] = ps_h[j - 1] + high_fi[j - 1];
    }

    let mut ps_l: Vec<i64> = vec![0; nn_plus + 1];
    for j in 1..=nn_plus {
        ps_l[j] = ps_l[j - 1] + low_fi[j - 1];
    }

    // case1
    let mut fuel1: i64 = 0;
    let mut max_c: i64 = num_l0;
    if nn_plus > 0 {
        let j = get_max_num(&ps_l, s);
        max_c += j as i64;
        fuel1 = ps_l[j];
    }
    let m1: i64 = max_c;

    // case2
    let mut m2: i64 = 0;
    let mut fuel2: i64 = 0;
    if num_h0 + m_plus as i64 > 0 {
        let r_free: i64 = num_h0 + num_l0;
        let h_free: i64 = num_h0;
        let mut max_add: i64 = 0;

        if h_free > 0 {
            let mut all_plus: Vec<i64> = high_fi.clone();
            all_plus.extend_from_slice(&low_fi);
            all_plus.sort();
            let all_size = all_plus.len();
            let mut ps_all: Vec<i64> = vec![0; all_size + 1];
            for j in 1..=all_size {
                ps_all[j] = ps_all[j - 1] + all_plus[j - 1];
            }
            let j = get_max_num(&ps_all, s);
            max_add = j as i64;
        } else {
            for hpp in 1..=m_plus {
                let cost_h = ps_h[hpp];
                if cost_h > s {
                    continue;
                }
                let rem = s - cost_h;
                let c = get_max_num(&ps_l, rem) as i64;
                let add = hpp as i64 + c;
                if add > max_add {
                    max_add = add;
                }
            }
        }

        let r: i64 = r_free + max_add;
        let w: i64 = cmp::min(total_nc, r + l_val);
        m2 = total_easy + w;

        // compute fuel2
        let is_full = w == total_nc;
        let minhh = if h_free == 0 { 1 } else { 0 };
        let mut min_cost: i64 = i64::MAX / 2;

        if is_full {
            let target_r = cmp::max(0, total_nc - l_val);
            let target_add = cmp::max(0, target_r - r_free);
            for hpp in minhh..=m_plus as i64 {
                if hpp > m_plus as i64 {
                    break;
                }
                let t = cmp::max(0, target_add - hpp);
                if t > nn_plus as i64 {
                    continue;
                }
                let cost_low = if t == 0 { 0 } else { ps_l[t as usize] };
                let cost = ps_h[hpp as usize] + cost_low;
                if cost < min_cost {
                    min_cost = cost;
                }
            }
        } else {
            let target_add = max_add;
            for hpp in minhh..=m_plus as i64 {
                if hpp > m_plus as i64 {
                    break;
                }
                let t = target_add - hpp;
                if t < 0 || t > nn_plus as i64 {
                    continue;
                }
                let cost = ps_h[hpp as usize] + ps_l[t as usize];
                if cost < min_cost {
                    min_cost = cost;
                }
            }
        }
        if min_cost <= s {
            fuel2 = min_cost;
        } else {
            fuel2 = 0;
            m2 = 0;
        }
    }

    let max_m_val = cmp::max(m1, m2);
    if max_m_val == 0 {
        println!("0 0");
        return;
    }

    let mut min_fuel: i64 = i64::MAX;
    if m1 == max_m_val {
        min_fuel = cmp::min(min_fuel, fuel1);
    }
    if m2 == max_m_val {
        min_fuel = cmp::min(min_fuel, fuel2);
    }

    println!("{} {}", max_m_val, min_fuel);
}

fn get_max_num(ps: &Vec<i64>, lim: i64) -> usize {
    if lim < 0 {
        return 0;
    }
    let mut left: usize = 0;
    let mut right: usize = ps.len() - 1;
    let mut res: usize = 0;
    while left <= right {
        let mid = left + (right - left) / 2;
        if ps[mid] <= lim {
            res = mid;
            left = mid + 1;
        } else {
            if mid == 0 {
                break;
            }
            right = mid - 1;
        }
    }
    res
}

