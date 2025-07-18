#include <iostream>

#include<vector>

#include<algorithm>

#include<math.h>

using namespace std;

typedef vector<int> vi;







int main()

{

int n,d;

cin >>n >>d;

if(n==1 && d==10 )

{

    cout <<-1;

    return 0;

}

if(d==10)

{



    n--;

    cout <<1;

    while(n--)

        cout<<0;

}

else{

    while(n--)

    {

        cout <<d;

    }

}





}