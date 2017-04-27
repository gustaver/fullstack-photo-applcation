package me.oskareriksson.photofullstack;

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

    public Photo(String jpgBase64, String title, String description, String date) {
        this.jpgBase64 = jpgBase64;
        this.title = title;
        this.description = description;
        this.date = date;
    }

}
