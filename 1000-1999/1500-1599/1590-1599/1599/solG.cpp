#include <bits/stdc++.h>



using namespace std;



struct pont{

    int x, y;

    int ind;



    void kiir(){

        cout<<x<<' '<<y<<'\n';

    }

};



double tav(pont a, pont b){

    return sqrt(pow(a.x-b.x, 2.0) + pow(a.y-b.y, 2.0));

}



bool jo(pont a, pont b, pont c){

    long long x1 = b.x-a.x, x2 = c.x-a.x, y1 = b.y-a.y, y2 = c.y-a.y;

    return x1*y2 == x2*y1;

}



int main(){

    ios_base::sync_with_stdio(false);

    cin.tie(0);



    int n, k;

    cin>>n>>k;

    --k;



    vector<pont> v(n);

    for(int i = 0; i < n; i++){

        cin>>v[i].x>>v[i].y;

        v[i].ind = i;

    }



    if(n > 3){

        for(int i = 0; i < 4; i++){

            if(jo(v[0], v[1], v[2]))

                break;

            rotate(v.begin(), v.begin()+1, v.begin()+4);

        }

    }



    for(int i = 3; i < n; i++){

        if(!jo(v[0], v[1], v[i])){

            swap(v[0], v[i]);

            break;

        }

    }



    sort(v.begin()+1, v.end(), [](pont a, pont b) -> bool{ return (a.x == b.x ? a.y < b.y : a.x < b.x); });



    int j = 0;

    while(v[j].ind != k) ++j;



    if(j == 0){

        cout<<fixed<<setprecision(10)<<min(tav(v[0], v[1]) + tav(v[1], v.back()), tav(v[0], v.back()) + tav(v[1], v.back()))<<'\n';

        return 0;

    }



    double meg = 1e18;



    meg = min(meg, tav(v[1], v[0]) + tav(v[1], v[j]) + min(j == n-1 ? 0.0 : tav(v[0], v[j+1]), tav(v[0], v.back())) + (j == n-1 ? 0.0 : tav(v[j+1], v.back())));

    meg = min(meg, tav(v[n-1], v[0]) + tav(v[j], v[n-1]) + min(j == 1 ? 0.0 : tav(v[0], v[j-1]), tav(v[0], v[1])) + (j == 1 ? 0.0 : tav(v[j-1], v[1])));



    for(int i = 1; i <= j; i++)

        meg = min(meg, tav(v[j], v.back())*2.0 + tav(v[i], v[j]) + tav(v[i], v[0]) + (i == 1 ? 0.0 : min(tav(v[i-1], v[0]), tav(v[1], v[0])) + tav(v[1], v[i-1])));

    for(int i = j; i < n; i++)

        meg = min(meg, tav(v[1], v[j])*2.0 + tav(v[j], v[i]) + tav(v[i], v[0]) + (i == n-1 ? 0.0 : min(tav(v[i+1], v[0]), tav(v.back(), v[0])) + tav(v[i+1], v.back())));



    cout<<fixed<<setprecision(10)<<meg<<'\n';



    return 0;

}