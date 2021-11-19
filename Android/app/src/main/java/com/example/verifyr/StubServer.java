package com.example.verifyr;

import android.content.Context;

import java.util.HashMap;

public class StubServer implements Server{

    private final int pufId;

    public StubServer(Context context, int pufId){
        this.pufId = pufId;
    }
    @Override
    public void getChallenge(int pufID,final ChallengeVolleyCallback callback) {
        callback.onSuccess(pufID*4);
    }

    @Override
    public void verify(int pufID, int challenge, int response, VerifyVolleyCallback callback) {
//        if(response==challenge/2){
//            callback.onSuccess(true);
//        }else{
//            callback.onSuccess(false);
//        }
    }

    @Override
    public void sendPk(String pk, PkCallback callback) {

    }

}
