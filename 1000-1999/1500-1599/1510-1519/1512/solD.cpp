#include<bits/stdc++.h>

using namespace std;

#define f(a,i,n) for(int i=a;i<n;i++)

#define f1(a,i,n) for(int i=a;i<=n;i++)

#define int long long

#define ll long long

#define vi vector<int>

#define vll vector<long long>

#define vs vector<string>

#define vvi vector<vector<int>>

#define pii pair<int,int>

#define vpii vector<pii>

#define pb push_back

#define pf push_front

#define mci map<char,int>

#define umci unordered_map<char,int>

#define mii map<int,int>

#define umii unordered_map<int,int>

#define all(a) a.begin(),a.end()

#define ret return

const int INF =1e9+7;

const int N =1e5+9;

void solve(){

int n;cin>>n;

    vi b(n+2);

    int hi=0,sechi=0;

    int j,k;

    int sum=0;

    f(0,i,n+2){

        cin>>b[i];

        sum+=b[i];

        if(hi<b[i]){

            hi=b[i];

            j=i;

        }

    }

    f(0,i,n+2){

        if(i!=j){

            if(sechi<b[i]){

                sechi=b[i];

                k=i;

            }

        }

    }

    int cnh=0;

    f(0,i,n+2){

        if(b[i]>=sum-hi-sechi){

            cnh++;

        }

    }

    

    if(sum-hi==2*sechi and cnh){

        f(0,i,n+2){

            if(i!=j and i!=k){

                cout<<b[i]<<" ";

            }

        }

        cout<<endl;

        ret;

    }

    else if(cnh){

        int p,cn=0;

        f(0,i,n+2){

            if(sum-b[i]==2*hi){

                p=i;

                cn++;

                break;

            }

        }

        if(cn){



        f(0,i,n+2){

            if(i!=p and i!=j){

            cout<<b[i]<<" ";

            }

        }

        cout<<endl;

        ret;

        }

        else{

            cout<<-1<<endl;

              ret;

        }

    }

    else{

        cout<<-1<<endl;

    }

}

main(){

ios_base::sync_with_stdio(false);

cin.tie(0);

cout.tie(0);

int t; cin>>t;

while(t--){

    solve();

    

}

return 0;

}