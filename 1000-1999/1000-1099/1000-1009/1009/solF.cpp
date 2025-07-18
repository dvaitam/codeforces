#include<cmath>
#include<cstdio>
#include<cstdlib>
#include<cstring>
#include<algorithm>
#define N 1000005 
using namespace std;
namespace LonelyDonkey
{
	int in()
	{
		char c = getchar(); int r = 0;
		for(; c < '0' || c > '9'; c = getchar());
		for(; c >='0' && c <='9'; r = r * 10 + c - '0', c = getchar());
		return r;
	}
	
	int n, last[N], ecnt;
	
	struct edge{int next, to;}e[N<<1];
	
	void addedge(int a, int b)
	{
		e[++ecnt] = (edge){last[a], b};
		last[a] = ecnt;
	}
	
	int mem[N*23], pos[N], len[N], mx[N], mxid[N], ans[N];
	
	void merge(int x, int y) // x <- y
	{
		if(len[x] < len[y])
		{
			swap(pos[x], pos[y]);
			swap(len[x], len[y]);
			swap(mx[x], mx[y]);
			swap(mxid[x], mxid[y]);
		}
		for(int i = 0; i < len[y]; i++)
		{
			mem[pos[x]+i] += mem[pos[y]+i];
			if(mem[pos[x]+i] > mx[x]) mx[x] = mem[pos[x]+i], mxid[x] = i;
			else if(mem[pos[x]+i] == mx[x] && i < mxid[x]) mxid[x] = i;
		}
	}
	
	void dfs(int x, int fa)
	{
		mx[x] = 1;
		mxid[x] = 0;
		
		len[x] = 1;
		mem[pos[x]] = 1;

//printf("pos %d = %d\n",x,pos[x]);
		
		for(int i = last[x]; i; i = e[i].next)
		{
			int y = e[i].to; if(y == fa) continue;
			
			pos[y] = pos[x] + len[x] + 1;
			
			dfs(y, x);
			
			--pos[y];
			++len[y];
			mem[pos[y]] = 0;
			mxid[y]++;
			
			merge(x, y);
		}
//printf("pos %d = %d\n",x,pos[x]);
		
//printf("x = %d mx = %d mxid = %d\n",x,mx[x],mxid[x]);
//for(int i = 0; i < len[x]; i++) printf("(%d)%d ",pos[x]+i,mem[pos[x]+i]); puts("");
		
		ans[x] = mxid[x];
	}
	
	void main()
	{
//		freopen("data.in","r",stdin);
		
		n = in();
		for(int i = 1; i < n; i++)
		{
			int a = in(), b = in();
			addedge(a, b);
			addedge(b, a);
		}
		pos[1] = 1;
		dfs(1,0);
		for(int i = 1; i <= n; i++) printf("%d\n",ans[i]);
	}
}
int main()
{
	LonelyDonkey::main();
}