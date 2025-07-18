#include <iostream>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
using namespace std;
const int maxn=5005,maxm=200005,oo=0xfffffff;
int f[maxn],d[maxn],u[maxn],vis[maxn],num[maxn];
int pos[maxn];
int vv[maxn];
int n,m,s,p,e=0,limit;
int getf(int x)
{
    if(f[x]==x)return x;
    f[x]=getf(f[x]);
    return f[x];
}
int first[maxn],tot;
struct node{
       int u,v,c,b,q,tip,next;
       }a[maxm],g[maxm];
void add(int u,int v,int c,int q)
{
     ++tot;g[tot].u=u;g[tot].v=v;g[tot].q=q;g[tot].c=c;g[tot].next=first[u];g[tot].b=tot+1;first[u]=tot;
     ++tot;g[tot].u=v;g[tot].v=u;g[tot].q=q;g[tot].c=c;g[tot].next=first[v];g[tot].b=tot-1;first[v]=tot;
}
int cmp(const void *A,const void *B)
{
    return (*(node *)A).c-(*(node *)B).c;
}
void init()
{
     scanf("%d%d%d",&n,&m,&limit);
     int i;
     for(i=1;i<=m;i++)
     {scanf("%d%d%d",&a[i].u,&a[i].v,&a[i].c);a[i].q=i;}
     for(i=1;i<=m;i++)
     {
     if(a[i].u==1){u[a[i].v]=a[i].c;num[a[i].v]=i;}
     if(a[i].v==1){u[a[i].u]=a[i].c;num[a[i].u]=i;}
     }
     qsort(a+1,m,sizeof(a[1]),cmp);
     for(i=1;i<=n;i++)f[i]=i;
}
void mst()
{
     int i,f1,f2,rp;
     for(i=1;i<=m;i++){
                       if(a[i].u==1||a[i].v==1)continue;
                       f1=getf(a[i].u);
                       f2=getf(a[i].v);
                       if(f1==f2)continue;
                       s+=a[i].c;
                       if(f1>f2){rp=f1;f1=f2;f2=rp;}
                       a[i].tip=1;
                       f[f2]=f1;
                       }
     for(i=1;i<=m;i++)
     {
     if(a[i].u==1){
                   f1=getf(a[i].v);
                   if(f1==1)continue;
                   s+=a[i].c;
                   a[i].tip=1;
                   vis[a[i].v]=1;
                   f[f1]=1;
                   ++e;
                   }
     else if(a[i].v==1){
                   f1=getf(a[i].u);
                   if(f1==1)continue;
                   s+=a[i].c;
                   a[i].tip=1;
                   vis[a[i].u]=1;
                   f[f1]=1;
                   ++e;
                   }
     if(e==limit)break;
     }
     for(i=1;i<=m;i++)
     if(a[i].tip&&a[i].u!=1&&a[i].v!=1)
     add(a[i].u,a[i].v,a[i].c,a[i].q);
}
void dfs1(int x,int s,int l)
{
     int i,k;
     vv[x]=1;
     d[x]=s;
     pos[x]=l;
     for(i=first[x];i!=0;i=g[i].next)
     if(!vv[g[i].v]&&g[i].tip==0)
     {
     k=l;if(g[i].c>s)k=i;
     dfs1(g[i].v,max(s,g[i].c),k);
     }
}
void work()
{
     int i,k,mins,pp;
     memset(vv,0,sizeof(vv));
     for(i=1;i<=n;i++)if(vis[i])dfs1(i,0,0);
     while(e<limit){
                    k=0;
                    mins=oo;
                    for(i=2;i<=n;i++)
                    {
                                     if(vis[i]||!u[i])continue;
                                     if(u[i]-d[i]<=mins){
                                                        mins=u[i]-d[i];
                                                        k=i;
                                                        }
                    }
                    if(k==0)break;
                    s+=mins;
                    memset(vv,0,sizeof(vv));
                    pp=pos[k];
                    g[pp].tip=g[g[pp].b].tip=1;
                    vis[k]=1;
                    d[k]=0;
                    memset(vv,0,sizeof(vv));
                    for(i=1;i<=n;i++)if(vis[i])dfs1(i,0,0);
                    ++e;
                    }
}
bool output1()
{
     int i,f1;
     for(i=1;i<=n;i++){
                       f1=getf(i);
                       if(f1!=1)return 1;
                       }
     return 0;
}
void output()
{
     if(e!=limit){printf("-1\n");return ;}
     int i;
     printf("%d\n",n-1);
     for(i=2;i<=n;i++)
     if(vis[i])printf("%d ",num[i]);
     for(i=1;i<=tot;i++)
     {
     if(g[i].tip)continue;
     printf("%d ",g[i].q);
     g[i].tip=g[g[i].b].tip=1;
     }
}
int main(void)
{
    init();
    mst();
    if(output1()){printf("-1\n");return 0;}
    work();
    output();
    return 0;
}