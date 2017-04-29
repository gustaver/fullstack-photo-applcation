package me.oskareriksson.photofullstack;

import android.content.Intent;
import android.graphics.Bitmap;
import android.provider.MediaStore;
import android.support.v7.app.AppCompatActivity;
import android.util.Base64;
import android.util.Log;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.net.HttpURLConnection;
import java.util.ArrayList;
import java.util.Collections;
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

    public PhotoHandler() { this.activity = null; }

    /**
     * Get photos from the backend
     */
    public void getPhotos() {
        try {
            Request.Response response = sendRequest("/get", null);
            Log.d(Models.FEEDBACK_SUCCESS, "Got photos from backend");
            JSONArray jsonArray = new JSONArray(response.getResponseString());

            // Clear and populate the photos ArrayList
            Models.PHOTO_ARRAYLIST = new ArrayList<>();
            for (int i = 0; i < jsonArray.length(); i++) {
                // Create a new Photo object and add to the ArrayList
                Models.PHOTO_ARRAYLIST.add(new Photo(jsonArray.getJSONObject(i)));
            }
            Collections.reverse(Models.PHOTO_ARRAYLIST);
        } catch (JSONException e) {
            Log.d(Models.FEEDBACK_ERROR, "JSONException getting photos: " + e.getMessage());
        }
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
    public void handlePhoto(int requestType, Intent intent) {
        // Try to convert the photo in the intent to a Base64 string
        String jpgBase64;
        try {
            if (requestType == Models.GALLERY_PHOTO_REQUEST) {
                // The photo is from the gallery, get the Base64 string
                jpgBase64 = bitmapToBase64(MediaStore.Images.Media.getBitmap(
                        activity.getContentResolver(), intent.getData()));
            } else if (requestType == Models.CAMERA_PHOTO_REQUEST) {
                // The photo is from the camera, get the Base64 string
                jpgBase64 = bitmapToBase64((Bitmap) intent.getExtras().get("data"));
            } else {
                // The request type was not valid, return
                return;
            }
        } catch (IOException e) {
            Log.d(Models.FEEDBACK_ERROR, "IOException handlePhoto: " + e.getMessage());
            return;
        }

        // Create a new photo and a new intent to edit the photo
        Models.PHOTO_EDIT = new Photo(jpgBase64);
        activity.startActivity(new Intent(activity, EditorActivity.class));
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

        // Compress the bitmap to a jpg
        ByteArrayOutputStream outputStream = new ByteArrayOutputStream();
        bitmap.compress(Bitmap.CompressFormat.JPEG, 100, outputStream);

        // Convert jpg to Base64
        String jpgBase64 = Base64.encodeToString(outputStream.toByteArray(), Base64.DEFAULT);
        Log.d(Models.FEEDBACK_SUCCESS, "Bitmap converted to Base64");

        return jpgBase64;
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

    /**
     * Deletes a given photo
     *
     * @param pos The position of the Photo object in the ArrayList of photos
     */
    public void deletePhoto(int pos) {
        // Get the photo at the given position
        Photo photo = Models.PHOTO_ARRAYLIST.get(pos);

        // If there is no photo to be deleted, return
        if (photo == null) {
            return;
        }

        // Try to send the request to upload the photo
        Request.Response response;
        try {
            response = new Request().execute(
                    "http://" + Models.IP + ":" + Models.PORT + "/remove", photo.toJSON()).get();
            if (response.getResponseCode() == HttpURLConnection.HTTP_OK) {
                Log.d(Models.FEEDBACK_SUCCESS, "Successfully deleted photo");
            } else {
                // Invalid response
                Log.d(Models.FEEDBACK_ERROR, "HTTP response not OK");
            }
        } catch (InterruptedException e) {
            Log.d(Models.FEEDBACK_ERROR, "InterruptedException delete photo: " + e.getMessage());
        } catch (ExecutionException e) {
            Log.d(Models.FEEDBACK_ERROR, "ExecutionException delete photo: " + e.getMessage());
        }
    }
}
