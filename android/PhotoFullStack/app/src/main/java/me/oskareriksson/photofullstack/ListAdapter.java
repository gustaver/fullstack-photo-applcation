package me.oskareriksson.photofullstack;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.util.Base64;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ImageView;
import android.widget.TextView;

import java.net.HttpURLConnection;
import java.util.ArrayList;
import java.util.Map;
import java.util.concurrent.ExecutionException;

/**
 * ListAdapter is used to populate the ListView of Photo objects in PhotoListActivity
 *
 * @author Oskar Eriksson
 * @version 1.0
 */
public class ListAdapter extends ArrayAdapter<Photo> {
    PhotoHandler photoHandler;

    /**
     * Constructor for the ListAdapter class
     *
     * @param context The context where the ListAdapter is
     * @param photos The ArrayList of Photo objects to populate the ListView
     */
    public ListAdapter(Context context, ArrayList<Photo> photos) {
        super(context, 0, photos);

        photoHandler = new PhotoHandler();
    }

    /**
     * Used to populate the ListView
     *
     * @param position The item's position in the ArrayList
     * @param convertView The current view
     * @param parent The parent ViewGroup
     * @return The convertView, the current view
     */
    @Override
    public View getView(int position, View convertView, final ViewGroup parent) {
        final Photo photo = getItem(position);

        if (convertView == null) {
            convertView = LayoutInflater.from(
                    getContext()).inflate(R.layout.item_photolist, parent, false);
        }

        ImageView imageView = (ImageView) convertView.findViewById(R.id.photolist_item_image);

        // Convert Base64 string to bitmap
        if (photo != null) {
            byte[] imageAsBytes = Base64.decode(photo.getJpgBase64().getBytes(), Base64.DEFAULT);
            Bitmap bitmap = BitmapFactory.decodeByteArray(imageAsBytes, 0, imageAsBytes.length);
            imageView.setImageBitmap(bitmap);
        }

        ((TextView) convertView.findViewById(R.id.photolist_item_title)).setText(photo.getTitle());
        ((TextView) convertView.findViewById(R.id.photolist_item_description)).setText(photo.getDescription());
        ((TextView) convertView.findViewById(R.id.photolist_item_date)).setText(photo.getDate());

        return convertView;
    }
}
