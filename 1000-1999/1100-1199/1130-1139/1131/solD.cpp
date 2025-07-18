#include<bits/stdc++.h>
#define mp make_pair
#define fi first
#define se second
#define debug(x) cerr<<#x<<" = "<<(x)<<endl
#define eps 1e-8
#define pi acos(-1.0)
using namespace std;
void test(){cerr<<"\n";}
template<typename T,typename... Args>void test(T x,Args... args){cerr<<x<<" ";test(args...);}
typedef long long ll;
typedef pair<int,int> pii;
typedef pair<ll,ll> pll;
const int MAXN=(int)1e5+5;
const int MOD=(int)1e9+7;
char g[1005][1005];
int a[1005],b[1005];
int main()
{
    int n,m;
    scanf("%d%d",&n,&m);
    vector<pii>v1(n);
    for(int i=0;i<n;i++){
        scanf("%s",g[i]);
        v1[i]={0,i};
        for(int j=0;j<m;j++){
            if(g[i][j]=='>')v1[i].fi++;
            else if(g[i][j]=='<')v1[i].fi--;
        }
    }
    sort(v1.begin(),v1.end());
    bool flag=1;
    a[v1[0].se]=1;
    for(int i=1;i<n;i++){
        int u=v1[i].se,v=v1[i-1].se;
        if(v1[i].fi==v1[i-1].fi)a[u]=a[v];
        else a[u]=a[v]+2;
    }
    for(int j=0;j<m;j++){
        int f=0;
        for(int i=0;i<n;i++){
            if(g[i][j]=='='){
                b[j]=a[i];
                f=1;
                break;
            }
        }
        if(!f){
            for(int i=0;i<n;i++){
                if(g[i][j]=='<'){
                    b[j]=max(b[j],a[i]+1);
                }
            }
        }
    }
    for(int i=0;i<n;i++){
        for(int j=0;j<m;j++){
            char tmp;
            if(a[i]>b[j])tmp='>';
            else if(a[i]==b[j])tmp='=';
            else tmp='<';
            if(tmp!=g[i][j])flag=0;
        }
    }
    if(flag==0)printf("No\n");
    else {
        printf("Yes\n");
        vector<int>ve;
        for(int i=0;i<n;i++){
            ve.push_back(a[i]);
        }
        for(int i=0;i<m;i++){
            ve.push_back(b[i]);
        }
        sort(ve.begin(),ve.end());
        ve.erase(unique(ve.begin(),ve.end()),ve.end());
        for(int i=0;i<n;i++){
            a[i]=lower_bound(ve.begin(),ve.end(),a[i])-ve.begin()+1;
            printf("%d ",a[i]);
        }
        printf("\n");
        for(int i=0;i<m;i++){
            b[i]=lower_bound(ve.begin(),ve.end(),b[i])-ve.begin()+1;
            printf("%d ",b[i]);
        }
    }
    return 0;
}