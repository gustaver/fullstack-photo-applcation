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

    func getPhotos(completeCallback: @escaping (_ success: Bool) -> Void) {
        // Clear PhotosArray for each request
        PhotoArray.removeAll()
        // Create headers with Token
        let headers: HTTPHeaders = ["Token": AuthenticationManager.sharedInstance.Token]
        // Create url for request 
        let url = AuthenticationManager.sharedInstance.baseUrl + AuthenticationManager.sharedInstance.ip + ":" + AuthenticationManager.sharedInstance.port + "/get"
        // Make get request
        Alamofire.request(url, method: .get, headers: headers).responseJSON { response in
            if response.response === nil {
                // Invalid url, request unsuccesful
                completeCallback(false)
            }
            // Check the result of the response and handle accordingly
            switch response.result {
            case .success(let value):
                let json = JSON(value)
                // Go through JSON response and add photos (know that JSON response will be an array if 200 OK)
                for (_, photo) in json {
                    let jsonPhoto: Photo = Photo(data: photo)
                    self.PhotoArray.append(jsonPhoto)
                }
                // Request succesful, toggle callback
                completeCallback(true)
            case .failure(let error):
                print(error)
                completeCallback(false)
            }
        }
    }

    func uploadPhoto() {
    }

    func removePhoto() {
    }
}
