#include<bits/stdc++.h>

using namespace std;



int n,p,q;



int main(){



    string s;

    int fl = 0;



    cin >> n >> p >> q;

    cin >> s;



    int xx = -1,yy = -1;



    for(int x = 0; x * q <= n; x++){

        int y = ( n - (q * x) );

        if(y % p == 0){

            xx = x;

            yy = y / p;

            break;

        }



    }





    if(xx == -1 && yy == -1){

        cout << "-1" << endl;

        return 0;

    }



    int total = xx + yy;



    int i = 0;



    cout << total << endl;



    for(int k = 0; k < xx; k++){

        string t = "";

        for(int l = 0; l < q; l++,i++){

            t+=s[i];

        }

        cout << t << endl;

    }

    for(int k = 0; k < yy; k++){

        string t = "";

        for(int l = 0; l < p; l++,i++){

            t+=s[i];

        }

        cout << t << endl;

    }



}