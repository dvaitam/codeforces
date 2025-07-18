#include <bits/stdc++.h>
#define buff ios::sync_with_stdio(false), cin.tie(nullptr), cout.tie(nullptr);
#define endl '\n'
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
const ll CF = 3e5 + 9;
const ll mod = 1e9 + 7;
struct q{
	int id,val;
} a[CF];
bool cmp(q x,q y) {
	return x.val > y.val;
}
bool cmpid(q x,q y) {
	return x.id < y.id;
}
void solve(){
	int n; cin >> n;
	for(int i = 1;i <= n;i++) cin >> a[i].val,a[i].id = i;
	sort(a + 1,a + n + 1,cmp);
	for(int i = 1;i <= n;i++) a[i].val = i;
	sort(a + 1,a + n + 1,cmpid);
	for(int i = 1;i <= n;i++) cout << a[i].val << " ";
	cout << "\n";
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