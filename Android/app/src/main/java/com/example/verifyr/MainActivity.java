package com.example.verifyr;

import androidx.annotation.RequiresApi;
import androidx.appcompat.app.AppCompatActivity;

import android.os.Build;
import android.os.Bundle;
import android.security.keystore.KeyGenParameterSpec;
import android.security.keystore.KeyProperties;
import android.util.Base64;
import android.util.Log;
import android.widget.Button;

import java.io.IOException;
import java.security.InvalidAlgorithmParameterException;
import java.security.InvalidKeyException;
import java.security.KeyPairGenerator;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.NoSuchProviderException;
import java.security.Signature;
import java.security.SignatureException;
import java.security.UnrecoverableEntryException;
import java.security.cert.CertificateException;
import java.util.Arrays;

public class MainActivity extends AppCompatActivity {
    private Puf puf;
    private Server server;
    private final String CRP = "CRP";


    @RequiresApi(api = Build.VERSION_CODES.M)
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        server = new RealServer(this);
        puf = new StubPuf(2);

        Button verifyButton = findViewById(R.id.verify);
        verifyButton.setOnClickListener(view -> verify());

        Button keyGenButton = findViewById(R.id.keyGen);
        keyGenButton.setOnClickListener(view -> genKeyPair());

        Button sendPkButton = findViewById(R.id.sendPkButton);
        sendPkButton.setOnClickListener(view -> sendPk());

        Button requestButton = findViewById(R.id.requestButton);
        requestButton.setOnClickListener(view -> server.requestOwnership(puf.getPufId(),1,(response) -> Log.i("req",response)));

        Button transferButton = findViewById(R.id.transferButton);
        transferButton.setOnClickListener(view -> server.transferOwnership(puf.getPufId(),1,(response) -> Log.i("transfer",response)));

    }

    private void sendPk() {
        server.sendPk(() -> Log.i("pkResp","Done"));
    }

    @RequiresApi(api = Build.VERSION_CODES.M)
    private void genKeyPair() {
        KeyPairGenerator kpg = null;
        try {
            kpg = KeyPairGenerator.getInstance(
                    KeyProperties.KEY_ALGORITHM_RSA, "AndroidKeyStore");
        } catch (NoSuchAlgorithmException | NoSuchProviderException e) {
            e.printStackTrace();
        }

        try {
            assert kpg != null;
            kpg.initialize(new KeyGenParameterSpec.Builder(
                    "test",
                    KeyProperties.PURPOSE_SIGN | KeyProperties.PURPOSE_VERIFY)
                    .setDigests(KeyProperties.DIGEST_SHA256, KeyProperties.DIGEST_SHA512)
                    .setKeySize(2048)
                    .setSignaturePaddings(KeyProperties.SIGNATURE_PADDING_RSA_PSS)
                    .build());
        } catch (InvalidAlgorithmParameterException e) {
            e.printStackTrace();
        }

        kpg.generateKeyPair();
        // Log.i("Key",keyPair.toString());
    }


    private void verify() {
        int pufId = puf.getPufId();
        server.getChallenge(pufId, challenge -> {
            Log.i(CRP, "Challenge: " + challenge);
            String response = puf.doChallenge(challenge);
            Log.i(CRP, "Response: " + response);
            server.verify(pufId, challenge, response, verified -> Log.i(CRP, "Verified: " + verified));
        });
    }


}