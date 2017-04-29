package me.oskareriksson.photofullstack;

import android.content.Intent;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.EditText;
import android.widget.TextView;

import org.json.JSONException;
import org.json.JSONObject;

import java.net.HttpURLConnection;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ExecutionException;

/**
 * The main activity is the login screen
 *
 * @author Oskar Eriksson
 * @version 1.0
 */
public class MainActivity extends AppCompatActivity {
    private TextView status;

    /**
     * Called when the activity starts
     *
     * @param savedInstanceState The saved instance state
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        // Get the status TextView
        status = (TextView) findViewById(R.id.status_text);

        // The user is already logged in, start intent to go to PhotoListActivity
        if (!Models.TOKEN.equals("")) {
            startActivity(new Intent(this, PhotoListActivity.class));
        }
    }

    /**
     * Requests a token from the backend, if successful the user is redirected to their photos
     *
     * @param view The login button
     */
    public void login(View view) {
        // The user is already logged in, start intent to go to PhotoListActivity
        if (!Models.TOKEN.equals("")) {
            startActivity(new Intent(this, PhotoListActivity.class));
        }

        Log.d(Models.FEEDBACK_SUCCESS, "Called login method");
        Request.Response response = sendRequest("/login");

        // Return if the request failed
        if (response == null) {
            return;
        }

        // Handle response depending on the status code returned from the backend
        if (response.getResponseCode() == HttpURLConnection.HTTP_OK &&
                !response.getResponseString().equals("")) {
            // Login attempt successful, token received
            Log.d(Models.FEEDBACK_SUCCESS, "Received: " + response.getResponseString());
            try {
                // Set the token to the value of the token field in the JSON response
                Models.TOKEN = (String) new JSONObject(
                        response.getResponseString()).get("token");
                setStatus(R.string.status_login_successful);

                // If the token is set, the login was successful, go to PhotoListActivity
                if (!Models.TOKEN.equals("")) {
                    Log.d(Models.FEEDBACK_SUCCESS, "Token set in app: " + Models.TOKEN);
                    startActivity(new Intent(this, PhotoListActivity.class));
                }
            } catch (JSONException e) {
                Log.d(Models.FEEDBACK_ERROR, "JSONException: " + e);
                setStatus(R.string.status_internal_error);
            }
        } else if (response.getResponseCode() == HttpURLConnection.HTTP_BAD_REQUEST) {
            Log.d(Models.FEEDBACK_ERROR, "Malformed login credentials");
            setStatus(R.string.status_malformed_credentials);
        } else if (response.getResponseCode() == HttpURLConnection.HTTP_UNAUTHORIZED) {
            Log.d(Models.FEEDBACK_ERROR, "Invalid username/password combination on login");
            setStatus(R.string.status_malformed_credentials);
        } else if (response.getResponseCode() == HttpURLConnection.HTTP_NOT_FOUND) {
            Log.d(Models.FEEDBACK_ERROR, "Error connecting to server");
            setStatus(R.string.status_malformed_url);
        } else {
            Log.d(Models.FEEDBACK_ERROR, "Error with response code "
                    + response.getResponseCode());
            setStatus(R.string.status_internal_error);
        }
    }

    /**
     * Sends a registration request to the server, signs up a new user if everything went good
     *
     * @param view The registration button
     */
    public void register(View view) {
        Log.d(Models.FEEDBACK_SUCCESS, "Called register method");
        Request.Response response = sendRequest("/signup");

        // Return if the request failed
        if (response == null) {
            return;
        }

        // Handle response depending on the status code returned from the backend
        if (response.getResponseCode() == HttpURLConnection.HTTP_OK &&
                !response.getResponseString().equals("")) {
            Log.d(Models.FEEDBACK_ERROR, "Registration successful");
            setStatus(R.string.status_register_successful);
        } else if (response.getResponseCode() == HttpURLConnection.HTTP_BAD_REQUEST ||
                response.getResponseCode() == HttpURLConnection.HTTP_UNAUTHORIZED) {
            Log.d(Models.FEEDBACK_ERROR, "Username already taken/malformed credentials");
            setStatus(R.string.status_username_taken);
        } else if (response.getResponseCode() == HttpURLConnection.HTTP_NOT_FOUND) {
            Log.d(Models.FEEDBACK_ERROR, "Error connecting to server");
            setStatus(R.string.status_malformed_url);
        } else {
            Log.d(Models.FEEDBACK_ERROR, "Error with response code "
                    + response.getResponseCode());
            setStatus(R.string.status_internal_error);
        }
    }

    /**
     * Sends a login or a registration request to the backend, depending on the specified API
     *
     * @param api The API to send the request to
     * @return The message returned from the Request class
     */
    private Request.Response sendRequest(String api) {
        // Get input and make sure that no input is empty
        Map<String, String> input = getFields();
        if (Models.validateInput(input)) {
            // If input was valid, set constant strings to the desired values
            Models.IP = input.get("ip");
            Models.PORT = input.get("port");
            Models.USER = input.get("username");
        } else {
            // If the input was incomplete, display a status message and return null
            Log.d(Models.FEEDBACK_ERROR, "Input fields incomplete");
            setStatus(R.string.status_incomplete_input);
            return null;
        }

        // Create a JSON object with the login credentials, return if it fails
        JSONObject credentials = createJSON(input.get("username"), input.get("password"));
        if (credentials == null) {
            setStatus(R.string.status_malformed_credentials);
            return null;
        }

        // Send request to the backend
        try {
            Request.Response response = new Request().execute(
                    "http://" + Models.IP + ":" + Models.PORT + api,
                    credentials).get();
            return response;
        } catch (ExecutionException e) {
            Log.d(Models.FEEDBACK_ERROR, "Login ExecutionException " + e.getMessage());
            setStatus(R.string.status_internal_error);
            return null;
        } catch (InterruptedException e) {
            Log.d(Models.FEEDBACK_ERROR, "Login InterruptedException " + e.getMessage());
            setStatus(R.string.status_internal_error);
            return null;
        }
    }

    /**
     * Gets the IP, port, username and password fields on the login screen
     *
     * @return A map of strings with the values from the input fields
     */
    private Map<String, String> getFields() {
        // Get input from all the input fields
        Map<String, String> input = new HashMap<>();
        input.put("ip", ((EditText) findViewById(R.id.backend_ip_input)).getText().toString());
        input.put("port", ((EditText) findViewById(R.id.backend_port_input)).getText().toString());
        input.put("username", ((EditText) findViewById(R.id.username_input)).getText().toString());
        input.put("password", ((EditText) findViewById(R.id.password_input)).getText().toString());
        Log.d(Models.FEEDBACK_SUCCESS, "Got text from input fields");
        return input;
    }

    /**
     * Generate a JSON object based on a username/password combination
     *
     * @param username The username as a string
     * @param password The password as a string
     * @return The generated JSON object
     */
    private JSONObject createJSON(String username, String password) {
        // Create a JSON object with the username and password
        JSONObject json = null;
        try {
            json = new JSONObject();
            json.put("username", username);
            json.put("password", password);
            Log.d(Models.FEEDBACK_SUCCESS, "Credentials JSON created " + json.toString());
        } catch (JSONException e) {
            Log.d("ERROR", e.getMessage());
        }

        return json;
    }

    /**
     * Sets the text of the status TextView to a given string
     *
     * @param stringCode The code of the string
     */
    public void setStatus(int stringCode) {
        status.setText(getString(stringCode));
    }
}
