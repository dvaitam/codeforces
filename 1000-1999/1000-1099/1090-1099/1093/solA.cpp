#include <bits/stdc++.h>
#define endl '\n';
using namespace std;

const int INF = 1e9 + 7;
const int MAX_IDX = 6 + 5;
const int MAX_VAL = 200 + 5;
const int S[MAX_IDX]{2, 3, 4, 5, 6, 7};

int memo[MAX_IDX][MAX_VAL];

int dp(int idx, int val) {
    if (val < 0) return INF;
    if (memo[idx][val] == -1)
        memo[idx][val] = min(1 + dp(idx, val - S[idx]), dp(idx + 1, val));
    return memo[idx][val];
}

signed main()
{
    cout.sync_with_stdio(0);
    cout.tie(0);
    
    int t;
    cin >> t;
    for (int i = 0; i < t; ++i) {
        int x;
        cin >> x;

        for (auto &m : memo) fill_n(m, MAX_VAL, -1);
        fill_n(memo[6] + 1, MAX_VAL, INF);
        memo[6][0] = 0;

        cout << dp(0, x) << endl;
    }
}