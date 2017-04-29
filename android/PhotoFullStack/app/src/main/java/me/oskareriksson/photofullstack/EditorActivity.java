package me.oskareriksson.photofullstack;

import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Base64;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.TextView;

import java.net.HttpURLConnection;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ExecutionException;

/**
 * The EditorActivity lets the user see their photo and set it's title and description
 *
 * @author Oskar Eriksson
 * @version 1.0
 */
public class EditorActivity extends AppCompatActivity {
    TextView status;
    Photo photo;

    /**
     * Called when the activity starts
     *
     * @param savedInstanceState The saved instance state
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_editor);

        status = (TextView) findViewById(R.id.editor_status);

        // Engage photo editing
        photo = Models.PHOTO_EDIT;
        if (photo == null) {
            // If there is no photo to be edited, finish this activity
            this.finish();
        } else {
            // If there is a photo, display it in the ImageView
            displayBase64(photo.getJpgBase64());
        }
    }

    /**
     * Add action buttons to the action bar
     *
     * @param menu The menu object
     * @return Always return true
     */
    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        getMenuInflater().inflate(R.menu.menu_editor, menu);
        return true;
    }

    /**
     * Handle clicks on the icons in the action bar
     *
     * @param item The item clicked
     * @return True if recognized, otherwise let superclass handle it
     */
    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        switch (item.getItemId()) {
            case R.id.upload_action_button:
                Log.d(Models.FEEDBACK_SUCCESS, "Upload action button pressed");
                if (uploadEditedPhoto() || photo == null) {
                    // Upload successful, go to PhotoListActivity
                    startActivity(new Intent(this, PhotoListActivity.class));
                }
                return true;
            default:
                return super.onOptionsItemSelected(item);
        }
    }

    /**
     * Display the image in the ImageView
     *
     * @param jpgBase64 The Base64 string of the image to be displayed
     */
    private void displayBase64(String jpgBase64) {
        // Create bitmap from Base64
        byte[] imageAsBytes = Base64.decode(jpgBase64.getBytes(), 0);
        Bitmap bitmap = BitmapFactory.decodeByteArray(imageAsBytes, 0, imageAsBytes.length);

        // Set the ImageView to the bitmap
        ((ImageView) findViewById(R.id.image_display)).setImageBitmap(bitmap);
    }

    /**
     * Get title and description input field values
     *
     * @return A map of strings with the values from the input fields
     */
    private Map<String, String> getFields() {
        // Get input from all the input fields
        Map<String, String> input = new HashMap<>();
        input.put("title", ((EditText) findViewById(R.id.title_input)).getText().toString());
        input.put("description", ((EditText) findViewById(
                R.id.description_input)).getText().toString());
        Log.d(Models.FEEDBACK_SUCCESS, "Got text from input fields in editor");
        return input;
    }

    private boolean uploadEditedPhoto() {
        // If there is no photo to be uploaded, return
        if (photo == null) {
            return false;
        }

        // Check if there if all fields are set
        Map<String, String> input = getFields();
        if (!Models.validateInput(input)) {
            setStatus(R.string.status_incomplete_input);
            return false;
        }

        // Set title and description of the photo
        photo.setTitle(input.get("title"));
        photo.setDescription(input.get("description"));

        // Try to send the request to upload the photo
        Request.Response response;
        try {
            response = new Request().execute(
                    "http://" + Models.IP + ":" + Models.PORT + "/upload", photo.toJSON()).get();
            if (response.getResponseCode() == HttpURLConnection.HTTP_OK) {
                Log.d(Models.FEEDBACK_SUCCESS, "Successfully uploaded photo");
                setStatus(R.string.upload_success);
                Models.PHOTO_EDIT = null;
                photo = null;
                return true;
            } else {
                // Invalid response
                Log.d(Models.FEEDBACK_ERROR, "HTTP response not OK");
                setStatus(R.string.upload_failure);
            }
        } catch (InterruptedException e) {
            Log.d(Models.FEEDBACK_ERROR, "InterruptedException EditorActivity: " + e.getMessage());
        } catch (ExecutionException e) {
            Log.d(Models.FEEDBACK_ERROR, "ExecutionException EditorActivity: " + e.getMessage());
        }

        // If the method didn't return before this, the upload failed
        return false;
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
