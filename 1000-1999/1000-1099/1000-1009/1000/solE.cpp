#include<bits/stdc++.h>
inline int read(){
	char c=getchar();while (c!='-'&&(c<'0'||c>'9'))c=getchar();
	int k=1,kk=0;if (c=='-')c=getchar(),k=-1;
	while (c>='0'&&c<='9')kk=kk*10+c-'0',c=getchar();return kk*k;
}using namespace std;
void write(int x){if (x<0)putchar('-'),x=-x;if (x/10)write(x/10);putchar(x%10+'0');}
void writeln(int x){write(x);puts("");}
const int maxn=300010,maxe=600010;
int times,in[maxn],low[maxn];
bool isb[maxe];
int n,e,x[maxe],y[maxe];
int tot,son[maxe],lnk[maxn],nxt[maxe];
void add(int x,int y){son[++tot]=y;nxt[tot]=lnk[x];lnk[x]=tot;}
void tarjan(int x,int fa){
    low[x]=in[x]=++times;
    for (int j=lnk[x];j;j=nxt[j])
     if (!in[son[j]]){
        tarjan(son[j],x);
        low[x]=min(low[x],low[son[j]]);
        if (low[son[j]]>in[x]) isb[j]=isb[j^1]=1;
     }else if (son[j]!=fa) low[x]=min(low[x],in[son[j]]);else fa=0;
}
int BCC,bcc[maxn];
void bl(int x){bcc[x]=BCC;for (int j=lnk[x];j;j=nxt[j])if (!isb[j]&&!bcc[son[j]]) bl(son[j]);}
int a[maxe][2],last[maxn],kk,ma[maxn],ans;
void doit(int x,int y){a[++kk][0]=last[x];a[kk][1]=y;last[x]=kk;} 
void dfs(int x,int fa){
	for (int i=last[x];i;i=a[i][0])if (a[i][1]!=fa){
		dfs(a[i][1],x);ans=max(ans,ma[a[i][1]]+ma[x]+1);
		ma[x]=max(ma[a[i][1]]+1,ma[x]);
	}
}
signed main(){
    n=read(),e=read();tot=1;
    for (int i=1;i<=e;i++) x[i]=read(),y[i]=read(),add(x[i],y[i]),add(y[i],x[i]);
    for (int i=1;i<=n;i++)if (!in[i]) tarjan(i,0);
    for (int i=1;i<=n;i++)if (!bcc[i]) BCC++,bl(i);
    for (int i=1;i<=e;i++)if (bcc[x[i]]!=bcc[y[i]])doit(bcc[x[i]],bcc[y[i]]),doit(bcc[y[i]],bcc[x[i]]);
    dfs(1,0);writeln(ans);
}