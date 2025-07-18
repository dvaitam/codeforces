#include<bits/stdc++.h>

using namespace std;

using i64=long long;

void solve(){

    int n;

    cin>>n;

    auto isprime=[](int x){

        if(x<2){

            return false;

        }

        for(int i=2;i<=sqrt(x);i++){

            if(x%i==0){

                return false;

            }

        }

        return true;

    };

    int a,b;

    for(int i=1;;i++){

        if(isprime(i+n-1)&&!isprime(i)){

            a=i;

            break;

        }

    }

    for(int i=1;;i++){

        if(isprime((n-1)*a+i)&&!isprime(i)){

            b=i;

            break;

        }

    }

    for(int i=1;i<=n-1;i++){

        for(int j=1;j<=n-1;j++){

            cout<<1<<' ';

        }

        cout<<a<<endl;

    }

    for(int i=1;i<=n-1;i++){

        cout<<a<<' ';

    }

    cout<<b<<endl;

}

int main(){

    ios::sync_with_stdio(false);

    cin.tie(nullptr);

    int T=1;

    cin>>T;

    while(T--) solve();

    return 0;

}