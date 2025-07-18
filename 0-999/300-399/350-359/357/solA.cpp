#include <iostream>



using namespace std;



int main()

{

    int number,num1,num2;

    long long sum1=0,sum2=0;

    cin>>number;

    int arr[number];

    for (int i=0;i<number;i++)

    {

        cin>>arr[i];

        sum1+=arr[i];

    }

    cin>>num1>>num2;

    for (int i=0;i<number;i++)

    {

        sum2+=arr[i];

        sum1-=arr[i];

        if (sum2>=num1 && sum2<=num2 && sum1>=num1 && sum1<=num2)

        {

            cout<<i+2;

            return 0;

        }

    }

    cout<<0;

    return 0;

}