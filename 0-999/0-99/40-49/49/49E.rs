use std::io::{self, Read};

fn compute_masks(s: &Vec<usize>, rules: &Vec<(usize, usize, usize)>) -> Vec<Vec<u32>> {
    let n = s.len();
    let mut mask = vec![vec![0u32; n]; n];
    for i in 0..n {
        mask[i][i] = 1u32 << s[i];
    }
    for len in 2..=n {
        for i in 0..=n - len {
            let j = i + len - 1;
            let mut m = 0u32;
            for k in i..j {
                let left = mask[i][k];
                let right = mask[k + 1][j];
                if left == 0 || right == 0 {
                    continue;
                }
                for &(a, b, c) in rules.iter() {
                    if ((left >> b) & 1) != 0 && ((right >> c) & 1) != 0 {
                        m |= 1u32 << a;
                    }
                }
            }
            mask[i][j] = m;
        }
    }
    mask
}

fn main() {
    let mut input = String::new();
    io::stdin().read_to_string(&mut input).unwrap();
    let mut lines = input.lines().filter(|l| !l.trim().is_empty());

    let s1_str = lines.next().unwrap().trim().to_string();
    let s2_str = lines.next().unwrap().trim().to_string();
    let n: usize = lines.next().unwrap().trim().parse().unwrap();

    let mut rules: Vec<(usize, usize, usize)> = Vec::new();
    for _ in 0..n {
        if let Some(line) = lines.next() {
            let letters: Vec<usize> = line
                .chars()
                .filter(|c| c.is_ascii_lowercase())
                .map(|c| (c as u8 - b'a') as usize)
                .collect();
            if letters.len() >= 3 {
                rules.push((letters[0], letters[1], letters[2]));
            }
        }
    }

    let s1: Vec<usize> = s1_str.bytes().map(|b| (b - b'a') as usize).collect();
    let s2: Vec<usize> = s2_str.bytes().map(|b| (b - b'a') as usize).collect();
    let n1 = s1.len();
    let n2 = s2.len();

    let masks1 = compute_masks(&s1, &rules);
    let masks2 = compute_masks(&s2, &rules);

    let inf = i32::MAX / 4;
    let mut dp = vec![vec![inf; n2 + 1]; n1 + 1];
    dp[n1][n2] = 0;

    for i in (0..=n1).rev() {
        for j in (0..=n2).rev() {
            if i == n1 && j == n2 {
                continue;
            }
            if i == n1 || j == n2 {
                continue;
            }
            let mut best = inf;
            for k1 in i..n1 {
                for k2 in j..n2 {
                    if (masks1[i][k1] & masks2[j][k2]) != 0 {
                        let next = dp[k1 + 1][k2 + 1];
                        if next != inf && next + 1 < best {
                            best = next + 1;
                        }
                    }
                }
            }
            dp[i][j] = best;
        }
    }

    if dp[0][0] >= inf {
        println!("-1");
    } else {
        println!("{}", dp[0][0]);
    }
}
