#include <bits/stdc++.h>

using namespace std;

typedef long long ll;

#define speed      ios::sync_with_stdio(false);cin.tie(NULL); cout.tie(NULL);

#define rng(i,a,b) for (int i=a;i>=int(b);i--)

#define rep(i,a,b) for(int i=a;i<=int(b);++i)

#define lop(i,n)   for(int i=0;i<int(n);i++)

#define cler(x,a)  memset(x,a,sizeof(x))

#define all(x)     (x).begin(),(x).end()

#define pb(x)      push_back(x)

#define sz(x)      (x.size())

#define f first

#define s second

#define sc(x) scanf("%d",&x)

#define pr(x) printf("%d ",x)

#define prln(x) printf("%d\n",x)

/*

#define scc(x) scanf("%1c",&x)

#define scl(x) scanf("%lld",&x)

#define prl(x) printf("%I64d\n",x)

typedef pair<int,int> ii;

typedef vector <int>  vi;

ll MOD=1000003;

freopen("output.txt","w",stdout);

freopen("input.txt","r",stdin);

*/

/*

for (int i=0;i<n;i++)            Hi    ^_^  hack me if you can :P



*/

int n;

int c[51];

int a[51][51];

int mn,mnx,mny;

vector<int> x,y,inx,iny;

int ans;

int main()

{

    sc(n);



    rep(i,1,n)

    sc(c[i]);



    rep(i,1,n)

    rep(j,1,c[i])

    sc(a[i][j]);



    rep(i,1,n)

    {

        rep(j,1,c[i])

        {

            mn=a[i][j];

            mnx=i,mny=j;

            rep(ii,i,n)

            rep(jj,j,c[ii])

            {

                if(a[ii][jj]<mn)

                {

                    mn=a[ii][jj];

                    mnx=ii;

                    mny=jj;

                }

            }

            if(mn^a[i][j])

            //else

            {

                ans++;

                x.pb(i);

                y.pb(j);

                inx.pb(mnx);

                iny.pb(mny);

                swap(a[i][j],a[mnx][mny]);

            }

        }

    }

    prln(ans);

    lop(i,ans)

    {

        pr(x[i]);

        pr(y[i]);

        pr(inx[i]);

        prln(iny[i]);

    }

    return 0;

}

/*



int dx[]= {0,0 ,1,-1, 1,1,-1,-1};

int dy[]= {1,-1,0,0 , 1,-1,1,-1};

char a[101][101];

bool vis[101][101];

int n,m,ans;

bool in (int i,int j) {return 0<=i&&0<=j&&i<n&&j<m ; }

void dfs(int i,int j){

    vis[i][j]=1;

    for (int k=0;k<4;k++){

        int x=i+dx[k];

        int y=j+dy[k];

        if(in(x,y))

        if(!vis[x][y])

        if(a[x][y]=='B')

        dfs(x,y);

    }



}







ll powmod(ll a,ll b)

{

    ll res=1LL;

    a%=MOD;

    for(; b; b>>=1)

    {

        if(b&1)res=res*a%MOD;

        a=a*a%MOD;

    }

    return res;

}

ll pow(long double a,ll b)

{

    ll res=1LL;



    for(; b; b>>=1)

    {

        if(b&1)res=res*a;

        a=a*a;

    }

    return res;

}



 fflush(stdout);

 cout.flush();





    ll A,B,n,x;

    scanf("%I64d%I64d%I64d%I64d",&A,&B,&n,&x);

    if(A==1)printf("%I64d",(x+n%MOD*B)%MOD);

    else

    {

        ll res=fp(A,n)*x%MOD;

        res+=(fp(A,n)-1)*fp(A-1,MOD-2)%MOD*B;

        printf("%I64d",(res%MOD+MOD)%MOD);

    }





bool cmp(string x,string y)   {return x+y<y+x;}

int dx[]= {0,0 ,1,-1, 1,1,-1,-1};

int dy[]= {1,-1,0,0 , 1,-1,1,-1};

bool inside (int i,int j) {return (i>=1 && i<=n && j>=1 && j<=m);}





bool prime[100007];

void sieve()

{

    for(int i=4; i<=100005; i+=2)

        prime[i]=false,prime[i-1]=true;

    prime[2]=true;

    for(int i=3; i<=1000; i+=2)

    {

        if(prime[i])

        {

            for(int j=i*i; j<=100005; j+=i*2)

                prime[j]=false;



        }



    }



}



    ll a,b,n,x,na,nb;

    cin>>a>>b>>n>>x;

    while(n){

        if (n&1){

            x=(a*x+b)%MOD;

        }

        n>>=1;

        na=a*a%MOD;

        nb=(a*b+b)%MOD;

        a=na;

        b=nb;

    }

    cout<<x;



*/