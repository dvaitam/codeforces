#include <bits/stdc++.h>



using namespace std;



#define int long long



int n;

vector <int> vec;

bool used[1005];



bool prime(int x)

{

    if (x==1)

        return 0;

    if (x==2)

        return 1;

    if (x%2==0)

        return 0;

    for (int i=3;i<=sqrt(x);i+=2)

        if (x%i==0)

            return 0;

    return 1;

}



main()

{

    cin >> n;

    for (int i=1;i<=n;i++)

    {

        if (prime(i))

        {

            int p=i;

            while (p<=n)

            {

                vec.push_back(p);

                p*=i;

            }

        }

    }

    cout << vec.size() << "\n";

    for (int i=0;i<vec.size();i++)

    {

        cout << vec[i] << " ";

    }

}