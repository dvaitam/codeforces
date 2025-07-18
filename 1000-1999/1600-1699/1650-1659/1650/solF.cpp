#include <bits/stdc++.h>

// #define int long long

#define x first
#define y second

using namespace std;

typedef long long LL;

void solve()
{
    int n, m;
    cin >> n >> m;

    vector<int> lim(n);
    for (int i = 0; i < n; i ++ ) cin >> lim[i];

    vector<vector<array<int, 3>>> a(n, vector<array<int, 3>>(1));
    for (int i = 0; i < m; i ++ ) {
        int e, t, p;
        cin >> e >> t >> p;
        e -- ;
        a[e].push_back({t, p, i});
    }

    vector<int> ans;
    int cur = 0;
    for (int i = 0; i < n; i ++ ) {
        int N = a[i].size();
        vector<vector<int>> dp(N, vector<int>(101, 1e9 + 1));
        for (int j = 0; j < N; j ++ ) dp[j][0] = 0;
        N -- ;
        for (int j = 1; j <= N; j ++ ) {
            auto [t, p, _] = a[i][j];
            for (int k = 0; k <= 100; k ++ ) {
                dp[j][k] = dp[j - 1][k];
                dp[j][k] = min(dp[j][k], dp[j - 1][max(0, k - p)] + t);
            }
        }

        // cout << dp[N][100] << "\n";
        // return;

        cur += dp[N][100];

        if (cur > lim[i]) {
            cout << "-1\n";
            return;
        }

        int k = 100;
        for (int j = N; j; j -- ) {
            auto [t, p, id] = a[i][j];
            if (k >= 0 && dp[j][k] == dp[j - 1][max(0, k - p)] + t) {
                ans.push_back(id);
                k -= p;
            }
        }
    }

    cout << ans.size() << "\n";
    for (auto id : ans) cout << id + 1 << " ";
    cout << "\n";
}

signed main()
{
    std::ios::sync_with_stdio(false);
    std::cin.tie(nullptr);

    int T;
    cin >> T;

    while (T -- ) {
        solve();
    }

    return 0;
}