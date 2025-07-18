#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
template<class T>bool chmax(T &a, const T &b) { if (a<b) { a=b; return true; } return false; }
template<class T>bool chmin(T &a, const T &b) { if (b<a) { a=b; return true; } return false; }
#define all(x) (x).begin(),(x).end()
#define fi first
#define se second
#define mp make_pair
#define si(x) int(x.size())
const int mod=998244353,MAX=5005,INF=1<<30;

int main(){
    
    std::ifstream in("text.txt");
    std::cin.rdbuf(in.rdbuf());
    cin.tie(0);
    ios::sync_with_stdio(false);
    
    int Q;cin>>Q;
    while(Q--){
        int N,K;cin>>N>>K;
        if(K==0||K==1){
            for(int i=0;i<N;i++){
                if(i) cout<<" ";
                cout<<i+1;
            }
            cout<<"\n";
            cout<<N<<"\n";
            for(int i=0;i<N;i++){
                if(i) cout<<" ";
                cout<<i+1;
            }
            cout<<"\n";
        }else{
            vector<int> ans;
            for(int i=0;i<N;i+=K){
                int s=i,t=min(N-1,i+K-1);
                int len=t-s+1;
                if(len==1) ans.push_back(s);
                else if(len==2){
                    ans.push_back(s);
                    ans.push_back(t);
                }else if(len&1){
                    int x=len/2,y=len-x;
                    for(int j=s+x-1;j>=s;j--) ans.push_back(j);
                    for(int j=t;j>s+x-1;j--) ans.push_back(j);
                }else{
                    int x=len/2-1,y=len-x;
                    for(int j=s+x-1;j>=s;j--) ans.push_back(j);
                    for(int j=t;j>s+x-1;j--) ans.push_back(j);
                }
            }
            for(int i=0;i<N;i++){
                if(i) cout<<" ";
                cout<<ans[i]+1;
            }
            cout<<"\n";
            cout<<(N+K-1)/K<<"\n";
            for(int i=0;i<N;i++){
                if(i) cout<<" ";
                cout<<i/K+1;
            }
            cout<<"\n";
        }
    }
    
    /*
    for(int N=1;N<=10;N++){
        pair<int,vector<int>> MI=mp(INF,vector<int>{});
        vector<int> P(N);iota(all(P),0);
        
        do{
            int ma=0;
            for(int i=0;i<N;i++){
                for(int j=i+1;j<N;j++){
                    chmax(ma,abs(i-j)+abs(P[i]-P[j]));
                }
            }
            chmin(MI,mp(ma,P));
            
        }while(next_permutation(all(P)));
        
        cout<<N<<" "<<MI.fi<<" ";
        for(int x:MI.se) cout<<x+1<<" ";
        cout<<endl;
    }
     */
}