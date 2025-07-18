#include <iostream>

#include <cstdio>

#include <cmath>

#include <vector>

#include <queue>

#include <stack>

#include <set>

#include <map>

#include <string>

#include <algorithm>

#include <fstream>

#include <utility>

#include <ctime>

#define ft first

#define sc second

#define mp make_pair

#define forn(i, n) for (int i = 0; i < n; (i++))

#define two(k) (k << 1)

using namespace std;

ifstream fin("INPUT.txt");

ofstream fout("OUTPUT.txt");

typedef int INT;

#define int short

typedef pair<int, int> pii;

typedef pair<long long, long long> pll;

typedef long long ll;

typedef unsigned long long ull;





const int N =500;

const long M = 10000000;

const long long INF = 1000000000000000000;

int mod = 1 << 15;



int n, m, k, f, s, suc = true;



vector<int> v[N];



int d[N][N], used[N][N];



int queuef[M];



int queues[M];



int pf[N][N], ps[N][N];



int d1[N], u1[N], d2[N], u2[N];



void bfs()

{

	for (int i = 0; i<N; i++)

		fill(d[i], d[i] + N, mod);

	long F = 0, S = 0;

	d[0][n-1] = 0;

	used[0][n-1] = 1;

	queuef[S] = 0;

	queues[S++] = n - 1;

	while (F != S)

	{

		int t1 = queuef[F], t2 = queues[F++];

		if (F > M)

			F = 0;

//		set<pair<pii, pii>> s;

		int cnt = 0;

		set<pii> s;

		for (int i = 0; i<v[t1].size(); i++)

			s.insert(mp(d2[v[t1][i]], v[t1][i]));

		for (auto i = s.begin(); i != s.end(); i++, cnt++)

			v[t1][cnt] = i->second;

		cnt = 0;

		s.clear();

		for (int i = 0; i<v[t2].size(); i++)

			s.insert(mp(d1[v[t2][i]], v[t2][i]));

		for (auto i = s.begin(); i != s.end(); i++, cnt++)

			v[t2][cnt] = i->second;

		int cnt1 = 0, cnt2 = 0;

		for (int i = 0;cnt1 < 4 && i<v[t1].size(); i++, cnt1++)

			for (int j = 0;cnt2<4 && j<v[t2].size(); j++, cnt2++)

			{

				int from = v[t1][i], to = v[t2][j];

				if (used[from][to] == 0 && from != to)

				{

					used[from][to] = 1;

					pf[from][to] = t1;

					ps[from][to] = t2;

					d[from][to] = d[t1][t2] + 1;

					queuef[S] = from;

					queues[S++] = to;

					if (from == n - 1 && to == 0)

						return;

					if (S > M)

						S = 0;

					if (S == F)

						return;

				}

			}

	}

}



void BFS(int k)

{

	for (int i = 0; i<n; i++)

		d1[i] = mod;

	d1[k] = 0;

	u1[k] = 1;

	queue<int> q;

	q.push(k);

	while (q.size())

	{

		int p = q.front();

		q.pop();

		for (int i = 0; i<v[p].size(); i++)

			if (u1[v[p][i]] == 0)

			{

				u1[v[p][i]] = 1;

				d1[v[p][i]] = 1 + d1[p];

				q.push(v[p][i]);

			}

	}

}



void BFS2(int k)

{

	for (int i = 0; i<n; i++)

		d2[i] = mod;

	d2[k] = 0;

	u2[k] = 1;

	queue<int> q;

	q.push(k);

	while (q.size())

	{

		int p = q.front();

		q.pop();

		for (int i = 0; i<v[p].size(); i++)

			if (u2[v[p][i]] == 0)

			{

				u2[v[p][i]] = 1;

				d2[v[p][i]] = 1 + d2[p];

				q.push(v[p][i]);

			}

	}

}



INT main()

{

//	freopen("INPUT.txt", "r", stdin);

//	freopen("OUTPUT.txt", "w", stdout);

	ios::sync_with_stdio(false);

	srand(time(0));

	cout.precision(10);

	cin >> n >> m;

	for (int i = 0; i < m; i++)

	{

		cin >> f >> s;

		v[f-1].push_back(s-1);

		v[s-1].push_back(f-1);

	}

	//for (int i = 0; i<n; i++)

	//	for (int j = 1; j<v[i].size(); j++)

	//		swap(v[i][rand()%(j+1)], v[i][j]);

	BFS(0);

	BFS2(n-1);

	bfs();

	if (d[n-1][0] == mod)

		cout << -1;

	else

	{

		cout << d[n-1][0] << endl;

		int ansf[10000], anss[10000];

		int fsize = 0, ssize = 0;

		int curf = n - 1, curs = 0;

		while (curf != 0 || curs != 0)

		{

			ansf[fsize++] = curf + 1;

			anss[ssize++] = curs + 1;

			int curt = curf;

			curf = pf[curf][curs];

			curs = ps[curt][curs];

		}

		for (int i = fsize - 1; i>= 0; i--)

			cout << ansf[i] << ' ';

		cout << endl;

		for (int i = ssize -1; i>=0; i--)

			cout << anss[i] << ' ';

	}

//	cin >> n;

	return 0;

}