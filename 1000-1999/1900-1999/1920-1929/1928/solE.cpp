#include<bits/stdc++.h>
#include<ext/pb_ds/assoc_container.hpp>
#include<ext/pb_ds/tree_policy.hpp>
using namespace std;
using namespace __gnu_pbds;
#define ll long long
#define ld long double
#define int ll
#define ln '\n'
#define F first 
#define S second
#define pb push_back
#define ins insert 
#define all(v) v.begin(), v.end()
#define alr(v) v.rbegin(), v.rend()
#define whole(a, n) a + 1, a + n + 1
#define pii pair<int, int>
#define pll pair<ll, ll>
#define oset tree<int, null_type, less<int>, rb_tree_tag, tree_order_statistics_node_update>
const int sz = 2e5 + 5;
const int inf = 1e9 + 7;
const ll infll = 1e18 + 7;
int get(int x, int y){
    int p1 = (x + y) * (x + y + 1) / 2;
    int p2 = (x - 1) * x / 2;
    return p1 - p2;
}
void solve(){
    int n, x, y, s;
    cin >> n >> x >> y >> s;
    int p = x % y;
    if(s < p * n or (s - p * n) % y){ cout << "NO" << ln; return;}
    s = (s - p * n) / y;
    vector<int>dp(s + 1, inf);
    vector<pii>par(s + 1, {0, 0});
    dp[0] = 0;
    for(int i = 1; i * (i - 1) / 2 <= s; i++){
        int k = i * (i - 1) / 2;
        for(int j = k; j <= s; j++){
            if(dp[j-k] + i <= dp[j]){
                dp[j] = dp[j-k] + i;
                par[j] = {j - k, i};
            }
        }
    }
    int be = x / y;
    for(int i = 0; i < n and get(be, i) <= s; i++){
        int t = get(be, i);
        if(dp[s-t] < n - i){
            vector<int>v;
            for(int j = 0; j <= i; j++) v.pb(be + j);
            int k = s - t;
            while(k){
                auto [l1, l2] = par[k];
                for(int g = 0; g < l2; g++) v.pb(g);
                k = l1;
            }
            while(v.size() < n) v.pb(0);
            cout << "YES" << ln;
            for(auto g: v) cout << g * y + p << ' ';
            cout << ln;
            return;
        }
    }
    cout << "NO" << ln;
}
signed main(){
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    int t; cin >> t;
    while(t--) solve();
    return 0;
}