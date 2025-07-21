#include<bits/stdc++.h>
using namespace std;
const int N = 1e6 + 5;
long long n, m, i, a[N], b, c, d[N], dp[N], x, ans;
int main(){
    ios::sync_with_stdio(0),cout.tie(0),cin.tie(0);
    memset(d, '?', sizeof d);
    for (cin >> n >> m; i < n; i++) cin >> a[i];
    for (i = 0; i < n; i++) {
        cin >> b;
        d[a[i]] = min(d[a[i]], a[i] - b);
    }
    for (i = 1; i < N; i++){
        d[i] = min(d[i], d[i - 1]);
        dp[i] = (d[i] <= i ? 2 + dp[i - d[i]] : 0);
    }
    while (m--){
        cin >> c;
        if (c >= N){
            x = (c - N) / d[N - 1] + 1;
            c -= x * d[N - 1];
            ans += 2 * x;
        }
        ans += dp[c];
    }
    cout << ans;
    return 0;
}
