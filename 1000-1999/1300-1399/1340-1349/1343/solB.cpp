#include<iostream>

#include<bits/stdc++.h>

using namespace std;

int main()

{

    int t;

    cin >> t;

    while(t--)

    {

        int n;

        cin >> n;

        int d =  n / 2;

        if(d % 2 == 1)

        {

            cout << "NO\n";

        }

        else

        {

            cout << "YES\n";

            int sum = 0,sum1 = 0;

            for(int i = n ; i >= 2 ; i = i - 2)

            {

                cout << i << " ";

                sum = sum + i;

            }

            for(int j = n-3 ; j >= 1; j = j - 2)

            {

                cout << j << " ";

                sum1 = sum1 + j;

            }

            cout << sum - sum1 << endl;

        }

    }

}