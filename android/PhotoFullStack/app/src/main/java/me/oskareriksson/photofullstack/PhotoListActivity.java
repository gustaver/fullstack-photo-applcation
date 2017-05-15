package me.oskareriksson.photofullstack;

import android.app.AlertDialog;
import android.content.DialogInterface;
import android.content.Intent;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ListView;

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

        // Populate the ListView with photos using the custom adapter
        ListAdapter listAdapter = new ListAdapter(this, Models.PHOTO_ARRAYLIST);
        ((ListView) findViewById(R.id.photo_listview)).setAdapter(listAdapter);

        // Set a long click listener for the photos in the list, which spawns a delete dialog
        ((ListView) findViewById(R.id.photo_listview)).setOnItemLongClickListener(
                new AdapterView.OnItemLongClickListener() {
                    @Override
                    public boolean onItemLongClick(AdapterView<?> adapter, View view, int pos, long id) {
                        PhotoListActivity.this.generateDeleteDialog(pos).show();
                        return true;
                    }
                });
    }

    /**
     * Generate a dialog that asks the user if they want to delete the photo
     *
     * @param pos The position of the photo to be deleted
     * @return A delete dialog with the option to delete the given photo
     */
    private AlertDialog generateDeleteDialog(final int pos) {
        AlertDialog.Builder alertBuilder = new AlertDialog.Builder(this);


        alertBuilder.setTitle(Models.PHOTO_ARRAYLIST.get(pos).getTitle());
        alertBuilder.setMessage("Do you want to delete the photo?");

        // Delete photo if the user clicks yes
        alertBuilder.setPositiveButton("YES", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                // Delete the photo
                PhotoListActivity.this.photoHandler.deletePhoto(pos);

                // Refresh the PhotoListActivity
                PhotoListActivity.this.finish();
                startActivity(PhotoListActivity.this.getIntent());
                dialog.dismiss();
            }
        });

        // Dismiss dialog if the user clicks no
        alertBuilder.setNegativeButton("NO", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                dialog.dismiss();
            }

        });

        return alertBuilder.create();
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
                Log.d(Models.FEEDBACK_SUCCESS, "Gallery action button pressed");
                photoHandler.galleryActivity();
                return true;
            case R.id.camera_action_button:
                Log.d(Models.FEEDBACK_SUCCESS, "Camera action button pressed");
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

        if ((requestCode == Models.GALLERY_PHOTO_REQUEST
                || requestCode == Models.CAMERA_PHOTO_REQUEST) && resultCode == RESULT_OK) {
            // The user wants to upload a photo from the gallery/camera
            photoHandler.handlePhoto(requestCode, data);
        }
    }
}
