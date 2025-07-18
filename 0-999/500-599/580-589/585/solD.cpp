#include<cstdio>
#include<cstring>
#include<algorithm>
const int maxn=30,mod=9999993,base=500000000,maxt=600010,inf=1e9;
typedef long long ll;
using namespace std;
int n,m,ans,sta1,sta2,pow[maxn];
struct data{int a,b,c;}p[maxn];

struct hash{
    int pre[maxt],now[mod+10],tot,val[maxt][3],sta[maxt];
    void insert(int a,int b,int c,int st){
        int x=(a-b+base)%mod;
        for (int y=now[x];y;y=pre[y])
            if (val[y][0]==a&&val[y][1]==b&&val[y][2]==c) return;
        pre[++tot]=now[x],now[x]=tot,val[tot][0]=a,val[tot][1]=b,val[tot][2]=c,sta[tot]=st;
    }
    void query(int a,int b,int c,int st){
        int x=(b-a+base)%mod,tmp=-inf*2,t1;
        for (int y=now[x];y;y=pre[y])
            if (val[y][1]-val[y][2]==c-b)
                if (val[y][0]>tmp) tmp=val[y][0],t1=sta[y];
        if (tmp+a>ans) ans=tmp+a,sta1=st,sta2=t1;
    }
}T;

void dfs(int k,int lim,int a,int b,int c,int code,int op){
    if (k==lim+1){
        if (!op) T.insert(a,b,c,code);
        else T.query(a,b,c,code);
        return;
    }
 dfs(k+1,lim,a+p[k].a,b+p[k].b,c,code*3,op);
    dfs(k+1,lim,a+p[k].a,b,c+p[k].c,code*3+1,op);
    dfs(k+1,lim,a,b+p[k].b,c+p[k].c,code*3+2,op);
}

int main(){
    scanf("%d",&n),m=(n+1)>>1,ans=-inf;
    pow[0]=1;for (int i=1;i<=m+1;i++) pow[i]=pow[i-1]*3;
    for (int i=1;i<=n;i++) scanf("%d%d%d",&p[i].a,&p[i].b,&p[i].c);
    dfs(m+1,n,0,0,0,0,0),dfs(1,m,0,0,0,0,1);
    if (ans==-inf){puts("Impossible");}
    else{
        for (int i=1;i<=m;i++){
            int t=(sta1/pow[m-i])%3;
   if (t==0) puts("LM");
   if (t==1) puts("LW");
            if (t==2) puts("MW");
        }
        for (int i=m+1;i<=n;i++){
            int t=(sta2/pow[n-i])%3;
            if (t==0) puts("LM");
   if (t==1) puts("LW");
            if (t==2) puts("MW");
        }
    }
    return 0;
}
/*
25
26668 10412 12658
25216 11939 10247
28514 22515 5833
4955 19029 22405
12552 6903 19634
12315 1671 505
20848 9175 6060
12990 5827 16433
9184 30621 25596
31818 7826 11221
18090 4476 30078
30915 11014 16950
3119 29529 21390
775 4290 11723
29679 14840 3566
4491 29480 2079
24129 5496 6381
20849 25772 9299
10825 30424 11842
18290 14728 30342
24893 27064 11604
26248 7490 18116
17182 32158 12518
23145 4288 7754
18544 25694 18784

*/