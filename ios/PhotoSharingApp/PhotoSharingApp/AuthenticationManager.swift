//
// Created by Gustave Rousselet on 2017-03-27.
// Copyright (c) 2017 Gustave Rousselet. All rights reserved.
//

import Foundation
import SwiftyJSON
import Alamofire

class AuthenticationManager {

    // Global variable "singleton"
    static let sharedInstance = AuthenticationManager()

    // API call settings
    var baseUrl: String
    var ip: String
    var port: String

    // Token for future API calls after login
    var Token: String
    
    // Currently logged in user 
    var username: String

    init() {
        // Assumed initial state, user needs to provide port and IP
        self.baseUrl = "http://"
        self.Token = ""
        self.ip = ""
        self.port = ""
        self.username = ""
    }

    func loginUser(username: String, password: String, completeCallback: @escaping (_ title: String, _ message: String, _ succesful: Bool) -> Void) {
        // Clear Token and username for each login
        self.Token = ""
        self.username = ""
        
        // Create JSON body of username and password
        let parameters: Parameters = ["username": username, "password": password]

        // Create url from parameters set in fields (by user from LoginView text fields)
        let url = baseUrl + ip + ":" + port + "/login"
        
        // Make request
        Alamofire.request(url, method: .post, parameters: parameters, encoding: JSONEncoding.default).response { response in
            let responseCode = response.response?.statusCode
            if responseCode == nil {
                // Invalid url 
                completeCallback("Login failed", "Invalid IP or port, try again", false)
            } else {
                if responseCode == 200 {
                    // Check decode response data to JSON
                    let json = JSON(response.data!)
                    if json == JSON.null {
                        // Error decoding JSON response 
                        completeCallback("Login failed", "Bad request", false)
                    } else {
                        // Request was succesful, set values
                        self.Token = json["token"].stringValue
                        self.username = username
                        // Use callback closure to send back response
                        completeCallback("Login Succesful", "Welcome!", true)
                    }
                } else if responseCode == 401 {
                    // Unathorized
                    completeCallback("Login failed", "Invalid login credentials, please try again", false)
                } else {
                    // Bad request
                    completeCallback("Login failed", "Invalid login credentials, please try again", false)
                }
            }
        }
    }

    func signupUser(username: String, password: String, completeCallback: @escaping (_ title: String, _ message: String, _ succesful: Bool) -> Void) {
        // Clear username for each signup 
        self.username = ""
        
        // Create JSON body of username and password
        let parameters: Parameters = ["username": username, "password": password]
        
        // Create url from parameters set in fields (by user from LoginView text fields)
        let url = baseUrl + ip + ":" + port + "/signup"
        
        // Make request
        Alamofire.request(url, method: .post, parameters: parameters, encoding: JSONEncoding.default).response { response in
            let responseCode = response.response?.statusCode
            if responseCode == nil {
                // Invalid url
                completeCallback("Singup failed", "Invalid IP or port, try again", false)
            } else {
                if responseCode == 200 {
                    // Signup successful 
                    completeCallback("Signup successful", "Welcome to PhotoSharing! Go ahead and login!", true)
                } else if responseCode == 401 {
                    // Unathorized
                    completeCallback("Signup failed", "Invalid signup credentials, please try again", false)
                } else {
                    // Bad request
                    completeCallback("Signup failed", "Invalid signup credentials, please try again", false)
                }
            }
        }
    }
}
