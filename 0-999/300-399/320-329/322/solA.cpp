#include <iostream>
using namespace std;
int main(){
    long m,n,i;
    cin>>m>>n;
    cout<<m+n-1<<endl;
    for(i=1;i<=n;i++)
        cout<<1<<' '<<i<<endl;
    for(i=2;i<=m;i++)
        cout<<i<<' '<<1<<endl;
}