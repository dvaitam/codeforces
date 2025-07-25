#include<iostream>
#include<cstdio>
#include<vector>
#define lld long long
using namespace std;
const int MaxN = 100000 + 10;
int n, k, d[MaxN];
vector<int>G[MaxN];
bool ok1()
{
	lld mxn = 1;
	for(int i = 0; i < n; ++i)
		if(G[i].size() > mxn) return false;
		else mxn = (lld)G[i].size() * (lld)(i == 0 ? k : (k - 1));
	return true;
}
void out()
{
	for(int i = 1; i < n; ++i)
	{
		int p = 0, r = (i == 1 ? k : k - 1);
		for(int j = 0; j < G[i].size(); ++j)
		{
			if(!r) {++p;r = (i == 1 ? k : k - 1);}
			printf("%d %d\n", G[i - 1][p], G[i][j]);
			--r;
		}
	}
}
int main()
{
	scanf("%d%d", &n, &k);
	for(int i = 1; i <= n; ++i) scanf("%d", &d[i]);
	for(int i = 1; i <= n; ++i)	G[d[i]].push_back(i);
	if(!ok1()){printf("-1\n");return 0;}
	printf("%d\n", n - 1);
	out();
	return 0;
}