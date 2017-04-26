package me.oskareriksson.photofullstack;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.EditText;
import android.widget.TextView;

import org.json.JSONException;
import org.json.JSONObject;

import java.util.concurrent.ExecutionException;

/**
 * The main activity is the login screen
 *
 * @author Oskar Eriksson
 * @version 1.0
 */
public class MainActivity extends AppCompatActivity {

    /**
     * Called when the activity starts
     *
     * @param savedInstanceState The saved instance state
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
    }

    /**
     * Requests a token from the backend, if successful the user is redirected to their photos
     *
     * @param view The login button
     */
    public void login(View view) {
        Log.d(Models.FEEDBACK_SUCCESS, "Called login method");

        // Get the status TextView and clear it
        TextView status = (TextView) findViewById(R.id.status_text);
        status.setText(null);

        // Get input from all the input fields
        String ip = ((EditText) findViewById(R.id.backend_ip_input)).getText().toString();
        String port = ((EditText) findViewById(R.id.backend_port_input)).getText().toString();
        String username = ((EditText) findViewById(R.id.username_input)).getText().toString();
        String password = ((EditText) findViewById(R.id.password_input)).getText().toString();
        Log.d(Models.FEEDBACK_SUCCESS, "Got text from input fields");

        // Make sure that no input is empty
        if (ip.length() == 0 || port.length() == 0 ||
                username.length() == 0 || password.length() == 0) {
            status.setText(getString(R.string.status_incomplete_input));
            return;
        }

        // Create a JSON object with the login credentials
        JSONObject credentials;
        try {
            credentials = new JSONObject();
            credentials.put("username", username);
            credentials.put("password", password);
            Log.d(Models.FEEDBACK_SUCCESS, "Credentials JSON created " + credentials.toString());
        } catch (JSONException e) {
            Log.d("ERROR", e.getMessage());
            status.setText(getString(R.string.status_malformed_credentials));
            return;
        }

        // Send login request to the backend
        String result;
        try {
            result = new Request().execute(
                    "http://" + ip + ":" + port + "/login", credentials.toString()).get();
        } catch (ExecutionException e) {
            Log.d(Models.FEEDBACK_ERROR, "Login ExecutionException " + e.getMessage());
            return;
        } catch (InterruptedException e) {
            Log.d(Models.FEEDBACK_ERROR, "Login InterruptedException " + e.getMessage());
            return;
        }

        // Display a message based on the status of the request
        if (result.equals(Request.errors[0])) {
            // No error, display login successful
            status.setText(getString(R.string.status_login_successful));
        } else if (result.equals(Request.errors[1])) {
            // Unable to contact server
            status.setText(getString(R.string.status_malformed_url));
        } else if (result.equals(Request.errors[2])) {
            // Invalid login credentials
            status.setText(getString(R.string.status_malformed_credentials));
        }
    }
}
