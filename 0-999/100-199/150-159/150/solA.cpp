#include <bits/stdc++.h>
#include <bits/stdc++.h>

using namespace std;



typedef long long ll;



int cnt;

ll a;

vector<ll>vec;

bool pr;



int check(ll n)

{

    int cnt = 0;

    while (n%2 == 0){

        cnt++;

        //cout << 2 << endl;

        vec.push_back(2LL);

        if(cnt > 2) return 1;

        n = n/2;

    }



    for (ll i = 3; i * i <= n; i = i+2){

        while (n%i == 0)

        {

            cnt++;

            //cout << i << endl;

            vec.push_back(i);

            if(cnt > 2) return 1;

            n = n/i;

        }

    }



    if(n > 2){

        cnt++;

        vec.push_back(n);

        //cout << n << endl;

    }

    if(cnt == 2){

        return 0;

    }

    else if(cnt == 1) {

        pr = 1;

        return 1;

    }

    else return 1;



}





int main()

{

    scanf("%lld",&a);



    if(a == 1){

        printf("1\n0\n");

        return 0;

    }



    int ch = check(a);



    if(ch){

        if(pr){

            printf("1\n0\n");

        }

        else{

            printf("1\n");

            ll ans = 1;

            for(int i = 0; i < 2; i++) ans = ans * vec[i];

            printf("%lld\n",ans);

        }

    }

    else printf("2\n");



    //sort(v.begin(),v.end());

    //cout << ch << endl;







    //for(int i = 0; i < sz; i++) cout << v[i] << endl;



    return 0;

}