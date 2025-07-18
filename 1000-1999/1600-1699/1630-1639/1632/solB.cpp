#include <bits/stdc++.h>

using namespace std;

#ifndef DEBUG
    #define cerr if (false) cerr
#endif

int min_ = INT32_MAX;
vector <pair <vector <int>, int>> all;

void rec(vector <int> a, vector <bool> used) {
    if (a.size() == used.size()) {
        int ans = -INT32_MAX;
        for (int i = 1; i < int(used.size()); ++i) {
            ans = ::max(a[i - 1] ^ a[i], ans);
        }
        min_ = ::min(ans, min_);
        all.push_back({a, ans});
    } else {
        for (int i = 0; i < int(used.size()); ++i) {
            if (!used[i]) {
                a.push_back(i);
                used[i] = true;
                rec(a, used);
                a.pop_back();
                used[i] = false;
            }
        }
    }
}

signed main() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);
    int t;
    cin >> t;
    while (t--){
        int n;
        cin >> n;
        --n;
        int d = 1 << int(log2(n));
        while (n >= d){
            cout << n << ' ';
            --n;
        }
        --d;
        cerr << d << '\n';
        cout << 0 << ' ';
        while (d > 0){
            cout << d << ' ';
            --d;
        }
        cout << '\n';
    }
}