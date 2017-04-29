package me.oskareriksson.photofullstack;

import android.util.Log;

import org.json.JSONException;
import org.json.JSONObject;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Locale;

/**
 * A Photo object
 *
 * @author Oskar Eriksson
 * @version 1.0
 */
public class Photo {
    private String jpgBase64;
    private String title;
    private String description;
    private String date;
    private String user;

    /**
     * Constructor for creating a new Photo object
     *
     * @param jpgBase64 The JPEG in Base64 encoding
     */
    public Photo(String jpgBase64) {
        this.jpgBase64 = jpgBase64;
        date = new SimpleDateFormat("dd/MM/yyyy", Locale.getDefault()).format(new Date());
        title = "";
        description = "";
        user = Models.USER;
    }

    /**
     * Constructor for creating a Photo object from a JSON object
     *
     * @param jsonObject
     */
    public Photo(JSONObject jsonObject) {
        try {
            jpgBase64 = jsonObject.getString("jpgbase64");
            title = jsonObject.getString("title");
            description = jsonObject.getString("description");
            date = jsonObject.getString("date");
            user = jsonObject.getString("user");
        } catch (JSONException e) {
            Log.d(Models.FEEDBACK_ERROR, "JSONException creating Photo object: " + e.getMessage());
        }
    }

    /**
     * Sets the title field of the photo object
     *
     * @param title The title
     */
    public void setTitle(String title) {
        this.title = title;
    }

    /**
     * Sets the description field of the photo object
     *
     * @param description The description
     */
    public void setDescription(String description) {
        this.description = description;
    }

    /**
     * @return The Base64 representation of the JPEG image
     */
    public String getJpgBase64() {
        return jpgBase64;
    }

    /**
     * @return The title of the photo
     */
    public String getTitle() {
        return title;
    }

    /**
     * @return The description of the photo
     */
    public String getDescription() {
        return description;
    }

    /**
     * @return The date of the photo
     */
    public String getDate() {
        return date;
    }

    public JSONObject toJSON() {
        try {
            JSONObject json = new JSONObject();
            json.put("jpgbase64", jpgBase64);
            json.put("title", title);
            json.put("description", description);
            json.put("date", date);
            json.put("user", user);
            Log.d(Models.FEEDBACK_SUCCESS, "Photo JSON created successfully");
            return json;
        } catch (JSONException e) {
            Log.d(Models.FEEDBACK_ERROR, "JSONException in Photo toJSON:" + e.getMessage());
            return null;
        }
    }
}
