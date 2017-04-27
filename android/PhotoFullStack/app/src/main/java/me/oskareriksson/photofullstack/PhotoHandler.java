package me.oskareriksson.photofullstack;

import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.provider.MediaStore;
import android.support.v7.app.AppCompatActivity;
import android.util.Base64;
import android.util.Log;
import android.widget.ImageView;

import org.json.JSONException;
import org.json.JSONObject;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.util.concurrent.ExecutionException;

/**
 * The PhotoHandler class handles taking, picking, uploading and getting photos
 *
 * @author Oskar Eriksson
 * @version 1.0
 */
public class PhotoHandler {
    AppCompatActivity activity;

    /**
     * Constructor for the PhotoHandler
     *
     * @param activity The activity that the PhotoHandler belongs to
     */
    public PhotoHandler(AppCompatActivity activity) {
        this.activity = activity;
    }

    /**
     * Get photos from the backend
     */
    public void getPhotos() {
        Request.Response response = sendRequest("/get", null);
        Log.d(Models.FEEDBACK_SUCCESS, "Got photos");
        // TODO Do something with photos
    }

    /**
     * Called when the user clicks the gallery action button,
     * opens the gallery and lets user pick photo
     */
    public void galleryActivity() {
        Intent pickPhotoIntent = new Intent();
        pickPhotoIntent.setType("image/*");
        pickPhotoIntent.setAction(Intent.ACTION_GET_CONTENT);
        activity.startActivityForResult(Intent.createChooser(
                pickPhotoIntent, "Pick a photo"), Models.GALLERY_PHOTO_REQUEST);
    }

    /**
     * Called when the user clicks the camera action button,
     * opens the camera and lets user take a photo
     */
    public void cameraActivity() {
        Intent takePhotoIntent = new Intent(MediaStore.ACTION_IMAGE_CAPTURE);
        if (takePhotoIntent.resolveActivity(activity.getPackageManager()) != null) {
            activity.startActivityForResult(takePhotoIntent, Models.CAMERA_PHOTO_REQUEST);
        }
    }

    /**
     * Called when the user has successfully picked a photo from the gallery
     *
     * @param intent The intent containing the photo
     */
    public void handlePhotoFromGallery(Intent intent) {
        try {
            // Convert bitmap to Base64
            bitmapToBase64(MediaStore.Images.Media.getBitmap(
                    activity.getContentResolver(), intent.getData()));
        } catch (IOException e) {
            Log.d(Models.FEEDBACK_ERROR, "IOException handlePhotoFromGallery: " + e.getMessage());
        }
    }

    /**
     * Called when a user has successfully taken a picture with the camera
     *
     * @param intent The intent containing the photo
     */
    public void handleCamera(Intent intent) {
        // Convert bitmap to Base64
        bitmapToBase64((Bitmap) intent.getExtras().get("data"));
    }

    /**
     * Scales a bitmap to a valid size, compresses it to JPEG and converts it to a Base64 string
     *
     * @param bitmap The bitmap to be converted to Base64
     */
    private String bitmapToBase64(Bitmap bitmap) {
        // Scale the bitmap to the given limit of pixels
        double limit = 500000;
        double scale = Math.sqrt(limit/(bitmap.getWidth()*bitmap.getHeight()));
        bitmap = Bitmap.createScaledBitmap(bitmap, (int) Math.round(bitmap.getWidth()*scale),
                (int) Math.round(bitmap.getHeight()*scale), true);
        ByteArrayOutputStream outputStream = new ByteArrayOutputStream();
        bitmap.compress(Bitmap.CompressFormat.JPEG, 100, outputStream);
        String jpgBase64 = Base64.encodeToString(outputStream.toByteArray(), Base64.DEFAULT);
        Log.d(Models.FEEDBACK_SUCCESS, "Bitmap converted to Base64");

        displayBase64(jpgBase64);
        return jpgBase64;
    }

    private void displayBase64(String jpgBase64) {
        byte[] imageAsBytes = Base64.decode(jpgBase64.getBytes(), 0);
        ((ImageView) activity.findViewById(R.id.displayPhotoView)).setImageBitmap(
                BitmapFactory.decodeByteArray(imageAsBytes, 0, imageAsBytes.length));
    }

    /**
     * Sends a login or a registration request to the backend, depending on the specified API
     *
     * @param api The API to send the request to
     * @return The message returned from the Request class
     */
    public Request.Response sendRequest(String api, JSONObject json) {
        // Send request to the backend
        try {
            Request.Response response = new Request().execute(
                    "http://" + Models.IP + ":" + Models.PORT + api, json).get();

            return response;
        } catch (ExecutionException e) {
            Log.d(Models.FEEDBACK_ERROR, "PhotoList: ExecutionException " + e.getMessage());
            return null;
        } catch (InterruptedException e) {
            Log.d(Models.FEEDBACK_ERROR, "PhotoList: InterruptedException " + e.getMessage());
            return null;
        }
    }
}
