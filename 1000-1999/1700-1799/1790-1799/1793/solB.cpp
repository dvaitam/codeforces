#include <bits/stdc++.h>
#define EslamAhmed171 ios::sync_with_stdio(0); cin.tie(0); cout.tie(0);
#define endl "\n"
#define pp push_back
#define ll long long
#define ld long double
#define ull unsigned long long
#define pii pair<int,int>
#define pll pair<ll,ll>
#define vi vector<int>
#define vii vector<pii>
#define vvi vector<vi>
#define all(x) (x).begin(), (x).end()
#define allr(x) (x).rbegin(), (x).rend()
#define fi first
#define se second
const int N = 2e5 + 10;
const int MAX_INF = 2e9 + 10;
const int MIN_INF = -2e9 - 10;
const ll MOD = 1e9 + 7;
const ld PI = 3.141592653589793238;
using namespace std;
ll pw(ll x, ll y) {
    ll result = 1;
    while (y){
        if (y & 1LL)
            result = (result * x);
        x = (x * x);
        y >>= 1;
    }
    return result;
}
ll gcd(ll a, ll b){
    if (b == 0)
        return a;
    return gcd(b,a % b);
}
ll lcm (ll a, ll b){
    // lcm(a, b) = a * b / gcd(a,b)
    return a / gcd(a,b) * b;
}
int dx[] = {+0, +0, -1, +1, +1, +1, -1, -1};
int dy[] = {-1, +1, +0, +0, +1, -1, +1, -1};
void solve () {
    int x, y; cin >> x >> y;
    int n = 2 * abs(x - y);
    cout << n << endl;
    for (int i = y; i <= x; i++){
        cout << i << " ";
    }
    for (int i = x - 1; i > y; i--){
        cout << i << " ";
    }
    cout << endl;
}
signed main(){
    EslamAhmed171
//    freopen("input.txt", "r", stdin);
//    freopen("output.txt", "w", stdout);
    int tt = 1; cin >> tt;
    while (tt--)
        solve();
    return 0;
}