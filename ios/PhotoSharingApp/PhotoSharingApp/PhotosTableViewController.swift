//
//  PhotosTableViewController.swift
//  PhotoSharingApp
//
//  Created by Gustave Rousselet on 2017-03-28.
//  Copyright © 2017 Gustave Rousselet. All rights reserved.
//

import UIKit

class PhotosTableViewController: UITableViewController, UIImagePickerControllerDelegate, UINavigationControllerDelegate {
    
    // Image that user takes in camera or in image picker, stored to be sent to photo editing view, does not need to be set
    private var image: UIImage!
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }
    
    override func viewWillAppear(_ animated: Bool) {
        // Load phtoos. List of photos should reload every time the view appears
        self.loadPhotos()
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
    }

    // MARK: - Table view data source

    override func numberOfSections(in tableView: UITableView) -> Int {
        // Only one section for photos
        return 1
    }

    override func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return PhotoManager.sharedInstance.PhotoArray.count
    }
    
    override func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        // Create cell
        let photoCell = tableView.dequeueReusableCell(withIdentifier: "PhotoCell") as! PhotoTableViewCell
        // Get current photo for easier reference 
        let photo = PhotoManager.sharedInstance.PhotoArray[indexPath.row]
        
        // Get base64 string of photo using current row as index for array
        let base64String = photo.JpgBase64
        // Decode Base64 string to image (if possible) and then set image
        if let decodedData = Data(base64Encoded: base64String!) {
            let image = UIImage(data: decodedData)
            photoCell.imageInCell.image = image
        }
        
        // Set remaining labels of photoCell
        photoCell.titleLabel.text = photo.Title
        photoCell.dateLabel.text = photo.Date
        photoCell.descriptionLabel.text = photo.Description
        return photoCell
    }
    
    @IBAction func onCameraButtonPress(_ sender: Any) {
        // Present the camera if ready
        if UIImagePickerController.isSourceTypeAvailable(UIImagePickerControllerSourceType.camera){
            let cameraView = UIImagePickerController()
            cameraView.delegate = self
            cameraView.sourceType = UIImagePickerControllerSourceType.camera;
            cameraView.allowsEditing = false
            self.present(cameraView, animated: true, completion: nil)
        }
    }
    
    @IBAction func onGalleryButtonPress(_ sender: Any) {
        if UIImagePickerController.isSourceTypeAvailable(UIImagePickerControllerSourceType.photoLibrary){
            let imagePicker = UIImagePickerController()
            imagePicker.delegate = self
            imagePicker.sourceType = UIImagePickerControllerSourceType.photoLibrary;
            imagePicker.allowsEditing = true
            self.present(imagePicker, animated: true, completion: nil)
        }
    }
    
    func imagePickerController(_ picker: UIImagePickerController, didFinishPickingImage image: UIImage!, editingInfo: [NSObject : AnyObject]!){
        self.image = image
        // Dismiss view
        self.dismiss(animated: true, completion: nil)
        // Perform segue to photo editing view
        performSegue(withIdentifier: "PhotoEditingView", sender: self)
    }
    
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        // Check the segue to be performed
        if segue.identifier == "PhotoEditingView" {
            // Check that next view controller can be set
            if let photoEditingViewController = segue.destination as? PhotoEditingViewController{
                // Set image of photo editing view to image picked by user
                photoEditingViewController.image = self.image
            }
        }
    }
    
    func loadPhotos() {
        // Display activity indicator and disable view perhaps
        
        // Create callback function to be called when request finishes
        func callbackPhotoRequest(succes: Bool) {
            if succes {
                // Request was succesful, load photos into tableview 
                self.onPhotosSuccessfullyLoaded()
            } else {
                // Request was unsuccesful, display error message 
                self.displayAlert(title: "Could not load photos", alertText: "Please logout and try again", buttonText: "Ok")
            }
        }
        PhotoManager.sharedInstance.getPhotos(completeCallback: callbackPhotoRequest)
    }
    
    func onPhotosSuccessfullyLoaded() {
        tableView.reloadData()
    }
    
    func displayAlert(title: String, alertText: String, buttonText: String) {
        let alert = UIAlertController(title: title, message: alertText, preferredStyle: UIAlertControllerStyle.alert)
        alert.addAction(UIAlertAction(title: buttonText, style: UIAlertActionStyle.default, handler: nil))
        self.present(alert, animated: true, completion: nil)
    }

    /*
    // Override to support conditional editing of the table view.
    override func tableView(_ tableView: UITableView, canEditRowAt indexPath: IndexPath) -> Bool {
        // Return false if you do not want the specified item to be editable.
        return true
    }
    */

    /*
    // Override to support editing the table view.
    override func tableView(_ tableView: UITableView, commit editingStyle: UITableViewCellEditingStyle, forRowAt indexPath: IndexPath) {
        if editingStyle == .delete {
            // Delete the row from the data source
            tableView.deleteRows(at: [indexPath], with: .fade)
        } else if editingStyle == .insert {
            // Create a new instance of the appropriate class, insert it into the array, and add a new row to the table view
        }    
    }
    */

    /*
    // Override to support rearranging the table view.
    override func tableView(_ tableView: UITableView, moveRowAt fromIndexPath: IndexPath, to: IndexPath) {

    }
    */

    /*
    // Override to support conditional rearranging of the table view.
    override func tableView(_ tableView: UITableView, canMoveRowAt indexPath: IndexPath) -> Bool {
        // Return false if you do not want the item to be re-orderable.
        return true
    }
    */
}
