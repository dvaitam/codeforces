#include <iostream>

#include <stdio.h>

#include <vector>

#include <algorithm>

#include <set>

#include <map>

#include <cmath>

#include <string>

#include <cstring>

#include <ctime>

#include <cassert>

#include <queue>

#include <stack>

#include <bitset>

#define rnd() ((rand() << 15) + rand())

#define y1 y11

#define fs first

#define sc second

#define mp make_pair

#define pb push_back

#define mt make_tuple

#define NAME ""



using namespace std;

	

typedef long long ll;

typedef long double ld;



const ld PI = acos(-1.0);



const int MAXN = 3001;

const int INF = 1000 * 1000 * 100 * 3 + 2;

int x[MAXN], y[MAXN], d[MAXN][MAXN];

unsigned int msk[MAXN][MAXN];



pair <int, pair <int, int> > a[MAXN * MAXN];

int ac = 0;

int main()

{

	//freopen("input.txt", "r", stdin);

	//freopen("output.txt", "w", stdout);

	int n;

	cin >> n;

	for (int i = 0; i < n; i++)

	{

		cin >> x[i] >> y[i];

		for (int j = 0; j < i; j++)

		{

			a[ac++] = mp((x[i] - x[j]) * (x[i] - x[j]) + (y[i] - y[j]) * (y[i] - y[j]), mp(i, j));

		}

	}

	int g = (n - 1) / 32 + 1;

	memset(msk, 0, sizeof(msk));

	sort(a, a + ac);

	reverse(a, a + ac);

	cout.setf(ios::fixed);

	cout.precision(20);

	for (int i = 0; i < ac; i++)

	{

		int x = a[i].sc.fs, y = a[i].sc.sc;

		for (int j = 0; j < g; j++)

		{

			if (msk[x][j] & msk[y][j])

			{

				cout << sqrtl((ld)(a[i].fs)) /  2.0 << endl;

				return 0;

			}

		}

		msk[x][y >> 5]  += ((unsigned int)1 << (y & 31));

		msk[y][x >> 5]  += ((unsigned int)1 << (x & 31));

	}

	return 0;

}