//
// Created by Gustave Rousselet on 2017-03-27.
// Copyright (c) 2017 Gustave Rousselet. All rights reserved.
//

import Foundation
import SwiftyJSON
import Alamofire

class PhotoManager {

    static let sharedInstance = PhotoManager()
    var PhotoArray = [Photo]()

    init() {
    }

    func getPhotos() {
        // Create headers with Token
        let headers: HTTPHeaders = ["Token": AuthenticationManager.sharedInstance.Token]
        // Create url for request 
        let url = AuthenticationManager.sharedInstance.baseUrl + AuthenticationManager.sharedInstance.ip + ":" + AuthenticationManager.sharedInstance.port + "/get"
        // Make get request
        Alamofire.request(url, method: .get, headers: headers).responseJSON { response in
            if response.response === nil {
                // Invalid url
            }
            // Check the result of the response and handle accordingly
            switch response.result {
            case .success(let value):
                let json = JSON(value)
                for (_, photo) in json {
                    let jsonPhoto: Photo = Photo(data: photo)
                    self.PhotoArray.append(jsonPhoto)
                }
                print(self.PhotoArray.count)
            case .failure(let error):
                print(error)
            }
        }
    }

    func uploadPhoto() {
    }

    func removePhoto() {
    }
}
