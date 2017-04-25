//
//  LoginViewController.swift
//  PhotoSharingApp
//
//  Created by Gustave Rousselet on 2017-03-24.
//  Copyright © 2017 Gustave Rousselet. All rights reserved.
//

import UIKit

class LoginViewController: UIViewController, UITextFieldDelegate {

    @IBOutlet var usernameTextInput: UITextField!
    @IBOutlet var passwordTextInput: UITextField!
    @IBOutlet var portTextInput: UITextField!
    @IBOutlet var ipTextInput: UITextField!
    
    
    override func viewDidLoad() {
        super.viewDidLoad()
        usernameTextInput.text = "gustave"
        passwordTextInput.text = "12345"
        portTextInput.text = "8080"
        ipTextInput.text = "192.168.1.2"
        
        // Set view controller to delegate for all text inputs
        self.usernameTextInput.delegate = self
        self.passwordTextInput.delegate = self
        self.portTextInput.delegate = self
        self.ipTextInput.delegate = self
        
        // Recognize tap on screen to dismiss keyboard
        let tap = UITapGestureRecognizer(target: self.view, action: #selector(UIView.endEditing(_:)))
        tap.cancelsTouchesInView = false
        self.view.addGestureRecognizer(tap)
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
    }

    @IBAction func onPressLogin(_ sender: Any) {
        // Check that there is input in both username field and password field
        if usernameTextInput.hasText && passwordTextInput.hasText {
            // Set the port and IP of AuthenticationManager (if bad, request will fail anyway) 
            AuthenticationManager.sharedInstance.port = portTextInput.text!
            AuthenticationManager.sharedInstance.ip = ipTextInput.text!
            // Callback function to be sent into authentication request
            func callbackRequestComplete(title: String, message: String, succesful: Bool) {
                if succesful {
                    self.onLoginSuccesful()
                } else {
                    displayAlert(title: title, alertText: message, buttonText: "OK")
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
    
    @IBAction func onPressSignup(_ sender: Any) {
        // Check that there is input in both username and password fields 
        if usernameTextInput.hasText && passwordTextInput.hasText {
            // Set the port and IP of AuthenticationManager (if bad, request will fail anyway)
            AuthenticationManager.sharedInstance.port = portTextInput.text!
            AuthenticationManager.sharedInstance.ip = ipTextInput.text!
            // Callback function to be sent into authentication request
            func callbackRequestComplete(title: String, message: String, succesful: Bool) {
                if succesful {
                    // Display alert of successful signup returned from request
                    displayAlert(title: title, alertText: message, buttonText: "Ok")
                } else {
                    // Display alert of unsuccessful singup returned from request
                    displayAlert(title: title, alertText: message, buttonText: "Ok")
                }
            }
            // Make authentication request, if succesful, callback will display accordingly
            AuthenticationManager.sharedInstance.signupUser(username: usernameTextInput.text!, password: passwordTextInput.text!, completeCallback: callbackRequestComplete)
        }
        else {
            // One of the fields didn't have any text, display alert, no callback
            self.displayAlert(title: "Invalid signup", alertText: "All signup parameters must be specified", buttonText: "Ok")
        }
    }
    
    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        if textField.returnKeyType == UIReturnKeyType.next {
            // Try to find next responder
            if let nextField = textField.superview?.viewWithTag(textField.tag + 1) as? UITextField {
                nextField.becomeFirstResponder()
            } else {
                // Not found, so remove keyboard.
                textField.resignFirstResponder()
            }
        } else {
            // Return key = GO 
            // Close keyboard and login 
            textField.resignFirstResponder()
            self.onPressLogin(self)
        }
        return false
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
        performSegue(withIdentifier: "ShowPhotoTableView", sender: self)
    }
}

