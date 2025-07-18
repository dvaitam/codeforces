#include<cstdio>
#include<cstdlib>
#include<iostream>
#include<vector>
#include<string>
#include<queue>
#include<stack>
#include<cstring>
#include<map>
#include<set>
#include<algorithm>
using namespace std;

#define INF 1000000000
#define NUM 10000

typedef pair<int,int> ii;
typedef vector<int> vi;
typedef vector<ii> vii;

int main(){

    set<int> lin,col;
    int n,m,cont=0;
    char c;
    bool poss=true;
    ii ans(-1,-1);

    cin >> n >> m;
    getchar();
    ii M[n+m+5];
    for(int i=0;i<n && poss;i++){
        for(int j=0;j<m && poss;j++){
            c=getchar();
            if(c=='*'){
                M[cont++]=ii(i,j);
                if(cont>n+m-1) poss=false;
            }
        }
        getchar();
    }

    for(int i=0;i<n && ans.first==-1;i++)
        for(int j=0;j<m && ans.first==-1;j++){
            bool util=true;
            for(int k=0;k<cont && util;k++){
                if(M[k].first!=i && M[k].second != j) util=false;
            }
            if(util) ans=ii(i,j);
        }

    if(ans.first!=-1){
        cout<<"YES"<<endl;
        cout<<ans.first+1<<" "<<ans.second+1<<endl;
    }else
        cout<<"NO"<<endl;



return 0;
}