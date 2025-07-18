#include <bits/stdc++.h>
using namespace std;

const int maxn = 1005;

int n, seq[maxn];
int L[maxn], R[maxn], cnt;

struct node
{
	int	w, id;
	node(int _w, int _id) : w(_w), id(_id) { }
	inline bool operator < (const node &T) const
	{
		return w < T.w;
	}
};

int main()
{
	scanf("%d", &n);
	for (int i = 1; i <= n; i++)
	{
		scanf("%d", &L[i]);
		if(L[i] > i - 1)
			return puts("NO"), 0;
	}
	for (int i = 1; i <= n; i++)
	{
		scanf("%d", &R[i]);
		if(R[i] > n - i)
			return puts("NO"), 0;
	}
	cnt = n;
	queue<node>	Q;
	for(int i = 1; i <= n; i++)
		if(L[i] == 0 && R[i] == 0)
		{
			seq[i] = n;
			Q.emplace(n, i);
		}
	cnt = n;
	while(!Q.empty())
	{
		while(!Q.empty())
		{
			node	x = Q.front();
			Q.pop();
			for(int i = x.id + 1; i <= n; i++)
			{
				if(seq[i])
					continue;
				--L[i];
			}
			for(int i = 1; i < x.id; i++)
			{
				if(seq[i])
					continue;
				--R[i];
			}
		}
		cnt--;
		for(int i = 1; i <= n; i++)
			if(!seq[i] && L[i] == 0 && R[i] == 0)
				seq[i] = cnt, Q.emplace(cnt, i);
	}
	for(int i = 1; i <= n; i++)
		if(!seq[i])
			return puts("NO"), 0;
	puts("YES");
	for(int i = 1; i <= n; i++)
		printf("%d%c", seq[i], " \n" [i == n]);
	return 0;
}