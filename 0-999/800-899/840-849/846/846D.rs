use std::io::{self, Read};

fn exists_square(limit: u64, times: &[Vec<u64>], n: usize, m: usize, k: usize) -> bool {
    let mut pref = vec![vec![0i32; m + 1]; n + 1];
    for i in 1..=n {
        for j in 1..=m {
            let broken = if times[i - 1][j - 1] <= limit { 1 } else { 0 };
            pref[i][j] = pref[i - 1][j] + pref[i][j - 1] - pref[i - 1][j - 1] + broken;
        }
    }
    for i in k..=n {
        for j in k..=m {
            let total = pref[i][j] - pref[i - k][j] - pref[i][j - k] + pref[i - k][j - k];
            if total as usize == k * k {
                return true;
            }
        }
    }
    false
}

fn main() {
    let mut input = String::new();
    io::stdin().read_to_string(&mut input).unwrap();
    let mut it = input.split_whitespace();

    let n: usize = it.next().unwrap().parse().unwrap();
    let m: usize = it.next().unwrap().parse().unwrap();
    let k: usize = it.next().unwrap().parse().unwrap();
    let q: usize = it.next().unwrap().parse().unwrap();

    if q == 0 {
        println!("-1");
        return;
    }

    let mut times = vec![vec![u64::MAX; m]; n];
    let mut moments = Vec::with_capacity(q);

    for _ in 0..q {
        let x: usize = it.next().unwrap().parse().unwrap();
        let y: usize = it.next().unwrap().parse().unwrap();
        let t: u64 = it.next().unwrap().parse().unwrap();
        let cell = &mut times[x - 1][y - 1];
        if t < *cell {
            *cell = t;
        }
        moments.push(t);
    }

    moments.sort();

    if !exists_square(*moments.last().unwrap(), &times, n, m, k) {
        println!("-1");
        return;
    }

    let mut l = 0usize;
    let mut r = moments.len() - 1;
    let mut answer = moments[r];

    while l <= r {
        let mid = (l + r) / 2;
        if exists_square(moments[mid], &times, n, m, k) {
            answer = moments[mid];
            if mid == 0 {
                break;
            }
            r = mid - 1;
        } else {
            l = mid + 1;
        }
    }

    println!("{}", answer);
}
