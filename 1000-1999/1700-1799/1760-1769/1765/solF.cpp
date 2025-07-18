#include <bits/stdc++.h>
using namespace std;
using ll = long long;
using ld = long double;
using pii = pair<int, int>;

const int N = 5050;

struct pt {
    int x;
    int contract;
    int y;
};
pt p[N];
ll dp[N];

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(nullptr);
    int n, k;
    cin >> n >> k;
    for (int i = 0; i < n; ++i) {
        cin >> p[i].x >> p[i].contract >> p[i].y;
    }
    sort(p, p + n, [](const pt& a, const pt& b) -> bool {
        if (a.x != b.x)
            return a.x < b.x;
        return a.y < b.y;
    });
    for (int i = 0; i < n; ++i) {
        dp[i] = 0;//p[i].x * (ll)p[i].y * (ll)k;
        for (int j = 0; j < i; ++j) {
            dp[i] = max(dp[i], dp[j] + (p[i].x - p[j].x) * (ll)(p[i].y + p[j].y) * (ll)k);
        }
        dp[i] -= p[i].contract * (ll)200;
    }
    ll ans = 0;
    for (int i = 0; i < n; ++i) {
        ans = max(ans, dp[i] /*+ (100 - p[i].x) * (ll)p[i].y * (ll)k*/);
    }
    cout << setprecision(17) << fixed << ans / (ld)200 << "\n";


}