#include <bits/stdc++.h>
using namespace std;
int main() {
    int t;cin>>t;
    while(t--){
        int n;cin>>n;
        string a,b;cin>>a>>b;
        int ans=0;
        for(int i=0;i<n-2;i++) if(a[i]!=b[i] && a[i]==a[i+2] && b[i]==b[i+2]  && a[i+1]=='.' && b[i+1]=='.') ans++;
        cout<<ans<<endl;
    }
}