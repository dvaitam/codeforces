#include <bits/stdc++.h> 
using namespace std;
#define endl '\n'
#define int long long
#define pii pair<int,int>
#define vi vector<int>
#define vvi vector<vi>
#define vpii vector<pii>
#define vvpii vector<vpii>
#define F first
#define S second
#define pb push_back
#define all(v) v.begin(),v.end()
#define loop(i,a,b,d) for(int i=a; i<b; i+=d)
#define fr(i,a,b) for(int i=a; i<b; i++)
#define rep(i,n) for(int i=0; i<n; i++)
#define ret(msg) {cout << msg << endl; return;}
inline int popcnt (int x) { return __builtin_popcountll(x); }
#define mod 1000000007
#define MOD 998244353

void solve() {
    int n;
    cin >> n;
    
    int m = (n*(n-1))/2;
    vi a(m);
    rep(i,m) {
        cin >> a[i];
    }

    sort(all(a));
    int k = 0;
    fr(i,1,n) {
        cout << a[k] << " ";
        k += n-i;
    }
    ret(a[m-1])
}

int32_t main() {
    ios_base::sync_with_stdio(false); cin.tie(NULL);
    int t;
    cin >> t;
    while(t--) {
        solve();
    }
}