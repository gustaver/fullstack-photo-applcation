//
//  PhotoEditingViewController.swift
//  PhotoSharingApp
//
//  Created by Gustave Rousselet on 2017-03-30.
//  Copyright © 2017 Gustave Rousselet. All rights reserved.
//

import UIKit

class PhotoEditingViewController: UIViewController, UITextFieldDelegate {
    
    // Image to be set from previous view controller (PhotoTableView)
    var image: UIImage!
    // ImageView displaying image taken by user (sent from previous view)
    @IBOutlet var imageView: UIImageView!
    
    // Title and description text fields to be filled by user 
    @IBOutlet var titleTextField: UITextField!
    @IBOutlet var descriptionTextField: UITextField!
    
    // Photo object, populated during editing and sent in request 
    var photo: Photo!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        // Set image passed from previous view controller (PhotoTableView)
        imageView.image = self.image
        
        // Create photo object from image 
        self.createPhoto()
        
        // Set delegates of all text fields 
        self.titleTextField.delegate = self
        self.descriptionTextField.delegate = self
        
        // Recognize tap on screen to dismiss keyboard
        let tap = UITapGestureRecognizer(target: self.view, action: #selector(UIView.endEditing(_:)))
        tap.cancelsTouchesInView = false
        self.view.addGestureRecognizer(tap)
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }
    
    @IBAction func onUploadPhotoPress(_ sender: Any) {
        // Check that there is input in both username field and password field
        if titleTextField.hasText && descriptionTextField.hasText {
            // Set the title and the description of the photo object in field
            photo.Title = titleTextField.text
            photo.Description = descriptionTextField.text
            
            // Callback function to be sent into upload request
            func callbackRequestComplete(title: String, message: String, succesful: Bool) {
                //TODO: Consider displaying an alert of succesful login here
                // Check if login was succesful
                if succesful {
                    self.onUploadSuccesful()
                }
                else {
                    displayAlert(title: title, alertText: message, buttonText: "Ok")
                }
            }
            
            // Make upload request, callback used on completion of request 
            PhotoManager.sharedInstance.uploadPhoto(completeCallback: callbackRequestComplete, photo: self.photo)
        }
        else {
            // One of the fields didn't have any text, display alert, no callback
            self.displayAlert(title: "Invalid request", alertText: "All photo parameters must be specified", buttonText: "Ok")
        }
    }
    
    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        // Check if title or description text field 
        if textField.returnKeyType == UIReturnKeyType.next {
            // Check for next responder
            if let nextField = textField.superview?.viewWithTag(textField.tag + 1) as? UITextField {
                nextField.becomeFirstResponder()
            } else {
                // Not found, so remove keyboard.
                textField.resignFirstResponder()
            }
        } else {
            // Description key, pressing done means upload
            // Dismiss keyboard 
            textField.resignFirstResponder()
            self.onUploadPhotoPress(self)
        }
        return false
    }
    
    func onUploadSuccesful() {
        // Callback for when OK is pressed on alert, pop back to photo table view
        func onPressOk(action: UIAlertAction) {
            navigationController?.popViewController(animated: true)
        }
        let alert = UIAlertController(title: "Upload successful", message: "Let's checkout your new photo!", preferredStyle: UIAlertControllerStyle.alert)
        alert.addAction(UIAlertAction(title: "Go!", style: UIAlertActionStyle.default, handler: onPressOk))
        self.present(alert, animated: true, completion: nil)
    }
    
    func displayAlert(title: String, alertText: String, buttonText: String) {
        let alert = UIAlertController(title: title, message: alertText, preferredStyle: UIAlertControllerStyle.alert)
        alert.addAction(UIAlertAction(title: buttonText, style: UIAlertActionStyle.default, handler: nil))
        self.present(alert, animated: true, completion: nil)
    }
    
    func createPhoto() {
        // Photo to be populated 
        let photo = Photo()
        
        // Get image base64 string
        if let base64String = UIImageJPEGRepresentation(self.image, 0.9)?.base64EncodedString() {
            // Get current date in correct format
            let date = Date()
            let formatter = DateFormatter()
            formatter.dateFormat = "dd/MM/yyyy"
            let photoDate = formatter.string(from: date)
            
            // Set photo fields from above values
            photo.JpgBase64 = base64String
            photo.Date = photoDate
            photo.User = AuthenticationManager.sharedInstance.username
            self.photo = photo
        }
    }
}
