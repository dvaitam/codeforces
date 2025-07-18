// Problem: D. Card Game
// Contest: Codeforces - Codeforces Round 927 (Div. 3)
// URL: https://codeforces.com/contest/1932/problem/D
// Memory Limit: 256 MB
// Time Limit: 2000 ms
// 
// Powered by CP Editor (https://cpeditor.org)

#include<bits/stdc++.h>
#include <ext/pb_ds/assoc_container.hpp>
using namespace __gnu_pbds;
using namespace std;
typedef long long ll;
#define speed ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0)
#define mp make_pair
#define pb push_back
#define ff first
#define ss second
#define vi vector<int>
#define vll vector<ll> 
#define all(x) (x).begin() , (x).end()
#define inf 1000000000
#define mod 1000000007

void dbg(){
	cerr << endl;
}
template<typename Head , typename... Tail>
void dbg(Head h , Tail... t){
	cerr << h << " ";
	dbg(t...);
}

#ifdef EMBI_DEBUG
#define debug(...) cerr << "(" << #__VA_ARGS__  << "): ", dbg(__VA_ARGS__)
#else 
#define debug(...)
#endif

const int max_n = 1e5 + 9;

typedef tree<int,null_type,less<int>,rb_tree_tag,tree_order_statistics_node_update> indexed_set;
ll power(ll a , ll b)
{
    ll prod = 1;
    while(b)
    {
        if(b&1)
        prod = (prod*a)%mod;
        a = (a*a)%mod;
        b >>= 1;
    }
    return prod;
}
void solve(){
	int n;
	cin >> n;
	
	n *= 2;
	
	char trump;
	cin >> trump;
	
	vector<string> a(n);
	for (int i = 0 ; i < n ; i++) {
		cin >> a[i];
	}
	
	sort(all(a));
	
	map<char, vector<string>> m;
	for (int i = 0 ; i < n ; i++) {
		m[a[i][1]].pb(a[i]);
	}
	
	int rem = 0, max1 = 0;
	vector<pair<string, string>> ans;
	vector<string> max_str;
	vector<string> rem_str;
	for (auto it : m) {
		if (it.ff == trump) {
			max1 = it.ss.size();
			max_str = it.ss;
		} else {
			for (int i = 0 ; i+1 < it.ss.size() ; i += 2) {
				ans.pb({it.ss[i], it.ss[i+1]});
			}
			rem += it.ss.size() % 2;
			if (it.ss.size() % 2) {
				rem_str.pb(it.ss.back());
			}
		}
	}
	
	if (max1 < rem) {
		cout << "IMPOSSIBLE\n";
		return;
	}
	
	if ((max1 - rem) % 2) {
		cout << "IMPOSSIBLE\n";
		return;
	}
	
	int idx = 0;
	
	for (auto it : rem_str) {
		ans.pb({it, max_str[idx++]});
	}
	
	while (idx < max_str.size()) {
		ans.pb({max_str[idx] , max_str[idx+1]});
		idx += 2;
	}
	
	for (auto it : ans) {
		cout << it.ff << " " << it.ss << "\n";
	}
	
}
signed main(){
    int t = 1;
    cin >> t;
    for (int i = 1 ; i <= t ; i++) {
        solve();
    }
}