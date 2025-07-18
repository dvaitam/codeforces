#include<bits/stdc++.h>

using namespace std;

#define forn(i,n) for (int i = 0; i < int(n); i++)
#define ford(i,n) for (int i = int(n) - 1; i >= 0; i--)
#define fore(i,l,r) for (int i = int(l); i <= int(r); i++)
#define all(a) a.begin(), a.end()
#define sz(a) int(a.size())
#define mp make_pair
#define pb push_back
#define ft first
#define sc second
#define x first
#define y second

template<typename X> inline X abs(const X& a) { return a < 0 ? -a : a; }
template<typename X> inline X sqr(const X& a) { return a * a; }

typedef long long li;
typedef long double ld;
typedef pair<int, int> pt;

const int INF = int(1e9);
const li INF64 = li(1e18);
const ld EPS = 1e-9;
const ld PI = acosl(ld(-1));

mt19937 mt(time(NULL));

const int N = 50 * 1000 + 13;
int n, m;
pt g[2 * N];
int l[N];
int a[N][21];
int maxT;

inline void gen()
{
	n = 50 * 1000, m = 100 * 1000;
	
	set<pt> q;
	while (sz(q) < m)
	{
		int x = mt() % n, y = mt() % n;
		
		if (x == y)
			continue;
			
		if (x > y)
			swap(x, y);
			
		q.insert(mp(x, y));
	}
	
	m = 0;
	for(auto it: q)
		g[m++] = it;
		
	int T = 1000;
	maxT = 0;
		
	forn (i, n)
	{
		l[i] = 16;
		
		set<int> q;
		while (sz(q) < l[i])
		{
			int x = mt() % T;
			q.insert(x);
		}
		
		l[i] = 0;
		for(auto it: q)
		{
			a[i][ l[i]++ ] = it;
			maxT = max(maxT, it);
		}
	}
	
	maxT++;
}

inline bool read()
{
	//gen();
	//return true;
	
	if (scanf ("%d%d", &n, &m) != 2)
		return false;
		
	forn (i, m)
	{
		int x, y;
		assert(scanf ("%d%d", &x, &y) == 2);
		
		x--, y--;
		
		g[i] = mp(x, y);
	}
	
	maxT = 0;
	
	forn (i, n)
	{
		assert(scanf ("%d", &l[i]) == 1);
		
		forn (j, l[i])
		{
			assert(scanf ("%d", &a[i][j]) == 1);
			
			a[i][j]--;
			
			maxT = max(maxT, a[i][j]);
		}
	}
	
	maxT++;
	
	return true;
}

int t[1000 * 1000 + 13];
int p[N];
int fr[21], szfr;
int col[N];

inline void solve()
{
	//cerr << clock() << endl;

	forn (i, n)
		p[i] = i;

	while (true)
	{
		fore (i, 1, n - 1)
		{
			int x = mt() % i;
			swap(p[i], p[x]);
		}
		
		memset(t, -1, sizeof t);
		
		forn (_i, n)
		{
			int i = p[_i];
			
			szfr = 0;
			int pos0 = -1, pos1 = -1;
			forn (j, l[i])
			{
				int x = a[i][j];
				
				if (t[x] == -1)
					fr[szfr++] = x;
				else if (t[x] == 0)
					pos0 = x;
				else if (t[x] == 1)
					pos1 = x;
				else
					throw;
			}
			
			if (szfr && pos0 == -1)
			{
				int x = mt() % szfr;
				
				t[ fr[x] ] = 0;
				
				pos0 = fr[x];

				swap(fr[x], fr[szfr - 1]);
				szfr--;
			}
			
			if (szfr && pos1 == -1)
			{
				int x = mt() % szfr;
				
				t[ fr[x] ] = 1;
				
				pos1 = fr[x];
			}
			
			if (pos0 != -1 && pos1 != -1)
			{
				if (mt() & 1)
					col[i] = pos0;
				else
					col[i] = pos1;
					
				continue;
			}
			
			if (pos0 != -1)
			{
				col[i] = pos0;
				continue;
			}
			
			if (pos1 != -1)
			{
				col[i] = pos1;
				continue;
			}
			
			throw;
		}
		
		int cnt = 0;
		forn (i, m)
		{
			int x = t[ col[ g[i].x ] ], y = t[ col[ g[i].y ] ];
			
			//cerr << x << ' ' << y << endl;
			
			if (x != y)
				cnt++;
		}
		
		//cerr << cnt << endl;
		
		if (2 * cnt >= m)
		{
			forn (i, n)
				printf ("%d ", col[i] + 1);
			puts("");
			forn (i, maxT)
				if (t[i] == -1)
					printf ("1 ");
				else
					printf ("%d ", t[i] + 1);
			puts("");

			return;
		}
	}
}

int main()
{
#ifdef SU2_PROJ
	assert(freopen("input.txt", "r", stdin));
	assert(freopen("output.txt", "w", stdout));
#endif

	cout << setprecision(25) << fixed;
	cerr << setprecision(10) << fixed;

	srand(int(time(NULL)));

	assert(read());
	solve();

#ifdef SU2_PROJ
	cerr << "TIME: " << clock() << endl;
#endif

	return 0;
}