#include<cstdio>

#include<queue>

#include<cstring>

#include<algorithm>

using namespace std;

struct node{

    int x,y,val;

}str[3010]; //储存每个湖泊的面积，以及他的其中一个坐标

int N,M,K,num,flag;

char st[55][55];

int fx[4] = {0,0,-1,1},fy[4] = {1,-1,0,0},vis[55][55];

bool cmp(node i,node j){

    return i.val > j.val;

}

void DFS(int x,int y) // 遍历该湖泊所有的点，以及他的最大面积

{

    if(vis[x][y] || st[x][y] == '*')

        return ;

    num++;

    vis[x][y] = 1;

    for(int i = 0 ; i < 4 ; i++){

        int nx = x + fx[i];

        int ny = y + fy[i];

        if(nx < 1 || ny < 1 || nx > N || ny > M){

           num = 0,flag  = 0;continue;

        }

        if(st[nx][ny] == '*' || vis[nx][ny])

            continue;

        DFS(nx,ny);

    }

}

void DFSL(int x,int y){ //把该湖泊变成陆地

     st[x][y] = '*';

     for(int i = 0 ; i < 4 ; i++){

        int nx = x + fx[i];

        int ny = y + fy[i];

        if(st[nx][ny] == '*')

            continue;

        DFSL(nx,ny);

     }

}

int main()

{

    int i,j,pl = 0;

    scanf("%d%d%d",&N,&M,&K);

    for(i = 1 ; i <= N ; i++)

        scanf("%s",st[i] + 1);

    for(i = 1 ; i <= N ; i++)

        for(j = 1 ; j <= M ; j++){

            num = 0; flag = 1;

            if(!vis[i][j] && st[i][j] == '.'){ // 遍历每个湖泊，以及记录该湖泊的其中一个位置和湖泊总个数

                    DFS(i,j);

                    if(flag){

                        str[++pl].x = i;

                        str[pl].y = j;

                        str[pl].val = num;

                    }

            }

        }

    int  ans = 0;

    sort(str + 1 , str + 1 + pl, cmp); // 降序排列

    for(i =  K + 1 ; i <= pl ; i++){ // 前 K 个湖泊不用改

        ans += str[i].val;

        DFSL(str[i].x,str[i].y);

    }

    printf("%d\n",ans);

    for(i = 1 ; i <= N ; i++)

        printf("%s\n",st[i] + 1);

    return 0;

}