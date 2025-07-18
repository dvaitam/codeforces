#pragma comment (linker, "/STACK:256000000")
#define _CRT_SECURE_NO_WARNINGS
#include <cstdio>
#include <vector>
#include <algorithm>
#include <iostream>
#include <string>

using namespace std;

int n, m;
vector <vector <int> > g;
vector <int> color;
vector <int> ts;

void dfs(int v)
{
	color[v] = 1;
	for (auto u = g[v].begin(); u != g[v].end(); ++u)
	{
		if (!color[*u])
			dfs(*u);
	}
	printf("%d ", v + 1);
}

int main()
{
	scanf("%d %d", &n, &m);
	g.resize(n);
	color.assign(n, 0);
	for (int i = 0; i < m; ++i)
	{
		int a, b;
		scanf("%d %d", &a, &b);
		--a, --b;
		g[a].push_back(b);
	}
	for (int v = 0; v < n; ++v)
		if (!color[v])
			dfs(v);
	printf("\n");
	return 0;
}