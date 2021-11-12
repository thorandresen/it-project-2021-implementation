package com.example.verifyr;

public interface Server {
    void getChallenge(int pufID,final ChallengeVolleyCallback callback);
    void verify(int pufID,int challenge, int response,final VerifyVolleyCallback callback);
}
