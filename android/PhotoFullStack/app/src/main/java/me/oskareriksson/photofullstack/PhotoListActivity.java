package me.oskareriksson.photofullstack;

import android.content.Intent;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;

/**
 * The PhotoListActivity is where the photos are displayed
 *
 * @author Oskar Eriksson
 * @version 1.0
 */
public class PhotoListActivity extends AppCompatActivity {
    PhotoHandler photoHandler;

    /**
     * Called when the activity starts
     *
     * @param savedInstanceState The saved instance state
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_photo_list);

        // Make sure there is a valid token, if not - send user back to login
        if (Models.TOKEN.equals("")) {
            startActivity(new Intent(this, MainActivity.class));
            this.finish();
        }

        // Create a new PhotoHandler and get photos
        photoHandler = new PhotoHandler(this);
        photoHandler.getPhotos();
    }

    /**
     * Add action buttons to the action bar
     *
     * @param menu The menu object
     * @return Always return true
     */
    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        getMenuInflater().inflate(R.menu.menu_photo_list, menu);
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
            case R.id.gallery_action_button:
                Log.d(Models.FEEDBACK_SUCCESS, "Gallery action button pressed, calling upload");
                photoHandler.galleryActivity();
                return true;
            case R.id.camera_action_button:
                Log.d(Models.FEEDBACK_SUCCESS, "Camera action button pressed, calling upload");
                photoHandler.cameraActivity();
                return true;
            default:
                return super.onOptionsItemSelected(item);
        }
    }

    /**
     * Get the photo from the gallery or from the camera
     *
     * @param requestCode The request code (which action to be executed)
     * @param resultCode The result code (for checking if the result is ok)
     * @param data The photo to be handled
     */
    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);

        if (requestCode == Models.GALLERY_PHOTO_REQUEST && resultCode == RESULT_OK) {
            // The user wants to pick an image from the gallery
            photoHandler.handlePhotoFromGallery(data);
        } else if (requestCode == Models.CAMERA_PHOTO_REQUEST && resultCode == RESULT_OK) {
            // The user wants to take a photo with the camera
            photoHandler.handleCamera(data);
        }
    }
}
