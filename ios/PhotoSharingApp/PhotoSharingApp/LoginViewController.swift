//
//  LoginViewController.swift
//  PhotoSharingApp
//
//  Created by Gustave Rousselet on 2017-03-24.
//  Copyright © 2017 Gustave Rousselet. All rights reserved.
//

import UIKit

class LoginViewController: UIViewController {

    @IBOutlet var usernameTextInput: UITextField!
    @IBOutlet var passwordTextInput: UITextField!
    @IBOutlet var portTextInput: UITextField!
    @IBOutlet var ipTextInput: UITextField!
    
    
    override func viewDidLoad() {
        super.viewDidLoad()
        // For testing purposes (for now) go straight to login, TODO: REMOVE THIS
        self.onLoginSuccesful()
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }

    @IBAction func onPressLogin(_ sender: Any) {
        // Check that there is input in both username field and password field
        if usernameTextInput.hasText && passwordTextInput.hasText {
            // Set the port and IP of AuthenticationManager (if bad, request will fail anyway) 
            AuthenticationManager.sharedInstance.port = portTextInput.text!
            AuthenticationManager.sharedInstance.ip = ipTextInput.text!
            // Callback function to be sent into authentication request
            func callbackRequestComplete(title: String, message: String, succesful: Bool) {
                //TODO: Consider displaying an alert of succesful login here
                // Check if login was succesful 
                if succesful {
                    self.onLoginSuccesful()
                }
            }
            // Make authentication request, if succesful, callback will display accordingly
            AuthenticationManager.sharedInstance.loginUser(username: usernameTextInput.text!, password: passwordTextInput.text!, completeCallback: callbackRequestComplete)
        }
        else {
            // One of the fields didn't have any text, display alert, no callback
            self.displayAlert(title: "Invalid login", alertText: "All login parameters must be specified", buttonText: "Ok")
        }
    }
    
    func displayAlert(title: String, alertText: String, buttonText: String) {
        let alert = UIAlertController(title: title, message: alertText, preferredStyle: UIAlertControllerStyle.alert)
        alert.addAction(UIAlertAction(title: buttonText, style: UIAlertActionStyle.default, handler: nil))
        self.present(alert, animated: true, completion: nil)
    }
    
    func onLoginSuccesful() {
        // Clear all text fields 
        usernameTextInput.text = ""
        passwordTextInput.text = ""
        portTextInput.text = ""
        ipTextInput.text = ""
        performSegue(withIdentifier: "ShowPhotoTableView", sender: nil)
    }
}

