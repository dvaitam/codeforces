#include<iostream>

#include<iomanip>

#define rep(i,a,b) for(int i = a;i < b;i++)



using namespace std;



/*double bus(double a)

{

    return (((double)l) - v*a)/((double)(b-v));

}



double returnb(double x)

{

    return ((x*(double)(b - v)))/((double)(v + b));

}

*/



int main()

{

    ios_base::sync_with_stdio(false);



    double l,v,b;

    long long n,k;



    cin >> n >> l >> v >> b >> k;



    int g = n / k;



    if(n % k > 0)

    {

        g++;

    }



  /*  while(high - low >= 0.0000000001)

    {

        double a = (high + low) / 2;

        double t = 0;

        double x = bus(a);

        double y = returnb(x);

        t = x*groups + y*(groups - 1);      answer = g*(bus(answer)) + returnb(bus(answer))*(g - 1)



        if(t > a)

        {

            low = a;

        }

        else

        {

            high = a;

        }

    }*/



    double a = (l*v + (2*b*g - b)*l)/((2*b*g-b)*v+b*b);



    cout << fixed << setprecision(10) << a;

}