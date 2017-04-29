package me.oskareriksson.photofullstack;

import android.os.AsyncTask;
import android.util.Log;

import org.json.JSONObject;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.URL;

/**
 * Handles requests which are sent to the backend
 *
 * @author Oskar Eriksson
 * @version 1.0
 */
public class Request extends AsyncTask<Object, Void, Request.Response> {

    /**
     * Sends a HTTP request to the backend
     *
     * @param params params[0] is the URL to the API, params[1] is the JSON to be sent
     * @return Returns the HTTP response code, HTTP_NOT_FOUND if unable to connect/malformed URL
     *         and -1 if there was an internal error
     */
    @Override
    protected Response doInBackground(Object[] params) {
        // The parameters are the string of the URL and the JSON
        String urlString = (String) params[0];
        JSONObject jsonObject = (JSONObject) params[1];

        // Try to create a URL object
        URL url;
        try {
            url = new URL(urlString);
            Log.d(Models.FEEDBACK_SUCCESS, "URL object created " + url.toString());
        } catch (MalformedURLException e) {
            Log.d("ERROR", "Malformed URL: " + e.getMessage());
            return new Response(HttpURLConnection.HTTP_NOT_FOUND);
        }

        // Try to open a URL connection, send the request and handle the response
        HttpURLConnection connection;
        DataOutputStream writer;

        try {
            // Open a URL connection and do set up for sending the request
            connection = (HttpURLConnection) url.openConnection();
            connection.setRequestMethod("POST");
            connection.setRequestProperty("Content-Type", "application/json");
            connection.setRequestProperty("Accept", "application/json");
            connection.setRequestProperty("Token", Models.TOKEN);
            connection.setDoOutput(true);
            connection.setDoInput(true);
            connection.setConnectTimeout(2500);
            Log.d(Models.FEEDBACK_SUCCESS, "URL connection opened with " + url.toString());

            // Open the output stream
            writer = new DataOutputStream(connection.getOutputStream());
            String toWrite = "";

            // Get JSON string if there is one
            if (jsonObject != null) {
                toWrite = jsonObject.toString();
            }

            // Send the request
            writer.writeBytes(toWrite);
            writer.flush();
            writer.close();
            Log.d(Models.FEEDBACK_SUCCESS, "DataOutputStream sent request");

            // Get the response code
            int responseCode = connection.getResponseCode();
            if (responseCode == HttpURLConnection.HTTP_OK) {
                Log.d(Models.FEEDBACK_SUCCESS, "Got successful response code " + responseCode);
            } else {
                Log.d(Models.FEEDBACK_ERROR, "Got error response code " + responseCode);
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

            String responseString = response.toString();
            Log.d(Models.FEEDBACK_SUCCESS, "Got response: " + responseString);

            return new Response(responseCode, responseString);
        } catch (IOException e) {
            Log.d(Models.FEEDBACK_ERROR, "Caught IOException: " + e.getMessage());
            if (e.getMessage().contains("Connection refused")
                    || e.getMessage().contains("Unable to resolve host")
                    || e.getMessage().contains("failed to connect")
                    || e.getMessage().contains("connect timed out")) {
                // Return a not found response on any of the error messages
                return new Response(HttpURLConnection.HTTP_NOT_FOUND);
            } else {
                // Return an unauthorized response on other IOExceptions
                return new Response(HttpURLConnection.HTTP_UNAUTHORIZED);
            }
        }
    }

    /**
     * A class that make up a HTTP response from the backend
     */
    public static final class Response {
        private int responseCode;
        private String responseString;

        /**
         * Constructor for the Response class (if there was a response string)
         *
         * @param responseCode The response code returned
         * @param responseString The response string returned
         */
        public Response(int responseCode, String responseString) {
            this.responseCode = responseCode;
            this.responseString = responseString;
        }
        /**
         * Constructor for the Response class (no response string)
         *
         * @param responseCode The response code returned
         */
        public Response(int responseCode) {
            this.responseCode = responseCode;
            this.responseString = "";
        }

        /**
         * @return The response code
         */
        public int getResponseCode() {
            return responseCode;
        }

        /**
         * @return The response string
         */
        public String getResponseString() {
            return responseString;
        }
    }
}