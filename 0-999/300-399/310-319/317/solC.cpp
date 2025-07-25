#include<cstdio>
#define rep(i,n) for (int i=1;i<=n;i++)
#define prt(x,y,z) if (z) S[++L]=x,T[L]=y,F[L]=z,a[x]-=z,a[y]+=z
#define N 605
int a[N],b[N],f[N],ed[N],next[N],son[N],S[N*N],T[N*N],F[N*N],n,m,v,l,x,y,L;
long long A,B,s[N];
int get(int x){return !f[x]?x:f[x]=get(f[x]);}
void pre(int x,int fa)
{s[x]=a[x]-b[x];for(int p=son[x],y;p;p=next[p])if((y=ed[p])!=fa)pre(y,x),s[x]+=s[y];}
void dfs(int x,int fa,int z)
{
	for (int p=son[x],y,t,w;p && z;p=next[p]) if ((y=ed[p])!=fa) if (s[y]>0 && (w=s[y]<z?s[y]:z)>0)
		if (z-=w,(t=a[y])>=w){prt(y,x,w); if (a[y]<b[y]) dfs(y,x,t<b[y]?w:b[y]-a[y]);}
		else {prt(y,x,t); dfs(y,x,w+(t<b[y]?0:b[y]-t)); prt(y,x,w-t);}
}
int main()
{
	scanf("%d%d%d",&n,&v,&m); rep(i,n) scanf("%d",a+i); rep(i,n) scanf("%d",b+i);
	rep(i,m){
		scanf("%d%d",&x,&y);
		if (get(x)!=get(y)) f[get(x)]=get(y),
			ed[++l]=y,next[l]=son[x],son[x]=l,
			ed[++l]=x,next[l]=son[y],son[y]=l;
		}
	rep(i,n) s[get(i)]+=a[i]-b[i]; rep(i,n) if (s[i]) return puts("NO"),0;
	rep(i,n) if (a[i]<b[i]) pre(i,0),dfs(i,0,b[i]-a[i]); printf("%d\n",L);
	rep(i,L) printf("%d %d %d\n",S[i],T[i],F[i]); return 0;
}