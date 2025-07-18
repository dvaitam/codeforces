#include <iostream>



using namespace std;

int f[100+10];

int main(){

    int n,a,b;

    cin>>n>>a>>b;

    for(int i=0;i<a;i++){

        int x;

        cin>>x;

        f[x]=1;

    }

    for(int i=0;i<b;i++){

        int x;

        cin>>x;

        f[x]=2;

    }

    for(int i=1;i<=n;i++)

        cout<<f[i]<<" ";

}