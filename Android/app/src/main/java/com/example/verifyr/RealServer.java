package com.example.verifyr;

import android.content.Context;
import android.util.Base64;
import android.util.Log;

import com.android.volley.AuthFailureError;
import com.android.volley.NetworkResponse;
import com.android.volley.Request;
import com.android.volley.RequestQueue;
import com.android.volley.Response;
import com.android.volley.VolleyError;
import com.android.volley.VolleyLog;
import com.android.volley.toolbox.HttpHeaderParser;
import com.android.volley.toolbox.JsonObjectRequest;
import com.android.volley.toolbox.StringRequest;
import com.android.volley.toolbox.Volley;

import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.security.InvalidKeyException;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.Signature;
import java.security.SignatureException;
import java.security.UnrecoverableEntryException;
import java.security.cert.CertificateException;
import java.util.HashMap;
import java.util.Map;

public class RealServer implements Server{

    private final Context context;
    private final RequestQueue queue;

    public RealServer(Context context){
        this.context = context;
        queue = Volley.newRequestQueue(context);
    }
    @Override
    public void getChallenge(int pufID,final ChallengeVolleyCallback callback) {
        String url = "https://ta.anrs.dk/challenge/" + pufID;

        // Request a string response from the provided URL.
        StringRequest stringRequest = new StringRequest(Request.Method.GET, url,
                new Response.Listener<String>() {
                    @Override
                    public void onResponse(String response) {
                        callback.onSuccess(Integer.parseInt(response));
                    }
                }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                System.out.println("Error fetching");
            }
        });

        // Add the request to the RequestQueue.
        queue.add(stringRequest);
    }

    @Override
    public void verify(int pufID, int challenge, String response, final VerifyVolleyCallback callback) {
        String url = "https://ta.anrs.dk/verify";

        Map<String, Object> params = new HashMap<String, Object>();
        params.put("id", pufID);
        params.put("challenge", challenge);
        params.put("response", response);

        JSONObject jsonObject = new JSONObject(params);

        JsonObjectRequest postRequest = new JsonObjectRequest(Request.Method.POST, url, jsonObject,
                new Response.Listener<JSONObject>() {
                    @Override
                    public void onResponse(JSONObject response) {
                        Log.i("skrt",response.toString());
                        callback.onSuccess(true);
                    }
                }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                Log.i("skrt",error.toString());
            }
        });
        queue.add(postRequest);
    }

    @Override
    public void sendPk(PkCallback callback) {
        String url = "https://ta.anrs.dk/create-user";

        KeyStore.Entry entry = Util.getKeyEntry();

        byte[] pk = ((KeyStore.PrivateKeyEntry) entry).getCertificate().getPublicKey().getEncoded();
        String publicKeyString = Base64.encodeToString(pk, 2);

        Map<String, Object> params = new HashMap<String, Object>();
        params.put("public_key", publicKeyString);
        params.put("uuid", 1);
        params.put("mitID_token", 1);

        JSONObject jsonObject = new JSONObject(params);

        JsonObjectRequest postRequest = new JsonObjectRequest(Request.Method.POST, url, jsonObject,
                new Response.Listener<JSONObject>() {
                    @Override
                    public void onResponse(JSONObject response) {
                        Log.i("skrt",response.toString());
                        callback.onSuccess();
                    }
                }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                Log.i("skrt",error.toString());
            }
        });
        queue.add(postRequest);

    }



    @Override
    public void requestOwnership(int pufID, int buyerId,TransferCallback callback) {
        String url = "https://ta.anrs.dk/request";

        String sig = Util.sign(Integer.toString(pufID));

        Map<String, Object> params = new HashMap<String, Object>();
        params.put("sig", sig);
        params.put("bid", buyerId);

        JSONObject jsonObject = new JSONObject(params);

        JsonObjectRequest postRequest = new JsonObjectRequest(Request.Method.POST, url, jsonObject,
                new Response.Listener<JSONObject>() {
                    @Override
                    public void onResponse(JSONObject response) {
                        callback.onSuccess(response.toString());
                    }
                }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                Log.i("skrt",error.toString());
            }
        });
        queue.add(postRequest);
    }

    @Override
    public void transferOwnership(int pufID, int buyerId,TransferCallback callback) {
        String url = "https://ta.anrs.dk/transfer";

        String sig = Util.sign(Integer.toString(pufID) + Integer.toString(buyerId));
        Map<String, Object> params = new HashMap<String, Object>();
        params.put("sig", sig);
        params.put("bid", buyerId);

        JSONObject jsonObject = new JSONObject(params);

        JsonObjectRequest postRequest = new JsonObjectRequest(Request.Method.POST, url, jsonObject,
                new Response.Listener<JSONObject>() {
                    @Override
                    public void onResponse(JSONObject response) {
                        callback.onSuccess(response.toString());
                    }
                }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                Log.i("skrt",error.toString());
            }
        });
        queue.add(postRequest);
    }
}


