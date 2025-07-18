#include <bits/stdc++.h>



using i64 = long long;



void solve() {

    i64 a, b, d;

    std::cin >> a >> b >> d;

    

    if (__builtin_ctz(a | b) < __builtin_ctz(d)) {

        std::cout << -1 << "\n";

        return;

    }

    

    int k = __builtin_ctz(d);

    i64 x = 0;

    

    for (int i = k; i < 30; i++) {

        if (~x >> i & 1) {

            x += d << (i - k);

        }

    }

    

     std::cout << x << "\n";

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