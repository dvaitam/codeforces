#include<bits/stdc++.h>
using namespace std;

int n;

void solve(){
    cin>>n;
    if(n&1)cout<<"NO"<<endl;
    else{
        cout<<"YES"<<endl;
        for(int i=1;i<=n/2;i++){
            if(i&1)cout<<"AA";
            else cout<<"BB";
        }
        cout<<endl;
    }
}

int main(){
    int t=1;
    cin>>t;
    while(t--){
        solve();
    }
    
    
    return 0;
}