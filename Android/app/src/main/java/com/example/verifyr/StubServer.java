package com.example.verifyr;

import android.content.Context;
import android.util.Log;

import java.security.InvalidKeyException;
import java.security.KeyStore;
import java.security.NoSuchAlgorithmException;
import java.security.Signature;

public class StubServer implements Server{


    public StubServer(Context context){

    }
    @Override
    public void getChallenge(int pufID,final ChallengeVolleyCallback callback) {
        callback.onSuccess(5);
    }

    @Override
    public void verify(int pufID, int challenge, String response, VerifyVolleyCallback callback) {
        callback.onSuccess(response.equals(Util.hash(Integer.toString(pufID) + Integer.toString(challenge))));
    }

    @Override
    public void sendPk(PkCallback callback) {
    }

    @Override
    public void requestOwnership(int pufID, int buyerId,TransferCallback callback) {
        callback.onSuccess("Transfer Requested");
    }

    @Override
    public void transferOwnership(int pufID, int buyerId,TransferCallback callback) {
        callback.onSuccess("Ownership transferred");
    }

}
