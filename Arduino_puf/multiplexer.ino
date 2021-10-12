int set1 = 6;
int set2 = 7;
int set3 = 8;
int set4 = 9;
int reset = 13;
int output = 10;
int zero = 0;
int one = 0;
int noOfRuns = 1000;
void setup()
{
    Serial.begin(9600);
    pinMode(set1, INPUT);
    pinMode(set2, INPUT);
    pinMode(set3, INPUT);
    pinMode(set4, INPUT);
    pinMode(output, INPUT);
    pinMode(reset, INPUT);
}

void loop()
{
    int resultArray[16];
    for (int i = 0; i < 16; i++)
    {
        boolean bits[] = {0, 0, 0, 0};
        for (int j = 4; j >= 0; j--)
        {
            bits[j] = (i & (1 << j)) != 0;
        }

        int result = run(bits[0], bits[1], bits[2], bits[3]);
        resultArray[i] = result;
    }
    for (int i = 0; i < 16; i++)
    {
        Serial.print(resultArray[i]);
    }
    delay(1000);
    Serial.println("___________________________");
}

int run(int s1, int s2, int s3, int s4)
{
    zero = 0;
    one = 0;
    digitalWrite(set1, s1);
    digitalWrite(set2, s2);
    digitalWrite(set3, s3);
    digitalWrite(set4, s4);
    int val;
    for (int i = 0; i < noOfRuns; i++)
    {
        digitalWrite(reset, HIGH);
        delay(1);
        val = digitalRead(output);
        if (val == 0)
        {
            zero++;
        }
        else
        {
            one++;
        }
        digitalWrite(reset, LOW);
        delay(1);
    }
    // Serial.println("Zero: ");
    // Serial.println(zero);
    // Serial.println("One: ");
    // Serial.println(one);
    if (zero >= one)
    {
        return 0;
    }
    else
    {
        return 1;
    }
}