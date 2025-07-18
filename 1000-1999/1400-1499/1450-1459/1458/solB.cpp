#include<bits/stdc++.h>

using namespace std;

#define AC return

#define Please 0

#define itn int

#define enld '\n'       //IO?

#define endl '\n'       //IO?

#define Endl '\n'       //IO?

#define esle else

#define ciN cin

#define YES "YES"

#define NO "NO"

#define eb emplace_back

#define wide(x,ch) setw(x)<<setfill(ch)

#define show "ans-------------> "<<

class MyCpiise{public:bool operator()(pair<int,int> a,pair<int,int> b){return a.second<b.second;}};

MyCpiise Mcse;

#define dbg(x...) do{cout<<#x<<" -> ";err(x);}while (0)

void err(){cout<<'\n';}

template<class T, class... Ts>

void err(T arg, Ts... args) {

    cout<<arg<< ' ';

    err(args...);

}



#define int long long                        //int main()->signed main()?

#define INF 0x3f3f3f3f

#define LLINF 9223372036854775807 

#define pii pair<int,int>

#define pip pair<int,pii >

#define ppi pair<pii ,int>

#define ppp pair<pii ,pii >

const int MOD1=998244353;

const int MOD2=1e9+7;

const int N4=1e4+10;

const int N5=1e5+10; 

const int NN5=2e5+10;

const int N6=1e6+7;

const int N7=1e7+10;

const int N8=1e8+10;

const int N9=1e9+5;

int dr[4]={0,1,0,-1};

int dc[4]={1,0,-1,0};

int dx[4]={1,0,-1,0};

int dy[4]={0,1,0,-1};

#define is_odd(x) (x&1)

#define lowbit(x) (x&(-x))

void bin_print(int now,int cnt){

    if(!cnt) return;

    bin_print(now>>1,cnt-1);

    cout<<(now&1);

}

int fast_pow(int x,int n=MOD2-2,int mod=MOD2,int ret=1){

    while(n){

        if(n&1) ret=ret*x%mod;

        x=x*x%mod;

        n>>=1;

    }

    return ret;

}

long long fast_ll_pow(long long x,long long n,long long ret=1)

{

    while(n){

        if(n&1) ret=ret*x;

        x=x*x;

        n>>=1;

    }

    return ret;

}

int exgcd(int a,int b,int &x,int &y){

    if(!b){x=1;y=0;return a;}

    else{int ret=exgcd(b,a%b,y,x);y-=a/b*x;return ret;}

}

int gcd(int a,int b){return !b ? a : gcd(b,a%b);}

typedef long long ll;

double dp[110][20010];   //dp[i][j][k],前i个杯子选择了j个杯子，总容积为k时的最大储水量，i滚动

int w[110];

double v[110];

vector<double>ans;

int sumw[110];

void run()

{

    int n;

    cin>>n;

    for(int i=1;i<=n;i++)

        cin>>w[i]>>v[i];

    memset(dp,-INF,sizeof(dp));

    dp[0][0]=0;

    // last[1][w[1]]=v[1];

    for(int i=1;i<=n;i++) sumw[i]=sumw[i-1]+w[i];

    double sumv=0;

    for(int i=1;i<=n;i++) sumv+=v[i];

    for(int i=1;i<=n;i++){

            for(int j=i;j>=1;j--){

                for(int k=10000;k>=w[i];k--){

                    dp[j][k]=max(dp[j][k],dp[j-1][k-w[i]]+v[i]);

                }

            }

        }

    // dbg("----------------------------------------------------------------------------------");

    for(int j=1;j<=n;j++)

    {

        double ans=0;

        for(int k=0;k<N4;k++) 

        {

            ans=max((double)ans,min((double)k,(dp[j][k]+sumv)*0.5));

            // dbg(j,k,last[j][k]);

        }

        cout<<ans<<' ';

    }

    

        





    return;

    //T==1?

}

signed main()

{

//    freopen("D:\\MY_std_input\\std_input.txt", "r", stdin);

//    freopen("D:\\MY_std_input\\std_output1.txt", "w", stdout);

//    freopen("D:\\MY_std_input\\std_output2.txt", "w", stdout);

    ios::sync_with_stdio(false);

    cin.tie(0);cout.tie(0);

    // system("cls");

    int T=1;

    // cin>>T;

    while(T--)

    {

        run();

    }







    AC Please;

}

// "str.size()-num" --> "(int)str.size()-num" ?

/*

    sort(edge[i].begin(), edge[i].end(), [&](pii &a, pii &b){

        ......

    });

*/