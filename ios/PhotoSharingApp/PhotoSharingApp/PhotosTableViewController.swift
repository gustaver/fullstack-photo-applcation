//
//  PhotosTableViewController.swift
//  PhotoSharingApp
//
//  Created by Gustave Rousselet on 2017-03-28.
//  Copyright © 2017 Gustave Rousselet. All rights reserved.
//

import UIKit

class PhotosTableViewController: UITableViewController, UIImagePickerControllerDelegate, UINavigationControllerDelegate {
    
    // Activity indicator to show when loading photos
    var activityIndicator: UIActivityIndicatorView = UIActivityIndicatorView()
    
    // Image that user takes in camera or in image picker, stored to be sent to photo editing view, does not need to be set
    private var image: UIImage!
    @IBOutlet var galleryBarButton: UIBarButtonItem!
    @IBOutlet var cameraBarButton: UIBarButtonItem!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        // Setup allowing to delete rows in photo table view 
        tableView.allowsSelectionDuringEditing = true
        
        // Setup for Activity Indicator 
        activityIndicator.center = self.view.center
        activityIndicator.hidesWhenStopped = true
        activityIndicator.activityIndicatorViewStyle = UIActivityIndicatorViewStyle.whiteLarge
        activityIndicator.color = UIColor.darkGray
        view.addSubview(activityIndicator)
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
        if let decodedData = Data(base64Encoded: base64String!, options: Data.Base64DecodingOptions.ignoreUnknownCharacters) {
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
    
    override func tableView(_ tableView: UITableView, canEditRowAt indexPath: IndexPath) -> Bool {
        return true
    }
    
    override func tableView(_ tableView: UITableView, commit editingStyle: UITableViewCellEditingStyle, forRowAt indexPath: IndexPath) {
        // Disable camera and gallery buttons
        cameraBarButton.isEnabled = false
        galleryBarButton.isEnabled = false
        
        // Create callback for reqest 
        func callbackPhotoRequest (title: String, message: String, success: Bool) {
            // Enable camera and gallery buttons
            cameraBarButton.isEnabled = true
            galleryBarButton.isEnabled = true
            
            if success {
                // Request was succesful, photo has been removed from array, display message and reload tableView
                self.displayAlert(title: title, alertText: message, buttonText: "Ok")
                self.onPhotosSuccessfullyLoaded()
            } else {
                // Request was unsuccesful, display error message from request callback
                self.displayAlert(title: title, alertText: message, buttonText: "Ok")
            }
        }
        PhotoManager.sharedInstance.removePhoto(completeCallback: callbackPhotoRequest, index: indexPath.row)
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
        // Show loading indicator
        activityIndicator.startAnimating()
        
        // Disable camera and gallery buttons 
        cameraBarButton.isEnabled = false
        galleryBarButton.isEnabled = false
        
        // Create callback function to be called when request finishes
        func callbackPhotoRequest(success: Bool) {
            // Stop loading and turn on interaction with application 
            activityIndicator.stopAnimating()
            
            // Enable camera and gallery buttons
            cameraBarButton.isEnabled = true
            galleryBarButton.isEnabled = true
            
            if success {
                // Request was succesful, load photos into tableview
                self.onPhotosSuccessfullyLoaded()
                if PhotoManager.sharedInstance.PhotoArray.isEmpty {
                    // No photos, display message
                    self.displayAlert(title: "Whoops!", alertText: "Looks like you don't have any photos, go ahead and add some!", buttonText: "Ok")
                }
            } else {
                // Request was unsuccesful, display error message
                self.displayAlert(title: "Could not load photos", alertText: "Please logout and try again", buttonText: "Ok")
            }
        }
        PhotoManager.sharedInstance.getPhotos(completeCallback: callbackPhotoRequest)
    }
    
    func onPhotosSuccessfullyLoaded() {
        DispatchQueue.main.async {
            self.tableView.reloadData()
        }
    }
    
    func displayAlert(title: String, alertText: String, buttonText: String) {
        let alert = UIAlertController(title: title, message: alertText, preferredStyle: UIAlertControllerStyle.alert)
        alert.addAction(UIAlertAction(title: buttonText, style: UIAlertActionStyle.default, handler: nil))
        self.present(alert, animated: true, completion: nil)
    }
}
