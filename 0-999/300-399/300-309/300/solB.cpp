#include<stdio.h>
#include<math.h>
#include<string.h>
#include<algorithm>
#include<stdlib.h>
#include<math.h>
using namespace std;
#define maxn 1100000
#define inf 10000007
typedef long long ll;
int a[100];
int fa[100];
int num[100];
int find(int x)
{
    if(x!=fa[x])
    fa[x]=find(fa[x]);
    return fa[x];
}
void merge(int x,int y)
{
    int fx=find(x);
    int fy=find(y);
    if(fx!=fy)
    {
      fa[fx]=fy;
      num[fy]+=num[fx];
    }
}
int main()
{
    //freopen("in.txt","r",stdin);
    int n,m;
    while(scanf("%d%d",&n,&m)==2)
    {
        int i,j;
        int a,b;
        for(i=1;i<=n;i++)
        {
          fa[i]=i;
          num[i]=1;
        }
        for(i=0;i<m;i++)
        {
            scanf("%d%d",&a,&b);
            merge(a,b);
        }
        int c[50],d[50];
        int flag=0;
        int u=0,v=0;
        for(i=1;i<=n;i++)
        {
            int f=find(i);
            if(i==f)
            {
               if(num[f]>3)
               flag=1;
               if(num[f]==2)
               c[u++]=f;
               if(num[f]==1)
               d[v++]=f;
            }
        }
        //printf("%d %d\n",u,v);
        if(flag==1)
        printf("-1\n");
        else if(v>u&&(v-u)%3!=0)
        printf("-1\n");
        else if(v<u)
        printf("-1\n");
        else
        {
            if(u==v)
            {
            for(i=0;i<u;i++)
            {
                merge(c[i],d[i]);
            }
            }
            else if(v>u)
            {
            int x=(v-u)/3;
           // printf("x=%d\n",x);
            int id=0;
            for(i=0;i<x;i++)
            {
                int id1=id;
                id1++;
                //printf("%d %d\n",d[id],d[id1]);
                merge(d[id],d[id1]);
                int id2=id1;
                id2++;
                merge(d[id1],d[id2]);
                id+=3;
            }
           //printf("fax=%d\n",find(2));
            for(i=0;i<u;i++)
            {
                merge(c[i],d[id++]);
            }
            }
            int k=0;
            int g[50];
            for(i=1;i<=n;i++)
            {
                if(i==find(i))
                g[k++]=find(i);
            }
            //printf("k=%d\n",k);
            for(j=0;j<k;j++)
            {
              for(i=1;i<=n;i++)
              {
                  if(find(i)==g[j])
                  printf("%d ",i);
              }
              printf("\n");
            }
        }
    }
    return 0;
}