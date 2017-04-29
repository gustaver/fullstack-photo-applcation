package me.oskareriksson.photofullstack;

import java.util.ArrayList;
import java.util.Map;

/**
 * Holds constants and static methods used throughout the application
 *
 * @author Oskar Eriksson
 * @version 1.0
 */
public class Models {
    public static String IP = "";
    public static String PORT = "";
    public static String TOKEN = "";
    public static String USER = "";

    public static final String FEEDBACK_SUCCESS = "SUCCESS";
    public static final String FEEDBACK_ERROR = "ERROR";

    public static final int GALLERY_PHOTO_REQUEST = 0;
    public static final int CAMERA_PHOTO_REQUEST = 1;

    public static ArrayList<Photo> PHOTO_ARRAYLIST = new ArrayList<>();
    public static Photo PHOTO_EDIT = null;


    /**
     * Make sure that all values in a map is set
     *
     * @param input The map to be checked
     * @return True if all values were set
     */
    public static boolean validateInput(Map<String, String> input) {
        for (String value : input.values()) {
            if (value.length() == 0) {
                return false;
            }
        }
        return true;
    }
}