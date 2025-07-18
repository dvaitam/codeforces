#include <bits/stdc++.h>
#define mp make_pair
#define pb push_back
#define F first
#define S second
#define all(x) x.begin(),x.end()
#define MAXN 1005
typedef long long ll;

using namespace std;

int n,m;
int mat[MAXN][MAXN];
char s[MAXN];
int gore[MAXN][MAXN],dole[MAXN][MAXN],levo[MAXN][MAXN],desno[MAXN][MAXN],pox[MAXN][MAXN],poy[MAXN][MAXN];
int res[MAXN][MAXN];
bool b[MAXN][MAXN];


int main()
{
    scanf("%d %d", &n, &m);
    for(int i=1;i<=n;i++){
        scanf("%s",s);
        for(int j=1;j<=m;j++){
            if(s[j-1]=='.')mat[i][j]=0;
            else mat[i][j]=1;
        }
    }
    for(int i=1;i<=n;i++){
        for(int j=1;j<=m;j++){
            levo[i][j]=mat[i][j]*(levo[i][j-1]+mat[i][j]);
            gore[i][j]=mat[i][j]*(gore[i-1][j]+mat[i][j]);
        }
    }
    for(int i=n;i>=1;i--){
        for(int j=m;j>=1;j--){
            desno[i][j]=mat[i][j]*(desno[i][j+1]+mat[i][j]);
            dole[i][j]=mat[i][j]*(dole[i+1][j]+mat[i][j]);
        }
    }
    for(int i=1;i<=n;i++){
        for(int j=1;j<=m;j++){
            int domet=min(min(gore[i][j],dole[i][j]),min(levo[i][j],desno[i][j]));
            if(domet<=1)continue;
            res[i][j]=domet-1;
            //printf("   %d %d %d\n",i,j,res[i][j]);
            pox[i][max(1,j-domet+1)]++;
            pox[i][min(m+1,j+domet)]--;
            poy[max(1,i-domet+1)][j]++;
            poy[min(n+1,i+domet)][j]--;

        }
    }
    for(int i=1;i<=n;i++){
        int tren=0;
        for(int j=1;j<=m;j++){
            //printf("%d ",tren);
            tren+=pox[i][j];
            if(tren>0)b[i][j]=true;
        }
        //printf("\n");

    }
    for(int i=1;i<=n;i++){
        int tren=0;
        for(int j=1;j<=m;j++){
            tren+=pox[i][j];
            if(tren>0)b[i][j]=true;
        }
    }

    for(int j=1;j<=m;j++){
        int tren=0;
        for(int i=1;i<=n;i++){
            tren+=poy[i][j];
            if(tren>0)b[i][j]=true;
        }
    }
    for(int i=1;i<=n;i++){
        for(int j=1;j<=m;j++){
            if(!b[i][j] && mat[i][j]){
                printf("-1");
                return 0;
            }
        }
    }
    int broj=0;
    for(int i=1;i<=n;i++){
        for(int j=1;j<=m;j++){
            if(res[i][j]){
                broj++;

            }
        }
    }
    printf("%d\n",broj);

    for(int i=1;i<=n;i++){
        for(int j=1;j<=m;j++){
            if(res[i][j]){
                printf("%d %d %d\n",i,j,res[i][j]);

            }
        }
    }




    return 0;

}