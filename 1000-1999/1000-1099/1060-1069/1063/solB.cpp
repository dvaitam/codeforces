#include<iostream>
#include<cmath>
#include<cstring>
#include<iomanip>
#include<cstdio>
#include<algorithm>
#include<queue>
#include<set>
#include<map>

using namespace std;

int n,m;
int sr, sc;
int x,y;
char maze[2010][2010];

struct cell {
    int r;
    int c;
    int lsteps;
    int rsteps;
};


int main() {
    ios::sync_with_stdio(false);
    scanf("%d %d", &n, &m);
    scanf("%d %d", &sr, &sc);
    scanf("%d %d", &x, &y);
    for(int i=0;i<n;i++) {
        scanf("%s", maze[i]);
    }
    int ans=0;
    queue<cell> q;
    cell start;
    start.r=sr-1;
    start.c=sc-1;
    start.lsteps=start.rsteps=0;
    q.push(start);
    maze[sr-1][sc-1]='v';
    ans=1;
    while(!q.empty()) {
        cell cur=q.front();
        q.pop();
        int r=cur.r;
        int c=cur.c;
        int lsteps=cur.lsteps;
        int rsteps=cur.rsteps;
        if(c>0 && lsteps<x && maze[r][c-1]=='.') {
            ans++;
            cell t;
            t.r=r;
            t.c=c-1;
            t.lsteps=lsteps+1;
            t.rsteps=rsteps;
            q.push(t);
            maze[r][c-1]='v';
        }
        if(c<m-1 && rsteps<y && maze[r][c+1]=='.') {
            ans++;
            cell t;
            t.r=r;
            t.c=c+1;
            t.lsteps=lsteps;
            t.rsteps=rsteps+1;
            q.push(t);
            maze[r][c+1]='v';
        }
        for(int i=r-1;i>=0 && maze[i][c]=='.';i--) {
            ans++;
            maze[i][c]='v';
            if(c>0 && lsteps<x && maze[i][c-1]=='.') {
                ans++;
                cell t;
                t.r=i;
                t.c=c-1;
                t.lsteps=lsteps+1;
                t.rsteps=rsteps;
                q.push(t);
                maze[i][c-1]='v';
            }
            if(c<m-1 && rsteps<y && maze[i][c+1]=='.') {
                ans++;
                cell t;
                t.r=i;
                t.c=c+1;
                t.lsteps=lsteps;
                t.rsteps=rsteps+1;
                q.push(t);
                maze[i][c+1]='v';
            }
        }
        for(int i=r+1;i<n && maze[i][c]=='.';i++) {
            ans++;
            maze[i][c]='v';
            if(c>0 && lsteps<x && maze[i][c-1]=='.') {
                ans++;
                cell t;
                t.r=i;
                t.c=c-1;
                t.lsteps=lsteps+1;
                t.rsteps=rsteps;
                q.push(t);
                maze[i][c-1]='v';
            }
            if(c<m-1 && rsteps<y && maze[i][c+1]=='.') {
                ans++;
                cell t;
                t.r=i;
                t.c=c+1;
                t.lsteps=lsteps;
                t.rsteps=rsteps+1;
                q.push(t);
                maze[i][c+1]='v';
            }
        }
    }
    printf("%d\n", ans);
}