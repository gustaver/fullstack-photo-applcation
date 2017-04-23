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

    // Token for future API calls after login
    var Token: String

    init() {
        self.baseUrl = "http://192.168.1.3:8080"
        self.Token = ""
    }

    func loginUser(username: String, password: String, complete: @escaping (_ response: String) -> Void) {
        // Create JSON body of username and password
        let parameters: Parameters = ["username": username, "password": password]

        // Create url
        let url = baseUrl + "/login"
        // Make request
        Alamofire.request(url, method: .post, parameters: parameters, encoding: JSONEncoding.default).validate().responseJSON { response in
            switch response.result {
            case .success(let value):
                let json = JSON(value)
                self.Token = json["token"].stringValue
                // Use callback closure to send back response
                complete("Login Succesful")
            case .failure:
                let statusCode = response.response?.statusCode
                if statusCode == 401 {
                    // Unathorized 
                    complete("Invalid login credentials")
                }
                else {
                    // Bad request 
                    complete("Invalid login credentials")
                }
            }
        }
    }

    func signupUser(username: String) {
    }
}
