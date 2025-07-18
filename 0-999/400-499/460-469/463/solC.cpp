#include <bits/stdc++.h>

using namespace std;

/* Template */

#define ll long long
#define frt(it, S) for(__typeof (S.begin()) it = S.begin(); it != S.end(); it++)
#define all(x) x.begin(), x.end()
#define rall(x) x.rbegin(), x.rend()
#define MN(a, b) a = min (a, b)
#define MX(a, b) a = max (a, b)
#define endl '\n'

#define fname ""

#define nxtnt getNext <int>()
#define nxtll getNext <ll>()

template <typename T> inline T sqr (T x) { return x * x; }

template <typename T>
inline T getNext () {
	T s = 1, x = 0, c = getc(stdin);
	while (c <= 32)
		c = getc(stdin);
	if (c == '-')
		s = -1, c = getc(stdin);
	while ('0' <= c && c <= '9')
		x = x * 10 + c - '0', c = getc(stdin);
	return x * s;
}

const int MaxN = int (3e3);

ll a[MaxN][MaxN], d1[MaxN][MaxN], d2[MaxN][MaxN];

int main() {
	#ifndef ONLINE_JUDGE
		freopen (".in", "r", stdin);
		// freopen("out", "w", stdout);
	#endif

	int n = nxtnt;
	for(int i = 1; i <= n ; ++i)
		for(int j = 1; j <= n; ++j)
			a[i][j] = nxtnt;

	for(int i = 1; i <= n; ++i)
		for(int j = 1; j <= n; ++j) {
			d2[i][j] = d2[i - 1][j - 1] + a[i][j];
			d1[i][j] = d1[i - 1][j + 1] + a[i][j];	
		}
	
	ll b = 0, w = 0;
	int x = 1, y=1, x2 = 2, y2 = 1;
	for(int i = 1; i <= n; ++i)
		for(int j = 1; j <= n; ++j)  {
			int k = min (n - i, n - j);
			ll tmp = d2[i + k][j + k];
			k = min (n - i, j - 1);
			tmp += d1[i + k][j - k];
			tmp -= a[i][j];
			if((i + j) % 2 == 0){
				if (b < tmp) {
					b = tmp;
					x = i;
					y = j;
		 	 	}
			} else {
				if(w < tmp) {
					w = tmp;
		 	  	 	x2 = i;
		 	  	 	y2 = j;
		 	  	}
			}
		}


	cout << b + w << endl;
	cout << x << " " << y << " " << x2 << " " << y2;
	return 0; 
}