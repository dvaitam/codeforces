#include<cstdio>
#include<algorithm>
#include<cstring>
#include<cstdlib>
using namespace std;
char a[4][15]=
{
    "AAA.A.A....A",
    ".A..A.AAAAAA",
    ".A.AAAA....A"
};
char p[10][10],ansp[10][10];
int n,m,ans,rem;
void saya(int x,int y,int move)
{
    if (x>=n-2)
    {
        ans=max(ans,move);
        if (move==ans)
            memcpy(ansp,p,sizeof(p));
        return;
    }

    if (y>=m-2)
    {
        rem-=p[x][y]=='.';
        rem-=p[x][y+1]=='.';
        saya(x+1,0,move);
        rem+=p[x][y]=='.';
        rem+=p[x][y+1]=='.';
        return;
    }
    if (rem/5<=ans-move) return;
    rem-=p[x][y]=='.';
    for(int d=0;d<12;d+=3)
    {
        int flag=0;
        for(int i=0;i<3;i++)
            for(int j=0;j<3;j++)
            {
                if (a[i][d+j]=='A'&&p[x+i][y+j]!='.')
                {
                    flag=1;
                    goto slb;
                }
            }
slb:    if (!flag)
        {
            for(int i=0;i<3;i++)
                for(int j=0;j<3;j++)
                    if (a[i][d+j]=='A') p[x+i][y+j]=a[i][d+j]+move;
            rem-=5;
            saya(x,y+1,move+1);
            rem+=5;
            for(int i=0;i<3;i++)
                for(int j=0;j<3;j++)
                    if (a[i][d+j]=='A') p[x+i][y+j]='.';
        }
    }
    saya(x,y+1,move);
    rem+=p[x][y]=='.';
}
int main()
{
    scanf("%d%d",&n,&m);
    for(int i=0;i<n;i++)
        for(int j=0;j<m;j++) ansp[i][j]=p[i][j]='.';
    rem=n*m;
    saya(0,0,0);
    printf("%d\n",ans);
    for(int i=0;i<n;i++) puts(ansp[i]);
}