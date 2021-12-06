package com.example.verifyr;

public interface Server {
    void getChallenge(int pufID,final ChallengeVolleyCallback callback);
    void verify(int pufID,int challenge, String response,final VerifyVolleyCallback callback);
    void sendPk(final PkCallback callback);
    void requestOwnership(int pufID,int buyerId,TransferCallback callback);
    void transferOwnership(int pufID,int buyerId,TransferCallback callback);
}
