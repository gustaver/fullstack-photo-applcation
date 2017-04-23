package me.oskareriksson.photofullstack;

import android.os.AsyncTask;
import android.util.Log;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.URL;

/**
 * Handles requests
 *
 * @author Oskar Eriksson
 * @version 1.0
 */
public class Request extends AsyncTask<String, Void, String> {
    public static final String[] errors = {
            "No error",
            "URL error",
            "Input error",
    };

    @Override
    protected String doInBackground(String[] params) {
        // Try to create a URL object
        URL url;
        try {
            url = new URL(params[0]);
            Log.d(Models.FEEDBACK_SUCCESS, "URL object created " + url.toString());
        } catch (MalformedURLException e) {
            Log.d("ERROR", e.getMessage());
            return errors[1];
        }

        // Try to open a URL connection, send the request and handle the response
        HttpURLConnection connection;
        DataOutputStream writer;
        try {
            // Open a URL connection
            connection = (HttpURLConnection) url.openConnection();
            connection.setRequestMethod("POST");
            connection.setRequestProperty("Content-Type", "application/json");
            connection.setRequestProperty("Accept", "application/json");

            connection.setDoOutput(true);
            connection.setDoInput(true);
            Log.d(Models.FEEDBACK_SUCCESS, "URL connection opened with " + url.toString());

            // Send the request
            writer = new DataOutputStream(connection.getOutputStream());

            writer.writeBytes(params[1]);
            writer.flush();
            writer.close();
            Log.d(Models.FEEDBACK_SUCCESS, "OutputStreamWriter sent " + params[1]);

            // Get the response code
            int responseCode = connection.getResponseCode();
            if (responseCode == HttpURLConnection.HTTP_OK) {
                Log.d(Models.FEEDBACK_SUCCESS, "Got successful response code " + responseCode);
            } else {
                Log.d(Models.FEEDBACK_ERROR, "Got error response code " + responseCode);
                if (responseCode == HttpURLConnection.HTTP_UNAUTHORIZED) {
                    return errors[2];
                }
            }

            // Read the response
            BufferedReader reader = new BufferedReader(
                    new InputStreamReader(connection.getInputStream()));
            String line;
            StringBuffer response = new StringBuffer();
            while ((line = reader.readLine()) != null) {
                response.append(line);
            }
            reader.close();

            Log.d(Models.FEEDBACK_SUCCESS, "Got response: " + response.toString());
        } catch (IOException e) {
            Log.d(Models.FEEDBACK_ERROR, "Caught IOException: " + e.getMessage());
            return errors[1];
        }

        // If all went well, return no error
        return errors[0];
    }
}