//
// Created by Gustave Rousselet on 2017-03-27.
// Copyright (c) 2017 Gustave Rousselet. All rights reserved.
//

import Foundation
import SwiftyJSON
import Alamofire

// Helper class (singleton) photomanager. Handles all requests which have to do with photos (get, upload, remove). Stores photos (array).
class PhotoManager {

    static let sharedInstance = PhotoManager()
    var PhotoArray = [Photo]()

    init() {
    }

    // Method used when getting photos for photo list. Uses logged in user to make request to backend to get all their photos from database.
    func getPhotos(completeCallback: @escaping (_ success: Bool) -> Void) {
        // Create headers with Token
        let headers: HTTPHeaders = ["Token": AuthenticationManager.sharedInstance.Token]
        // Create url for request 
        let url = AuthenticationManager.sharedInstance.baseUrl + AuthenticationManager.sharedInstance.ip + ":" + AuthenticationManager.sharedInstance.port + "/get"
        // Make get request
        Alamofire.request(url, method: .post, headers: headers).responseJSON { response in
            if response.response === nil {
                // Invalid url, request unsuccesful
                completeCallback(false)
            }
            // Check the result of the response and handle accordingly
            switch response.result {
            case .success(let value):
                // TODO: Consider changing this to only add new photos (arrays.map?)
                // Clear PhotosArray for succesful request and re-populate
                self.PhotoArray.removeAll()
                let json = JSON(value)
                // Go through JSON response and add photos (know that JSON response will be an array if 200 OK)
                for (_, photo) in json {
                    let jsonPhoto: Photo = Photo()
                    jsonPhoto.initialiseFromJSON(data: photo)
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

    // Method used when uploading photo. Request made to backend, with associated callback depeding on outcome, and index of photo 
    // in shared photo array.
    func uploadPhoto(completeCallback: @escaping (_ title: String, _ message: String, _ succesful: Bool) -> Void, photo: Photo) {
        // Create headers with Token
        let headers: HTTPHeaders = ["Token": AuthenticationManager.sharedInstance.Token]
        // Create url for request
        let url = AuthenticationManager.sharedInstance.baseUrl + AuthenticationManager.sharedInstance.ip + ":" + AuthenticationManager.sharedInstance.port + "/upload"
        // Create parameters to encode into JSON body 
        let parameters: Parameters = photo.toParameters()
        
        // Make request using above information
        Alamofire.request(url, method: .post, parameters: parameters, encoding: JSONEncoding.default,  headers: headers).validate().responseJSON { response in
            // FIXME: This entire request/response process needs to be fixed
            if response.response === nil {
                // Invalid url
                completeCallback("Upload failed", "Invalid IP or port, try again", false)
            }
            let statusCode = response.response?.statusCode
            if statusCode == 200 {
                completeCallback("Upload succesful", "Go checkout your new photo!", true)
            } else {
                completeCallback("Upload failed", "Invalid upload request, try again", false)
            }
        }
    }

    // Method used when removing photos. Request made to backend with callback, as well as index of photo in photo array.
    func removePhoto(completeCallback: @escaping (_ title: String, _ message: String, _ succesful: Bool) -> Void, index: Int) {
        // Create headers with Token
        let headers: HTTPHeaders = ["Token": AuthenticationManager.sharedInstance.Token]
        // Create url for request
        let url = AuthenticationManager.sharedInstance.baseUrl + AuthenticationManager.sharedInstance.ip + ":" + AuthenticationManager.sharedInstance.port + "/remove"
        // Get the photo to be removed 
        let photo: Photo = self.PhotoArray[index]
        let parameters: Parameters = photo.toParameters()
        
        // Make request using above information
        Alamofire.request(url, method: .post, parameters: parameters, encoding: JSONEncoding.default,  headers: headers).validate().responseJSON { response in
            // FIXME: This entire request/response process needs to be fixed
            if response.response === nil {
                // Invalid url
                completeCallback("Remove failed", "Invalid IP or port, try again", false)
            }
            let statusCode = response.response?.statusCode
            if statusCode == 200 {
                //Remove photo from PhotoArray, no need to reload all photos 
                self.PhotoArray.remove(at: index)
                completeCallback("Remove succesful", "Your photo has been succesfully removed", true)
            } else {
                completeCallback("Remove failed", "Invalid removal request, try again", false)
            }
        }
    }
}
