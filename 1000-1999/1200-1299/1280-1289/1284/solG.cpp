#include <bits/stdc++.h>

using namespace std;
const int N=10005;
const int M=105;
int head[N],nxt[N],edgenum=1,vet[N],num[N],val[N],dir[N],fa[N],flag[N],vis[N],n,m;
char ans[M][M],c[M][M];
const int dx[4]={1,-1,0,0};
const int dy[4]={0,0,1,-1};
void Add(int u,int v,int w) { vet[++edgenum]=v; nxt[edgenum]=head[u]; head[u]=edgenum; val[edgenum]=w; }
int find(int x) { if (x!=fa[x]) fa[x]=find(fa[x]); return fa[x]; }
bool dfs(int x) {
if (vis[x]) return 0;
vis[x]=1;
for (int e=head[x]; e; e=nxt[e]) {
int v=vet[e];
if (!vis[v] && (!num[v] || dfs(num[v]))) {
num[v]=x;
dir[v]=val[e]^1;
return 1;
}
}
return 0;
}
void dfs2(int x,int f) {
for (int e=head[x]; e; e=nxt[e]) {
int v=vet[e]; if (v==f) continue;
if (find(x)!=find(v)) {
fa[find(x)]=find(v);
ans[(x-1)/m*2+1+dx[val[e]]][(x-1)%m*2+1+dy[val[e]]]='O';
if (v!=1) dfs2(num[v],v);
}
}
}
int main() {
//	freopen("sample1.in","r",stdin);
//	freopen("1.out","w",stdout);
int deb=0;
int T; scanf("%d",&T);
if (deb) cout<<T<<endl;
while (T--) {
memset(fa,0,sizeof(fa)); memset(head,0,sizeof(head));
memset(dir,0,sizeof(dir));
edgenum=1; for (int i=0; i<M; i++) for (int j=0; j<M; j++) ans[i][j]=' ';
scanf("%d%d",&n,&m);
for (int i=1; i<=n*m; i++) num[i]=flag[i]=0,fa[i]=i;
for (int i=1; i<=n; i++) scanf("%s",c[i]+1);
if (deb) {
cout<<n<<' '<<m<<endl;
for (int i=1; i<=n; i++,puts("")) for (int j=1; j<=m; j++) printf("%c",c[i][j]);
}

for (int x=1; x<=n; x++)
for (int y=(x&1)+1; y<=m; y+=2)
if (c[x][y]=='O') {
for (int k=0; k<4; k++) {
int nx=x+dx[k],ny=y+dy[k];
if (nx<1 || ny<1 || nx>n || ny>m || (nx==1 && ny==1) || c[nx][ny]!='O') continue;
Add((x-1)*m+y,(nx-1)*m+ny,k);
}
}
for (int i=1; i<=n; i++)
for (int j=(i&1)+1; j<=m; j+=2)
if (c[i][j]=='O') {
for (int k=1; k<=n*m; k++) vis[k]=0;
flag[(i-1)*m+j]=dfs((i-1)*m+j);
}

for (int i=1; i<=n; i++)
for (int j=-(i&1)+2; j<=m; j+=2)
if (i+j>2 && c[i][j]=='O') {
if (!num[(i-1)*m+j]) { puts("NO"); goto out; }
fa[find((i-1)*m+j)]=find(num[(i-1)*m+j]);
ans[(i<<1)-1+dx[dir[(i-1)*m+j]]][(j<<1)-1+dy[dir[(i-1)*m+j]]]='O';
}

if (c[2][1]=='O') Add(m+1,1,1);
if (c[1][2]=='O') Add(2,1,3);
for (int i=1; i<=n; i++)
for (int j=(i&1)+1; j<=m; j+=2)
if (c[i][j]=='O' && !flag[(i-1)*m+j]) dfs2((i-1)*m+j,-1);

for (int i=1,rt=find(1); i<=n; i++)
for (int j=1; j<=m; j++) {
if (c[i][j]=='O' && find((i-1)*m+j)!=rt) { puts("NO"); goto out; }
ans[(i<<1)-1][(j<<1)-1]=c[i][j];
}
if (deb) { for (int i=1; i<(n<<1); i++) for (int j=1; j<(m<<1); j++) if (ans[i][j]==' ') ans[i][j]='K'; }
puts("YES");
for (int i=1; i<(n<<1); i++,puts(""))
for (int j=1; j<(m<<1); j++) printf("%c",ans[i][j]);
out:;
}
return 0;
}