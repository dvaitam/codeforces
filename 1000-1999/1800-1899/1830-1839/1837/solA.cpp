#include<bits/stdc++.h>
using namespace std;
#define ll long long
#define endl "\n" 
void solve(){
    int x,k;
    cin>>x>>k;
    if(x%k)cout<<1<<endl<<x<<endl;
    else{
        cout<<2<<endl;
        cout<<x+1<<" "<<-1<<endl;
    }
     return;
}
int main(void){
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);
    ll testcase=1;
    cin>>testcase;
    while(testcase--){
         solve();
    }
    
    return 0;
}