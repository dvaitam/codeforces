#include <cstdio>

#include <cstdlib>

#include <cstring>

#include <string>

#include <map>

#include <vector>

#include <iostream>

#include <set>

#include <algorithm>
#pragma comment(linker, "/STACK:67108864")
#define maxn 100010
#define maxm 601000

using namespace std;


int nxt[maxm],en[maxm],done[maxm],tot;
int first[maxn],dep[maxn],low[maxn],path[maxn];
bool visit[maxn],flag;
pair<int,int> ans[maxm];

void tjb(int x,int y){
    nxt[++tot]=first[x];
    first[x]=tot;
    en[tot]=y;
}

void dfs(int x){
    visit[x]=true;
    //cout<<x<<endl;
    int k,j,temp=0,num;
    low[x]=dep[x];
    k=first[x];
    while(k!=0){
        j=en[k];
        num=(k+1)/2;
        if(!visit[j]){
            if(!done[num]){
                done[num]=true;
                ans[num]=make_pair(x,j);
            }
            ++temp;
            dep[j]=dep[x]+1;
            path[j]=x;
            dfs(j);
            low[x]=min(low[x],low[j]);
            if(low[j]>dep[x])flag=true;
        }else if(j!=path[x]){
            low[x]=min(dep[j],low[x]);
            if(!done[num]){
                ans[num]=make_pair(x,j);
                done[num]=true;
            }
        }
        k=nxt[k];
    }
    //if(x==1 && temp>1)flag=true;
}
    

int main(){
    int n,m,x,y,i;

    scanf("%d%d",&n,&m);
    for(i=1;i<=m;++i){
        scanf("%d%d",&x,&y);
        tjb(x,y);
        tjb(y,x);
    }
    flag=false;
    dfs(1);
    if(flag)printf("0\n");else for(i=1;i<=m;++i)printf("%d %d\n",ans[i].first,ans[i].second);

    return 0;

}