#include<cstdio>
#include<algorithm>
#include<cmath>
#include<vector>
using namespace std;
const int MAX = 100010;
int a[MAX];
double p[MAX][110];
int main()
{
    int n,q;
    double ans = 0;
    scanf("%d",&n);
    for (int i=1;i<=n;++i) {
        scanf("%d",&a[i]);
        p[i][a[i]]=1;
        if (a[i]==0)
            ans += 1;
    }
    scanf("%d",&q);
    for (int i=0;i<q;++i) {
        int u,v,k;
        scanf("%d%d%d",&u,&v,&k);
        for (int j=0;j<k;++j) {
            ans -= p[u][0];
            for (int z=0;z<=min(100,a[u]);++z)
                p[u][z] = (p[u][z]*(a[u]-z)+p[u][z+1]*(z+1))/a[u];
            ans += p[u][0];
            a[u]--;
        }
        a[v]+=k;
        printf("%.10lf\n",ans);
    }
    return 0;
}