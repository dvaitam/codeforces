#include<cstdio>
#include<cstring>
#include<algorithm>
#include<cctype>
#include<ctime>
#include<cstdlib>
#include<string>
#include<queue>
#include<cmath>
#include<set>
#include<map>
#include<bitset>
#include<vector>
#define For(x,a,b) for (int x=a;x<=(int)b;x++)
#define cross(x,a) for (int x=head[a];~x;x=nx[x])
#define ll long long
#define inf (1<<29)
#define PII pair<int,int>
#define PDD pair<double,double>
#define mk(a,b) make_pair(a,b)
#define fr first
#define sc second
#define pb push_back
using namespace std;
inline ll read(){
	ll x=0;int ch=getchar(),f=1;
	while (!isdigit(ch)&&(ch!='-')&&(ch!=EOF)) ch=getchar();
	if (ch=='-'){f=-1;ch=getchar();}
	while (isdigit(ch)){x=(x<<1)+(x<<3)+ch-'0';ch=getchar();}
	return x*f;
}
inline void rt(ll x){
	if (x<0) putchar('-'),x=-x;
	if (x>=10) rt(x/10),putchar(x%10+'0');
		else putchar(x+'0');
}
const int N=200005;
int a[N],n,color[N];
int head[N],to[N<<1],nx[N<<1],cnt;
void insert(int u,int v){
	to[cnt]=v;nx[cnt]=head[u];head[u]=cnt++;
}
void Dfs(int u,int fa){
	if (color[u]!=1) a[u]=1;
	cross(i,u){
		int v=to[i];
		if (v!=fa) Dfs(v,u),a[u]|=a[v];
	}
}
void Solve(int u,int fa){
	if (u!=1) color[u]^=1;
	rt(u),putchar(' ');
	cross(i,u){
		int v=to[i];
		if (v==fa) continue;
		if (a[v]||!color[v]){
			Solve(v,u);color[u]^=1;rt(u),putchar(' ');
			if (!color[v]){color[u]^=1;rt(v),putchar(' '),rt(u),putchar(' ');}
		}
	}
}
int main(){
	n=read();
	For(i,1,n) color[i]=read();
	For(i,1,n) if (color[i]==-1) color[i]=0;
	memset(head,-1,sizeof head);
	For(i,1,n-1){
		int u=read(),v=read();
		insert(u,v),insert(v,u);
	}
	Dfs(1,0);
	Solve(1,0);
	if (!color[1]) printf("%d %d %d",to[head[1]],1,to[head[1]]);
}