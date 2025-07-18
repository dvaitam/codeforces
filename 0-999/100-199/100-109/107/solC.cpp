/*
 * 2011-08-23  Martin  <Martin@Martin-desktop>

 * 
 */
#include <algorithm>
#include <iostream>
#include <fstream>
#include <climits>
#include <vector>

using namespace std;

#define tr(i, x) for (typeof(x.begin()) i = x.begin(); i != x.end(); ++ i)
#define rep(i, n) for (int i = 0; i < n; ++ i)
#define pii pair <int, int>
#define mp make_pair
#define x first
#define y second
#define pb push_back
#define ll long long

#define MaxiN 20
#define MaxiSit 70000

int N, M, UsedNum;
ll Y;
vector <int> Adj[MaxiN];
bool Known[MaxiN];
int Num[MaxiN], Given[MaxiN];
int InE[MaxiN];
ll Opt[MaxiSit];

bool HasCircle(int Cur)
{
	if (Known[Cur])
		return 1;
	Known[Cur] = 1;
	tr (it, Adj[Cur])
		if (HasCircle(*it))
			return 1;
	Known[Cur] = 0;
	return 0;
}

bool Checked()
{
	fill(Known, Known + N, 0);
	rep (i, N)
		if (HasCircle(i))
			return 0;
	return 1;
}

ll Find()
{
	fill(Opt, Opt + (1 << N), 0);
	Opt[0] = 1;
	rep (i, (1 << N) - 1)
		if (Opt[i] > 0)
		{
			int c = __builtin_popcount(i);
			if (Given[c] == - 1)
			{
				rep (j, N)
					if (Num[j] == - 1 && (i & (1 << j)) == 0 && (i & InE[j]) == InE[j])
						Opt[i ^ (1 << j)] += Opt[i];
			}
			else
			{
				if ((i & InE[Given[c]]) == InE[Given[c]])
					Opt[i ^ (1 << Given[c])] += Opt[i];
			}
		}
	return Opt[(1 << N) - 1];
}

void GetAns()
{
	fill(Num, Num + N, - 1);
	fill(Given, Given + N, - 1);
	UsedNum = 0;
	rep (i, N)
	{
		bool Flag = 0;
		rep (j, N)
			if ((UsedNum & (1 << j)) == 0)
			{
				Num[i] = j;
				Given[j] = i;
				UsedNum ^= (1 << j);
				ll Cnt = Find();
				if (Cnt >= Y)
				{
					Flag = 1;
					break;
				}
				Y -= Cnt;
				Num[i] = - 1;
				Given[j] = - 1;
				UsedNum ^= (1 << j);
			}
		if (!Flag)
		{
			puts("The times have changed");
			return;
		}
	}
	rep (i, N)
		printf("%d%c", Num[i] + 1, (i == N - 1) ? '\n' : ' ');
}

int main()
{
	cin >> N >> Y >> M;
	Y -= 2000LL;
	fill(InE, InE + N, 0);
	rep (i, M)
	{
		int U, V;
		cin >> U >> V;
		-- U, -- V;
		Adj[U].pb(V);
		InE[V] |= (1 << U);
	}
	if (!Checked())
		puts("The times have changed");
	else
		GetAns();
}