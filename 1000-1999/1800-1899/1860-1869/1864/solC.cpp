#pragma GCC optimize("Ofast")
#pragma GCC target("avx,avx2,fma")
#pragma GCC optimize("unroll-loops")

#include <bits/stdc++.h>
using namespace std;

#define int long long

#define pb push_back
#define mp make_pair
#define fi first
#define se second

void dbg_out() { cerr << endl; }
template <typename H, typename... T>
void dbg_out(H h, T... t) { cerr << ' ' << h; dbg_out(t...); }
#define dbg(...) { cerr << #__VA_ARGS__ << ':'; dbg_out(__VA_ARGS__); }

using ll = long long;
using pi = pair<int,int>;
using pl = pair<ll,ll>;
using vi = vector<int>;
using vl = vector<ll>;

mt19937_64 rng(chrono::steady_clock::now().time_since_epoch().count());

const int MAXN = 2e5 + 5;
const int MOD = 1e9+7; //998244353;
const int INF = 0x3f3f3f3f;
const ll INF64 = ll(4e18) + 5;

void solve(){
	int x;
	cin >> x;
	vector<int> bin;
	int xx = x;
	while(xx > 0){
		bin.pb(xx%2);
		xx /= 2;
	}
	vector<int> ans;
	ans.pb(x);
	int sz = bin.size();
	int at = 1;
	for(int i = 0; i < sz-1; i++){
		if(bin[i]){
			x = x-at;
			ans.pb(x);
		}
		at *= 2;
	}
	while(x != 1){
		x /= 2;
		ans.pb(x);
	}
	cout << ans.size() << '\n';
	for(auto x : ans){
		cout << x << ' ';
	}
	cout << '\n';
}

signed main(){

    ios::sync_with_stdio(false); cin.tie(NULL);

    int t = 1;
    cin >> t;
    while(t--){
        solve();
    }

    return 0;
}