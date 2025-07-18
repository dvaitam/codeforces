#include <bits/stdc++.h>
#pragma comment(linker, "/stack:200000000")
#pragma GCC optimize ("Ofast")
#pragma GCC target ("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx,tune=native")

#include <bits/stdc++.h>
using namespace std;

#define rep(i, n) for(int i = 0; i < (n); ++i)
#define repA(i, a, n) for(int i = a; i <= (n); ++i)
#define repD(i, a, n) for(int i = a; i >= (n); --i)
#define trav(a, x) for(auto& a : x)
#define all(x) x.begin(), x.end()
#define sz(x) (int)(x).size()
#define fill(a) memset(a, 0, sizeof (a))
#define fst first
#define snd second
#define mp make_pair
#define pb push_back
typedef long long ll;
typedef pair<int, int> pii;
typedef vector<int> vi;
typedef vector<ll> vl;

int main() {
	cin.sync_with_stdio(0); cin.tie(0);
	cin.exceptions(cin.failbit);
	int n; cin >> n;
	vl x(n), y(n); rep(i,n) cin >> x[i] >> y[i];
	vl t(n); rep(i,n) t[i] = abs(x[i] - x[0])*abs(x[i] - x[0]) + abs(y[i] - y[0])*abs(y[i] - y[0]);
   	int e = 0, o = 0;
	rep(i,n){
		if(t[i]%2 == 1) o++;
		else e++;
	}
	while(o == 0){
		rep(i,n) t[i] /= 2;
		o = 0; e = 0;
		rep(i,n){
			if(t[i]%2 == 1) o++;
			else e++;
		}
	}
	cout << o << endl;
	rep(i,n) if(t[i]%2 == 1) cout << i+1 << " ";
	cout << endl;
			
	return 0;
}