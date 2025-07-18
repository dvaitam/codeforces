#include <bits/stdc++.h>

using namespace std;
//Mr abdoo 2018

int Hobi(int a, int b) {
    if (b == 0) return a;
    return Hobi(b, a%b);
}

int main()
{
    ios::sync_with_stdio(0);
    cin.tie(0);

    int t;
    cin>>t;

    while(t--){
        int Flag;
        cin>>Flag;
        int pgcd=Hobi(360,Flag*2);
        if (pgcd==1){cout<<"1\n";}
        else{
            int r=360/pgcd;
            if(Flag<90) cout<<r<<"\n";
            else{
                int m=min(2*Flag,360-Flag*2);
                if (m==360/r) r*=2;
                cout<<r<<"\n";
            }
        }
    }

    return 0;
}