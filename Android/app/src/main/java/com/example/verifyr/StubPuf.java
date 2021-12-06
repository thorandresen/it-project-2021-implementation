package com.example.verifyr;

import java.util.HashMap;

public class StubPuf implements Puf{
    private final int id;
    private final HashMap<Integer, Integer> puf;

    public StubPuf(int id){
        this.id = id;
        puf = new HashMap<Integer, Integer>() {{
            put(id*4,id*2);
        }};
    }


    @Override
    public String doChallenge(int challenge) {
        return Util.hash(Integer.toString(id) + Integer.toString(challenge) );
    }

    @Override
    public int getPufId() {
        return id;
    }


}
