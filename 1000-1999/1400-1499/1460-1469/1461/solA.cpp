#include <bits/stdc++.h>

#define lint unsigned long long int

#define dint long double

using namespace std;









int main()

{

    int t, n, k, flag;



    cin >> t;



    for (int i=0; i<t; i++) {

        cin >> n >> k;

        flag = 0;



        for (int j=0; j<k; j++) cout << 'a';

        for (int j=k; j<n; j++) {

            if (flag == 0) {

                flag = 1;

                cout << 'b';

            } else if (flag == 1) {

                flag = 2;

                cout << 'c';

            } else {

                flag = 0;

                cout << 'a';

            }

        }



        cout << endl;

    }

    







    return 0;

}