#include<bits/stdc++.h>



using namespace std;



int n,m,mn,mx,type;

string s;



int main(){

    cin.tie(NULL);

    cout.tie(NULL);

    ios_base::sync_with_stdio(false);



    cin >> n;



    while(n % 7 != 0 && n - 4 >= 0){

        n -= 4;

        s = s + '4';

    }



    if (n % 7 != 0) return cout << -1,0;



    cout << s;

    for(int i = 1; i <= n / 7; ++i)

        cout << 7;



}