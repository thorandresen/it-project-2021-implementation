package com.example.verifyr;

import android.content.Context;
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

import java.io.UnsupportedEncodingException;
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
        String url = "https://ta.anderswiggers.dk/challenge/" + pufID;

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
    public void verify(int pufID, int challenge, int response, final VerifyVolleyCallback callback) {
        String url = "https://ta.anderswiggers.dk/verify";

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
    public void sendPk(byte[] pk, PkCallback callback) {
        String url = "https://ta.anderswiggers.dk/create-user";

        Map<String, Object> params = new HashMap<String, Object>();
        params.put("public_key", pk);
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
}


