#include <cstdio>

using namespace std;
#define MAXN 510
#define MAXK 1000010
int n,m,len,k;
int g[4][2]={{0,1},{0,-1},{1,0},{-1,0}};
char Map[MAXN][MAXN];
bool vis[MAXN][MAXN],mark[MAXN][MAXN];
struct Point
{
    int x,y;
    char Opt;
    Point(int x=0,int y=0):x(x),y(y){}
}sta[MAXK],ans[MAXK];
void Init()
{
    scanf("%d%d",&n,&m);
    for(int i=1;i<=n;++i)
        scanf("%s",Map[i]+1);
}
void Solve()
{
    int i,j,p,x,y;
    Point a; bool flag;
    for(i=1;i<=n;++i)
        for(j=1;j<=m;++j){
            if(mark[i][j]||Map[i][j]=='#') continue;
            mark[i][j]=true;
            for(sta[++len]=Point(i,j);len;){
                a=sta[len];
                if(vis[a.x][a.y]){
                    if((--len)){
                        ans[++k]=a; ans[k].Opt='D';
                        ans[++k]=a; ans[k].Opt='R';
                    }
                }
                else{
                    vis[a.x][a.y]=true;
                    ans[++k]=a; ans[k].Opt='B';
                    flag=false;
                    for(p=0;p<4;++p){
                        x=a.x+g[p][0];
                        y=a.y+g[p][1];
                        if(x<1||x>n||y<1||y>m||mark[x][y]||Map[x][y]=='#')
                            continue;
                        flag=mark[x][y]=true;
                        sta[++len]=Point(x,y);
                    }
                    if(!flag){
                        if((--len)) ans[k].Opt='R';
                    }
                }
            }
        }
    printf("%d\n",k);
    for(i=1;i<=k;++i)
        printf("%c %d %d\n",ans[i].Opt,ans[i].x,ans[i].y);
}
int main()
{
    Init();
    Solve();
    return 0;
}