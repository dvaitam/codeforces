//              +-- -- --++-- +-In the name of ALLAH-+ --++-- -- --+              \\

#include <bits/stdc++.h>

#define int ll

#define F first
#define S second
#define _sz(x) (int)x.size()

using namespace std ;
using ll = long long ;
using ld = long double ;
using pii = pair <int , int> ;

int const N = 100 + 20 , M = 5000 + 220 , mod = 998244353 ;
int n , r , s ;
int ans ;
int C[M][N] ;

inline void rel (int &a) {
	if (a >= mod) a -= mod ;
}

inline int pw (int a , int b) {
	int res = 1 ;

	while (b) {
		if (b & 1) res = 1ll * res * a % mod ;
		a = 1ll * a * a % mod ;
		b >>= 1 ;
	}

	return res ;
}

int32_t main(){
	ios::sync_with_stdio(false) , cin.tie(0) , cout.tie(0) ;

	for (int i = 0 ; i < M ; i ++) {
		C[i][0] = 1 ;
		for (int j = 1 ; j < N && j <= i ; j ++) {
			C[i][j] = C[i - 1][j] + C[i - 1][j - 1] ;
			rel(C[i][j]) ;
		}
	}

	cin >> n >> s >> r ;

	if (r == 0) return cout << pw(n , mod - 2) << '\n' , 0 ;

	for (int p = 1 ; p <= n ; p ++) {
		int save = pw(p , mod - 2) ;

		for (int k = r ; p * k <= s ; k ++) {
			int SUM = s - p * k ;
			int t = n - p ;
			int res = 0 ;

			if (!t) {
				if (!SUM) ans = ans + save , rel(ans) ;
				continue ;
			}

			for (int bad = 0 ; bad <= t && bad * k <= SUM ; bad ++) {			
				int nsum = SUM - bad * k ;
		
				int add = 1ll * C[nsum + t - 1][t - 1] * C[t][bad] % mod ;

				if (bad & 1) res = (res - add + mod) , rel(res) ;
				else res = (res + add) , rel(res) ;
			}

			res = 1ll * res * C[n - 1][p - 1] % mod ;

			ans += 1ll * res * save % mod ;
			rel(ans) ;
		}
	}

	ans = 1ll * ans * pw(C[s - r + n - 1][n - 1] , mod - 2) % mod ;

	cout << ans << '\n' ;
}