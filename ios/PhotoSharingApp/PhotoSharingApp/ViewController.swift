//
//  ViewController.swift
//  PhotoSharingApp
//
//  Created by Gustave Rousselet on 2017-03-24.
//  Copyright © 2017 Gustave Rousselet. All rights reserved.
//

import UIKit

class ViewController: UIViewController {

    @IBOutlet var usernameTextInput: UITextField!
    @IBOutlet var passwordTextInput: UITextField!
    
    
    override func viewDidLoad() {
        super.viewDidLoad()
        // Do any additional setup after loading the view, typically from a nib.
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }

    @IBAction func onPressLogin(_ sender: Any) {
        // Check that there is input in both username field and password field
        if usernameTextInput.hasText && passwordTextInput.hasText {
            AuthenticationManager.sharedInstance.loginUser(username: usernameTextInput.text!, password: passwordTextInput.text!) { response in
                // Callback from response
                // Display message from callback 
                self.displayAlert(title: "Login", alertText: response, buttonText: "Ok")
            }
        }
        else {
            self.displayAlert(title: "Invalid login", alertText: "All login parameters must be specified", buttonText: "Ok")
        }
    }
    
    func displayAlert(title: String, alertText: String, buttonText: String) {
        let alert = UIAlertController(title: title, message: alertText, preferredStyle: UIAlertControllerStyle.alert)
        alert.addAction(UIAlertAction(title: buttonText, style: UIAlertActionStyle.default, handler: nil))
        self.present(alert, animated: true, completion: nil)
    }
}

