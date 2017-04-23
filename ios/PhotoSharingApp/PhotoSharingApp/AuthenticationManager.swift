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

    init() {
        // Assumed initial state, user needs to provide port and IP
        self.baseUrl = "http://"
        self.Token = ""
        self.ip = ""
        self.port = ""
    }

    func loginUser(username: String, password: String, completeCallback: @escaping (_ title: String, _ message: String, _ succesful: Bool) -> Void) {
        // Create JSON body of username and password
        let parameters: Parameters = ["username": username, "password": password]

        // Create url from parameters set in fields (by user from LoginView text fields)
        let url = baseUrl + ip + ":" + port + "/login"
        // Make request
        Alamofire.request(url, method: .post, parameters: parameters, encoding: JSONEncoding.default).validate().responseJSON { response in
            if response.response === nil {
                // Invalid url 
                completeCallback("Login failed", "Invalid IP or port, try again", false)
            }
            switch response.result {
            case .success(let value):
                let json = JSON(value)
                self.Token = json["token"].stringValue
                // Use callback closure to send back response
                completeCallback("Login Succesful", "Welcome!", true)
            case .failure:
                let statusCode = response.response?.statusCode
                if statusCode == 401 {
                    // Unathorized 
                    completeCallback("Login failed", "Invalid login credentials, please try again", false)
                }
                else {
                    // Bad request 
                    completeCallback("Login failed", "Invalid login credentials, please try again", false)
                }
            }
        }
    }

    func signupUser(username: String) {
    }
}
