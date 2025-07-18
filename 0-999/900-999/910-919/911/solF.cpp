#include <bits/stdc++.h>
#pragma GCC optimize("Ofast")
#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,avx,mmx,tune=native")
#define ll long long
using namespace std;
const int MAXBUF=1<<23;
char B[MAXBUF],*Si = B,*Ti = B;
inline char getc()
{if (Si==Ti) Ti=(Si=B)+fread(B,1,MAXBUF,stdin);
if (Si==Ti) return 0;
else return *Si++;
}
template <class T>
inline void read(T &a)
{static char c;
static int fh;
while (((c=getc())<'0'||c>'9')&&c!='-');
if (c=='-') fh=-1,a=0;
else fh=1,a=c-'0';
while ((c=getc())<='9'&&c>='0') a=(a<<3)+(a<<1)+c-'0';
if (fh==-1) a=-a;
}
char Buff[MAXBUF], *sti=Buff;
template <class T>
inline void write(T a)
{if (a==0) {*sti++='0';return;}
if (a<0) *sti++='-',a=-a;
static char c[20];
static int c0;c0=0;
while (a) c[c0++]=a%10+'0',a/=10;
while (c0--) *sti++=c[c0];
}
int n,fa[200005],head[200005],nxt[400005],v[400005],tot=0;
int dis[200005],maxn[200005],depa[200005],depb[200005];
int d[200005],x[200005],y[200005],del[200005],cnt=0;
int c[200005],tp=0;
inline void add(int a,int b)
{tot++;nxt[tot]=head[a];head[a]=tot;v[tot]=b;}
int getd(int pos,int f)
{int ret=pos;
maxn[pos]=dis[pos];
for (int i=head[pos];i;i=nxt[i])
{if (v[i]==f) continue;
dis[v[i]]=dis[pos]+1;
int tp=getd(v[i],pos);
if (maxn[v[i]]>maxn[pos])
{maxn[pos]=maxn[v[i]];ret=tp;}
}
return ret;
}
void dfs1(int pos,int f)
{for (int i=head[pos];i;i=nxt[i])
{if (v[i]==f) continue;
depa[v[i]]=depa[pos]+1;
fa[v[i]]=pos;d[pos]++;
dfs1(v[i],pos);
}
}
void dfs2(int pos,int f)
{for (int i=head[pos];i;i=nxt[i])
{if (v[i]==f) continue;
depb[v[i]]=depb[pos]+1;
dfs2(v[i],pos);
}
}
int main (){
 int i,a,b,u,v;
 read(n);
 for (i=1;i<n;i++)
 {read(a);read(b);
 add(a,b);add(b,a);
 }
 ll ans=0;
 dis[1]=0;u=getd(1,0);
 dis[u]=0;v=getd(u,0);
 dfs1(u,0);dfs2(v,0);
 for (i=1;i<=n;i++)
 {if (i==u||i==v) continue;
 if (!d[i]) c[++tp]=i;
 }
 while (tp) 
 {int p=c[tp--];d[fa[p]]--;
 if (!d[fa[p]]) c[++tp]=fa[p];
 int d1=depa[p],d2=depb[p];
 x[++cnt]=p;del[cnt]=p;
 if (d1>d2) {y[cnt]=u;ans+=d1;}
 else {y[cnt]=v;ans+=d2;}
 }
 while (u!=v)
 {x[++cnt]=u;y[cnt]=v;del[cnt]=v;
 ans+=depa[v];v=fa[v];
 }
 write(ans),*sti++='\n';
 for (i=1;i<n;i++)
 {write(x[i]),*sti++=' ';
 write(y[i]),*sti++=' ';
 write(del[i]),*sti++='\n';
 }
 fwrite(Buff,1,sti-Buff,stdout);
 return 0;
}