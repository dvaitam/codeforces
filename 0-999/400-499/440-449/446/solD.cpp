#include <algorithm>
#include <iostream>
#include <cstdlib>
#include <cstring>
#include <cstdio>
#include <cmath>
#define LL  long long
#define oo  (1<<30)
#include <queue>
using namespace std;

struct matrix{double n[105][105];}p;
struct edge{int s,t,n;}e[200010];
int n,m,k,W,B,a[505],c[505],b[105],h[505],s[100010],t[100010];
int cnt[105][105],Cnt[105][105],s1[105],s2[105];
double ans;

matrix operator * (matrix a,matrix b)
{
    matrix c;
    memset(c.n,0,sizeof(c.n));
    for (int i=1; i<=B; i++)
        for (int j=1; j<=B; j++)
            for (int k=1; k<=B; k++)
                c.n[i][j]+=a.n[i][k]*b.n[k][j];
    return c;
}

void BFS(int s)
{
    int x,y,i;
    queue <int> Q;
    Q.push(s),c[s]=W;
    while (!Q.empty())
        for (i=h[x=Q.front()],Q.pop(); y=e[i].t,i; i=e[i].n)
            if ((!a[y])&&(!c[y]))  Q.push(y),c[y]=W;
}

matrix pow(matrix a,int b)
{
    matrix x=a;
    for (int i=1; i<=B; i++)
        for (int j=1; j<=B; j++)
            a.n[i][j]=(i==j);
    while (b)
        {
            if (b&1)  a=a*x;
            x=x*x,b>>=1;
        }
    return a;
}

void work()
{
    scanf("%d %d %d",&n,&m,&k);
    for (int i=1; i<=n; i++)
        scanf("%d",&a[i]);
    for (int i=1,tot=0; i<=m; i++)
        {
            scanf("%d %d",&s[i],&t[i]);
            e[++tot]=(edge){s[i],t[i],h[s[i]]},h[s[i]]=tot;
            e[++tot]=(edge){t[i],s[i],h[t[i]]},h[t[i]]=tot;
        }
    for (int i=1; i<=n; i++)  if ((!a[i])&&(!c[i]))  W++,BFS(i);
    for (int i=1; i<=n; i++)  if (a[i])  c[i]=++B,b[B]=i;
    for (int i=1; i<=m; i++)
        if ((a[s[i]])&&(a[t[i]]))
            Cnt[c[s[i]]][c[t[i]]]++,Cnt[c[t[i]]][c[s[i]]]++;
        else  if ((a[s[i]])&&(!a[t[i]]))  cnt[c[s[i]]][c[t[i]]]++;
        else  if ((!a[s[i]])&&(a[t[i]]))  cnt[c[t[i]]][c[s[i]]]++;
    for (int i=1; s1[i]=0,i<=W; i++)
        for (int j=1; j<=B; j++)
            s1[i]+=cnt[j][i];
    for (int i=1; s2[i]=0,i<=B; i++)
        {
            for (int j=1; j<=B; j++)  s2[i]+=Cnt[j][i];
            for (int j=1; j<=W; j++)  s2[i]+=cnt[i][j];
        }
    for (int i=1; i<=B; i++)
        for (int j=1; j<=B; j++)
            {
                p.n[i][j]=1.0*Cnt[i][j]/s2[i];
                for (int k=1; k<=W; k++)
                    p.n[i][j]+=(1.0*cnt[i][k]/s2[i])*(1.0*cnt[j][k]/s1[k]);
            }
    p=pow(p,k-2);
    for (int i=1,j=c[1]; i<=B; i++)
        ans+=1.0*cnt[i][j]/s1[j]*p.n[j][B];
    printf("%lf",ans);
}

int main()
{
    work();
    return 0;
}