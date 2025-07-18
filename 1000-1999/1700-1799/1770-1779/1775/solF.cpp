#include<bits/stdc++.h>

using namespace std;



const int SQRTN = 1000;

long long dp[SQRTN][SQRTN],binom[SQRTN][SQRTN];



long long h(long long x){

    return binom[x+3][3];

}



long long g(long long num,long long siz,long long mod){

    if(num==0)return 1;

    if(siz==SQRTN)return 0;

    if(dp[num][siz]!=-1)return dp[num][siz];

    

    long long res=0;

    for(int i=0;i*siz<=num;i++){

        res+=g(num-i*siz,siz+1,mod)*h(i);

        res%=mod;

    }

    

    return dp[num][siz]=res;

}



long long f(int x,int y,int n,int mod){

    assert(x*y-n<SQRTN);

    return g(x*y-n,1,mod);

}



void solve(int mod){

    long long n;

    cin>>n;

    

    long long low=0,up=1500;

    

    while(up-low>1){

        long long mid = (up+low)/2;

        

        long long a = mid/2;

        long long b = mid-a;

        

        if(a*b>=n)up=mid;

        else low=mid;

    }

    long long a=up/2;

    long long b=up-a;

    

    if(mod==0){

        cout<<a<<" "<<b<<"\n";

        for(int i=0;i<a;i++){

            for(int j=0;j<b;j++){

                if(n>0){

                    cout<<"#";

                    n--;

                }else{

                    cout<<".";

                }

            }

            cout<<"\n";

        }

    }else{

        cout<<(a+b)*2<<" ";

        long long ans=0;

        

        for(long long i = a;(up-i)>=1 && i*(up-i)>=n;i++){

            ans+=f(i,up-i,n,mod);

            ans%=mod;

        }

        

        for(long long i = a-1;i>=1 && i*(up-i)>=n;i--){

            ans+=f(i,up-i,n,mod);

            ans%=mod;

        }

        

        cout<<ans<<"\n";

    }

}



int main(){

    ios::sync_with_stdio(false);

    cin.tie(0);

    

    memset(dp,-1,sizeof(dp));

    

    int t=1,u,m=0;

    cin>>t>>u;

    if(u==2){

        cin>>m;

        

        for(int i=0;i<SQRTN;i++){

            for(int j=0;j<=i;j++){

                if(j==0 || j==i)binom[i][j]=1;

                else binom[i][j]=(binom[i-1][j]+binom[i-1][j-1])%m;

            }

        }

    }

    for(int i=1;i<=t;i++)solve(m);

    

    return 0;

}