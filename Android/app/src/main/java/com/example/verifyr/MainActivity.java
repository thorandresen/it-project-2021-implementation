package com.example.verifyr;

import androidx.annotation.RequiresApi;
import androidx.appcompat.app.AppCompatActivity;

import android.os.Build;
import android.os.Bundle;
import android.security.keystore.KeyGenParameterSpec;
import android.security.keystore.KeyProperties;
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
        puf = new StubPuf(1);

        Button verifyButton = findViewById(R.id.verify);
        verifyButton.setOnClickListener(view -> verify());

        Button keyGenButton = findViewById(R.id.keyGen);
        keyGenButton.setOnClickListener(view -> genKeyPair());

        Button svButton = findViewById(R.id.signVerifyButton);
        svButton.setOnClickListener(view -> signAndVerify());

        Button sendPkButton = findViewById(R.id.sendPkButton);
        sendPkButton.setOnClickListener(view -> sendPk());

    }

    private void sendPk() {
        KeyStore.Entry entry = getKeyEntry();
        byte[] pk = ((KeyStore.PrivateKeyEntry) entry).getCertificate().getPublicKey().getEncoded();
        server.sendPk(pk, () -> Log.i("pkResp","Done"));
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
    public KeyStore.Entry getKeyEntry(){
        KeyStore ks;
        try {
            ks = KeyStore.getInstance("AndroidKeyStore");
            ks.load(null);
            KeyStore.Entry entry = ks.getEntry("test", null);
            if (!(entry instanceof KeyStore.PrivateKeyEntry)) {
                Log.w("pkfail", "Not an instance of a PrivateKeyEntry");
            }
            return entry;
        } catch (KeyStoreException | CertificateException | IOException | NoSuchAlgorithmException | UnrecoverableEntryException e) {
            e.printStackTrace();
        }
        return null;
    }

    public void signAndVerify(){
        try {
            KeyStore.Entry entry = getKeyEntry();
            // Real MSG
            Signature s = Signature.getInstance("SHA256withRSA/PSS");
            String msg = "skrt";
            byte[] data = msg.getBytes();
            Log.i("pkfail",((KeyStore.PrivateKeyEntry) entry).getPrivateKey().toString());
            s.initSign(((KeyStore.PrivateKeyEntry) entry).getPrivateKey());
            s.update(data);
            byte[] signature = s.sign();

            // Verify
            Signature sV = Signature.getInstance("SHA256withRSA/PSS");
            sV.initVerify(((KeyStore.PrivateKeyEntry) entry).getCertificate());
            sV.update(data);
            boolean valid = sV.verify(signature);
            Log.i("valid", String.valueOf(valid));


        } catch (NoSuchAlgorithmException | InvalidKeyException | SignatureException e) {
            e.printStackTrace();
        }
    }

    private void verify() {
        int pufId = puf.getPufId();
        server.getChallenge(pufId, challenge -> {
            Log.i(CRP, "Challenge: " + challenge);
            int response = puf.doChallenge(challenge);
            Log.i(CRP, "Response: " + response);
            server.verify(pufId, challenge, response, verified -> Log.i(CRP, "Verified: " + verified));
        });
    }


}