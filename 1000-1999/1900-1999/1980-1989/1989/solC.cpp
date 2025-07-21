#include <bits/stdc++.h>
using namespace std;
int main(){
    long long c,t,n,e,d;
    cin>>t;
    while(t--){
        cin>>n;
        long long a[n],b[n];
        e=0,c=0,d=0;
        for(int i=0;i<n;i++){
            cin>>a[i];
        }
        for(int i=0;i<n;i++){
            cin>>b[i];
            if(a[i]>b[i])
                c+=a[i];
            else if(a[i]<b[i])
                d+=b[i];
            else if(a[i]==1)
                e+=a[i];
            else if(a[i]==-1)
                e++,c--,d--;
        }
        cout<<min(c+e,min(d+e,(c+e+d)>>1))<<endl;
    }
    return 0;
}
