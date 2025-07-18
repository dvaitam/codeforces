#include <bits/stdc++.h>

using namespace std;

int main(){
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);
    cout.tie(NULL);
    long long int t,a,b,k,l,r;
    cin>>t;
    while(t--){
        cin>>a>>b>>k;
        if(k%2==0){
            l = k/2;
            r = k/2;
        }
        else{
            l = k/2;
            r = l+1;
        }
        long long int sum = r*a - l*b;
        cout<<sum<<endl;
    }
}