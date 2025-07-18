#include <bits/stdc++.h>
 
using i64 = long long;
 
void solve() {
    int n;
    std::cin >> n;
    
    int l = n / 2 + 1, r = n;
    int t = 1;
    std::vector<std::array<int, 3>> ans;
    int lastv = 0;
    int lastt = 0;
    int lastl = 0;
    std::vector<int> vis(n + 1);
    while (true) {
        if (r <= 11) {
            if (lastv != 0) {
                ans.push_back({lastt, 2 * lastt, lastv});
            }
            if (r == 3) {
                ans.push_back({t, 2 * t, 3 * t});
            } else if (r == 4) {
                ans.push_back({t, 3 * t, 4 * t});
            } else if (r == 5) {
                ans.push_back({3 * t, 4 * t, 5 * t});
            } else if (r == 6) {
                ans.push_back({3 * t, 4 * t, 5 * t});
                vis[6 * t] = 1;
            } else if (r == 7) {
                ans.push_back({t, 3 * t, 4 * t});
                ans.push_back({5 * t, 6 * t, 7 * t});
            } else if (r == 8) {
                ans.push_back({t, 5 * t, 7 * t});
                ans.push_back({2 * t, 6 * t, 8 * t});
            } else if (r == 9) {
                ans.push_back({t, 5 * t, 6 * t});
                ans.push_back({7 * t, 8 * t, 9 * t});
            } else if (r == 10) {
                ans.push_back({2 * t, 6 * t, 10 * t});
                ans.push_back({7 * t, 8 * t, 9 * t});
            } else if (r == 11) {
                ans.push_back({t, 10 * t, 11 * t});
                ans.push_back({7 * t, 8 * t, 9 * t});
                vis[6 * t] = 1;
            }
            break;
        }
        int skip = 0;
        if (l % 4 == 3) {
            if (r % 4 == 1) {
                assert(r == 2 * l - 1);
                ans.push_back({l * t, r * t, 2 * t});
                l++;
                skip = 1;
            }
        }
        while (l % 4 > 1) {
            l--;
        }
        int i = l / 4 * 4 + 1;
        while (i + 2 <= r) {
            ans.push_back({i * t, (i + 1) * t, (i + 2) * t});
            i += 4;
        }
        if (i <= r) {
            if (i + 1 <= r) {
                ans.push_back({t, i * t, (i + 1) * t});
            } else {
                if (!skip) {
                    if (lastv != 0) {
                        if (std::gcd(i * t, lastv) < lastl * lastt) {
                            ans.push_back({std::gcd(i * t, lastv), i * t, lastv});
                            lastv = 0;
                            lastt = 0;
                        } else {
                            ans.push_back({lastt, 2 * lastt, lastv});
                            lastv = i * t;
                            lastt = t;
                            lastl = l;
                        }
                    } else {
                        lastv = i * t;
                        lastt = t;
                        lastl = l;
                    }
                }
            }
        }
        l = (l - 1) / 4 + 1;
        r = r / 4;
        t *= 4;
    }
    
    assert(ans.size() <= n / 6 + 5);
    std::cout << ans.size() << "\n";
    std::vector<int> a(n + 1);
    std::iota(a.begin(), a.end(), 0);
    for (auto [i, j, k] : ans) {
        a[i] = std::lcm(j, k);
        a[j] = std::lcm(i, k);
        a[k] = std::lcm(i, j);
        assert(!vis[i]);
        vis[i] = 1;
        assert(!vis[j]);
        vis[j] = 1;
        assert(!vis[k]);
        vis[k] = 1;
        int g = std::gcd(std::gcd(i, j), k);
        assert(std::gcd(i, j) == g);
        assert(std::gcd(i, k) == g);
        assert(std::gcd(j, k) == g);
        std::cout << i << " " << j << " " << k << "\n";
    }
    for (int i = n / 2 + 1; i <= n; i++) {
        if (!vis[i]) {
            std::cerr << "n : " << n << "\n";
            std::exit(0);
        }
        assert(vis[i]);
    }
    // for (int i = 1; i <= n; i++) {
    //     int g = 0;
    //     int cnt = 0;
    //     for (int j = 1; j <= n; j++) {
    //         if (a[j] % i == 0) {
    //             g = std::gcd(g, a[j]);
    //             cnt++;
    //         }
    //     }
    //     if (cnt > 1 && g == i) {
    //         continue;
    //     }
    //     std::cerr << "i : " << i << "\n";
    // }
}
 
int main() {
    std::ios::sync_with_stdio(false);
    std::cin.tie(nullptr);
    
    int t;
    std::cin >> t;
    
    while (t--) {
        solve();
    }
    
    return 0;
}