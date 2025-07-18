#include<bits/stdc++.h>
using namespace std;

using ll = long long;

int n;

signed main() {
    ios_base::sync_with_stdio(0);
    cin.tie(0);

    cin >> n;
    queue<int> q;
    vector<pair<int, int>> res;
    int r = 0, pre = 0;
    for (int i = 1; i <= n; i++) {
        int a;
        cin >> a;
        if (a == 0)
            continue;

        if (a == 1) {
            if (!q.empty()) {
                res.push_back({q.front(), i});
                q.pop();
                continue;
            }
            else {
                if (pre == 3)
                    res.push_back({r, i});
                res.push_back({++r, i});
            }
        }

        if (a == 2) {
            if (pre == 3)
                res.push_back({r, i});
            res.push_back({++r, i});
            q.push(r);
        }

        if (a == 3) {
            if (pre == 3)
                res.push_back({r, i});
            res.push_back({++r, i});
        }

        pre = a;
    }
    if (!q.empty() || pre == 3) cout << -1;
    else {
        cout << (int) res.size() << '\n';
        for (auto [u, v] : res) {
            cout << u << ' ' << v << '\n';
        }
    }

    return 0;

}