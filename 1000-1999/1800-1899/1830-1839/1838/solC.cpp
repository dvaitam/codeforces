#include <bits/stdc++.h>
#define buff ios::sync_with_stdio(false), cin.tie(nullptr), cout.tie(nullptr);
//#define endl '\n'
#define pb push_back
#define eb emplace_back
#define all(x) x.begin(), x.end()
#define PI 3.1415926535
#define int long long
#define lowbit(x) ((x) & - (x))
/*
#pragma GCC optimize ("Ofast")
#pragma GCC optimize ("unroll-loops")
#pragma GCC optimize(3)
*/
using namespace std;
typedef long long ll;
typedef pair<ll, ll> pll;
typedef vector<ll> vll;
typedef vector<pair<ll, ll>> vpll;
typedef vector<string> vstr;
const ll MAX_INT = 0x3f3f3f3f;
const ll MAX_LL = 0x3f3f3f3f3f3f3f3f;
const ll CF = 2e5 + 9;
const ll N = 2e6 + 5;
const ll mod = 998244353;
int a[1005][1005];
void solve(){
	int n , m; cin >> n >> m;
	int now = 1;
	for(int i = 1;i <= n;i++) {
		for(int j = 1;j <= m;j++) {
			a[i][j] = now++;
		}
	}
	for(int i = 1;i <= n;i++) {
		if(i % 2 == 0) {
			for(int j = 1;j <= m;j++) cout << a[i][j] << " ";
			cout << "\n";
		}
	}
	for(int i = 1;i <= n;i++) {
		if(i & 1) {
			for(int j = 1;j <= m;j++) cout << a[i][j] << " ";
			cout << "\n";
		}
	}
}	
signed main()
{
	buff;
	int t = 1;
	cin >> t;
	while(t--) 
	{
		solve();
	}
    return 0; 
}
//