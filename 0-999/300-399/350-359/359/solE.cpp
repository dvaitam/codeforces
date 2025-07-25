#include "stdio.h"
#include "iostream"
#include "vector"
using namespace std;
const int MAX = 505;
int a[MAX][MAX]={0}, visit[MAX][MAX]={0};
int n = 0, total = 0;
char opt[3000000]={0};
int dx[4]={-1,1,0,0};
int dy[4]={0,0,-1,1};
char opts[4]={'U','D','L','R'};
int cnt = 0;
bool canGo(int x,int y, int dir)
{
    x+=dx[dir];
    y+=dy[dir];
    while(x>0 && y>0 && x<=n && y<=n){
        if(visit[x][y])return false;
        if(a[x][y])return 1;
        x+=dx[dir];
        y+=dy[dir];
    }
    return false;
}
void dfs(int x, int y)
{
    visit[x][y] = 1;
    if(a[x][y]==0){
        total ++;
        opt[cnt++] = '1';
    }
    for(int i=0;i<4;i++){
        if(canGo(x,y,i)){
            opt[cnt++] = opts[i];
            dfs(x+dx[i],y+dy[i]);
            opt[cnt++] = opts[i^1];
        }
    }
    opt[cnt++] = '2';
    total--;
    a[x][y] = 0;
    
}
int main()
{
    int x0, y0;
    cin>>n>>x0>>y0;;
    for(int i=1;i<=n;i++)
        for(int j=1;j<=n;j++){
            scanf("%d", &a[i][j]);
            total += a[i][j];
        }
    dfs(x0,y0);
    if(total>0)cout<<"NO"<<endl;
    else{
        cout<<"YES"<<endl;
        printf("%s\n",opt);
    }
}