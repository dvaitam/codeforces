#include <stdio.h>
#include <vector>
#include <algorithm>
using namespace std;
vector<int> in[100010],vw[100010];
int su[100010],zz[100010],jg[1000010];
bool bk[1000010];
int A[1000010],B[1000010];
int main()
{
    int n,m;
    scanf("%d%d",&n,&m);
    for(int i=0;i<m;i++)
    {
        int a,b;
        scanf("%d%d",&a,&b);
        A[i]=a;B[i]=b;
        if(a>b)swap(a,b);
        in[b].push_back(a);
        vw[b].push_back(i);
        su[a]+=1;su[b]+=1;jg[i]=1;
    }
    for(int i=1;i<=n;i++)
    {
        int s0=0,s1=0;
        for(int x:in[i])
        {
            bk[su[x]]=1;
            (zz[x]?s1:s0)+=1;
        }
        int w=su[i]-s0;
        while(bk[w])w+=1;
        int c=w-su[i];su[i]=w;
        for(int s=0;s<in[i].size();s++)
        {
            int x=in[i][s],y=vw[i][s];
            bk[su[x]]=0;
            if(c<0)
            {
                if(zz[x]==0)
                    jg[y]-=1,c+=1,zz[x]=1;
            }
            else if(c>0)
            {
                if(zz[x]==1)
                    jg[y]+=1,c-=1,zz[x]=0;
            }
        }
    }
    int cn=0;
    for(int i=1;i<=n;i++)cn+=zz[i];
    printf("%d\n",cn);
    for(int i=1;i<=n;i++)
        if(zz[i])printf("%d ",i);
    printf("\n");
    for(int i=0;i<m;i++)
        printf("%d %d %d\n",A[i],B[i],jg[i]);
    return 0;
}
//I love Set forever!