package com.example.verifyr;

import androidx.annotation.RequiresApi;
import androidx.appcompat.app.AppCompatActivity;

import android.os.Build;
import android.os.Bundle;
import android.security.keystore.KeyGenParameterSpec;
import android.security.keystore.KeyProperties;
import android.util.Log;
import android.view.View;
import android.widget.Button;

import java.io.IOException;
import java.security.InvalidAlgorithmParameterException;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.NoSuchProviderException;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.UnrecoverableEntryException;
import java.security.cert.CertificateException;

public class MainActivity extends AppCompatActivity {

    private Button verifyButton;
    private Puf puf;
    private Server server;
    private String CRP = "CRP";
    private Button keyGenButton;
    private Button getKeyButton;


    @RequiresApi(api = Build.VERSION_CODES.M)
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        server = new StubServer(this,1);
        puf = new StubPuf(1);

        verifyButton = findViewById(R.id.verify);
        verifyButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                verify();
            }
        });

        keyGenButton = findViewById(R.id.keyGen);
        keyGenButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                genKeyPair();
            }
        });

        getKeyButton = findViewById(R.id.printKey);
        getKeyButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                getKey();
            }
        });

    }

    @RequiresApi(api = Build.VERSION_CODES.M)
    private void genKeyPair() {
        KeyPairGenerator kpg = null;
        try {
            kpg = KeyPairGenerator.getInstance(
                    KeyProperties.KEY_ALGORITHM_RSA, "AndroidKeyStore");
        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
        } catch (NoSuchProviderException e) {
            e.printStackTrace();
        }

        try {
            kpg.initialize(new KeyGenParameterSpec.Builder(
                    "test",
                    KeyProperties.PURPOSE_SIGN | KeyProperties.PURPOSE_VERIFY)
                    .setDigests(KeyProperties.DIGEST_SHA256, KeyProperties.DIGEST_SHA512)
                    .setKeySize(2048)
                    .build());
        } catch (InvalidAlgorithmParameterException e) {
            e.printStackTrace();
        }

        KeyPair keyPair = kpg.generateKeyPair();
        // Log.i("Key",keyPair.toString());
    }

    private void getKey(){
        KeyStore keyStore = null;
        try {
            keyStore = KeyStore.getInstance("AndroidKeyStore");
        } catch (KeyStoreException e) {
            e.printStackTrace();
        }
        try {
            keyStore.load(null);
        } catch (CertificateException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
        }
        KeyStore.Entry entry = null;
        try {
            entry = keyStore.getEntry("test", null);
        } catch (KeyStoreException e) {
            e.printStackTrace();
        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
        } catch (UnrecoverableEntryException e) {
            e.printStackTrace();
        }
        PrivateKey privateKey = ((KeyStore.PrivateKeyEntry) entry).getPrivateKey();
        Log.i("Key",privateKey.toString());

        try {
            PublicKey publicKey = keyStore.getCertificate("test").getPublicKey();
            Log.i("Key",publicKey.toString());
        } catch (KeyStoreException e) {
            e.printStackTrace();
        }
    }

    private void verify() {
        int pufId = puf.getPufId();
        server.getChallenge(pufId, new ChallengeVolleyCallback() {
            @Override
            public void onSuccess(int challenge) {
                Log.i(CRP, "Challenge: " + challenge);
                int response = puf.doChallenge(challenge);
                Log.i(CRP, "Response: " + response);
                server.verify(pufId, challenge, response, new VerifyVolleyCallback() {
                    @Override
                    public void onSuccess(boolean verified) {
                        Log.i(CRP, "Verified: " + verified);
                    }
                });
            }
        });
        // int response = puf.doChallenge(challenge);
        // boolean verified = server.verify(pufId,challenge,response);
        // System.out.println(pufId+ " " + challenge+ " " + response + " " + verified);
    }


}