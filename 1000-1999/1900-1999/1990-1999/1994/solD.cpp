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
const int mod=998244353,MAX=300005,INF=15<<26;

struct UF{
    int n;
    vector<int> par,size,edge;
    
    void init(int n_){
        n=n_;
        par.assign(n,-1);
        size.assign(n,1);
        edge.assign(n,0);
        
        for(int i=0;i<n;i++){
            par[i]=i;
        }
    }
    
    int root(int a){
        if(par[a]==a) return a;
        else return par[a]=root(par[a]);
    }
    
    void unite(int a,int b){
        edge[root(a)]++;
        if(root(a)!=root(b)){
            size[root(a)]+=size[root(b)];
            edge[root(a)]+=edge[root(b)];
            par[root(b)]=root(a);
        }
    }
    
    bool check(int a,int b){
        return root(a)==root(b);
    }
};

int main(){
    
    std::ifstream in("text.txt");
    std::cin.rdbuf(in.rdbuf());
    cin.tie(0);
    ios::sync_with_stdio(false);
    
    int Q;cin>>Q;
    while(Q--){
        ll N;cin>>N;
        vector<ll> A(N);
        for(int i=0;i<N;i++) cin>>A[i];
        
        UF uf;uf.init(N);
        
        vector<pair<int,int>> ans;
        
        for(int t=N-1;t>=1;t--){
            vector<int> mita(t,-1);
            for(int i=0;i<N;i++){
                if(uf.root(i)==i){
                    if(mita[A[i]%t]==-1) mita[A[i]%t]=i;
                    else{
                        ans.push_back(mp(mita[A[i]%t],i));
                        uf.unite(mita[A[i]%t],i);
                        break;
                    }
                }
            }
        }
        
        reverse(all(ans));
        cout<<"Yes\n";
        for(auto [a,b]:ans) cout<<a+1<<" "<<b+1<<"\n";
    }
}