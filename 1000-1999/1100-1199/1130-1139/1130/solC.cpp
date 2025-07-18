#include<bits/stdc++.h>
using namespace std;
int n,sx,sy,ex,ey;
int nxt[5][5];
char mp[59][59];
bool used[59][59],f=0;
struct node
{
    int x,y;
};
vector<node> vec1,vec2;
void dfs1(int x,int y)
{
    used[x][y]=1;
    if(x==ex&&y==ey)
    {
        f=1;
        return;
    }
    for(int i=1; i<=4; i++)
    {
        int tx=x+nxt[i][1];
        int ty=y+nxt[i][2];
        if(used[tx][ty]||tx>n||ty>n||tx<1||ty<1)
            continue;
        if(mp[tx][ty]=='1')
        {
            node temp;
            temp.x=x;temp.y=y;
            vec1.push_back(temp);
            continue;
        }
        dfs1(tx,ty);
    }
}
void dfs2(int x,int y)
{
    used[x][y]=1;
    for(int i=1; i<=4; i++)
    {
        int tx=x+nxt[i][1];
        int ty=y+nxt[i][2];
        if(used[tx][ty]||tx>n||ty>n||tx<1||ty<1)
            continue;
        if(mp[tx][ty]=='1')
        {
            node temp;
            temp.x=x;temp.y=y;
            vec2.push_back(temp);
            continue;
        }
        dfs2(tx,ty);
    }
}
int sqr(int x)
{
    return x*x;
}
int main()
{
    cin>>n;
    cin>>sx>>sy;
    cin>>ex>>ey;
    nxt[1][1]=1;
    nxt[2][1]=-1;
    nxt[3][1]=0;
    nxt[4][1]=0;
    nxt[1][2]=0;
    nxt[2][2]=0;
    nxt[3][2]=1;
    nxt[4][2]=-1;
    for(int i=1; i<=n; i++)
        for(int j=1; j<=n; j++)
            cin>>mp[i][j];
    dfs1(sx,sy);
    if(f)
        return puts("0"),0;
    memset(used,0,sizeof(used));
    dfs2(ex,ey);
    int ans=0x3f3f3f3f;
    for(int i=0;i<vec1.size();i++)
    {
        for(int j=0;j<vec2.size();j++)
        {
            node u=vec1[i],v=vec2[j];
            ans=min(ans,sqr(u.x-v.x)+sqr(u.y-v.y));
        }
    }
    printf("%d\n",ans);
    return 0;
}