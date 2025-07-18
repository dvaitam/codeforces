#include <bits/stdc++.h>

using namespace std;

#define ll long long

int main()
{
    ll a1,b1,a2,b2;

    string s;

    cin >> s;

    a1 = (s[0]-'0')*10+(s[1]-'0');
    b1 = (s[3]-'0')*10+(s[4]-'0');

    cin >> s;

    a2 = (s[0]-'0')*10+(s[1]-'0');
    b2 = (s[3]-'0')*10+(s[4]-'0');


    ll a3,b3,aux;

    aux = (a2-a1)*60 + (b2-b1)%60;

    a3 = a1;

    b3 = b1;

    while(aux > 59){
        aux-=60;
        b3+=30;
        if(b3>59){a3++;b3-=60;}
    }
    b3 += aux/2;

    if(b3 > 59){b3 -=60; a3++;}

    if(a3 < 10){
        cout << 0;
    }

    cout << a3 << ":";

    if(b3 < 10){cout<<0;}
    cout << b3;



    return 0;
}