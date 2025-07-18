//求边双连通分量（e-DCC）
#include <cstdio>
#include <iostream>
#include <algorithm>
#include <cstring>
#include <queue>
#define rep(i,a,b) for(int i = a;i <= b;i++)
using namespace std;
const int N = 2*1e5+100;

struct Edge{	
	int to,next;
	int from;
	int vis;
}e[2*N]; //双向边
int head[N],deg[N];  //dfn:记录该点的dfs序		low:记录该点所能达到的最小dfn值
int tot,n,m,dcc,D; //tot记录边序，num记录dfs序，dcc记录有多少个边双连通分量
int c[N],ph[N],tnum,book[N],vvis[N];  //c[x]表示节点x所属的“边双连通分量”的编号

void add(int x,int y)
{
	e[++tot].to = y, e[tot].next = head[x], head[x] = tot, e[tot].vis = 0, e[tot].from = x;
}

void tarjan(int x)
{
	c[x] = dcc;
	for(int i = head[x]; i ; i = e[i].next)
	{
		int y = e[i].to;
		if(y == 1 || vvis[y]) continue;
		vvis[y] = 1;
		tarjan(y);
	}
}

// void dfs(int x)
// {
// 	c[x] = dcc;
// 	for(int i = head[x]; i ; i = e[i].next)
// 	{
// 		int y = e[i].to;
// 		if(c[y] || bridge[i]) continue;  //如果点y已经属于别的强连通分量，或者边i是割边，则continue
// 		dfs(y);
// 	}
// }

void solve()
{
	queue<int> q;
	while(q.size()) q.pop();
	book[1] = 1;
	tnum = 0;
	for(int i = head[1]; i; i = e[i].next)
	{
		int y = e[i].to;
		ph[c[y]]++;
		if(ph[c[y]] == 1) tnum++;
	}
	if(tnum > D){
		printf("NO\n");
		return;
	}
	memset(ph,0,sizeof ph);
	int c1 = tnum, c2 = D-tnum;
	for(int i = head[1]; i; i = e[i].next)
	{
		int y = e[i].to;
		ph[c[y]]++;
		// printf("y:%d,c[y]:%d\n",y,c[y]);
		if(ph[c[y]] == 1){
			c1--,book[y] = 1,e[i].vis = 1;
			q.push(y);
		} 
		else if(c2 > 0){
			c2--,book[y] = 1,e[i].vis = 1;
			q.push(y);	
		}
	}
	while(q.size())
	{
		int a = q.front();
		q.pop();
		// if(book[a]) continue;
		// book[a] = 1;
		for(int i = head[a]; i; i = e[i].next){
			int b = e[i].to;
			if(book[b]) continue;
			book[b] = 1;
			e[i].vis = 1;
			q.push(b);
		}
	}
	printf("YES\n");
	rep(i,1,tot){
		if(e[i].vis == 1){
			printf("%d %d\n",e[i].from,e[i].to);
		}
	}
}

int main()
{
	scanf("%d%d%d",&n,&m,&D);
	tot = 1; dcc = 0;
	rep(i,1,m){
		int x,y;
		scanf("%d%d",&x,&y);
		add(x,y); add(y,x);
		deg[x]++, deg[y]++;
	}
	if(deg[1] < D) printf("NO\n");
	else{
		rep(i,2,n)
			if(!vvis[i]){
				dcc++;
				tarjan(i); //将边序号传进去
			} 
		solve();
	}
	/*	for(int i = 2; i <= tot; i+=2)
			if(bridge[i])
				printf("%d %d\n",e[i^1].to,e[i].to);  //该割边：(x,y)*/
	return 0;
}

/*
11 13 2
1 6
1 5
1 4
5 11
6 11
4 10
11 10
1 2
1 3
2 7
3 8
7 8
8 9
*/