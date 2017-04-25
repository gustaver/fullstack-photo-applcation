//
// Created by Gustave Rousselet on 2017-03-27.
// Copyright (c) 2017 Gustave Rousselet. All rights reserved.
//

import Foundation
import SwiftyJSON
import Alamofire

class Photo: NSObject {

    var JpgBase64: String?
    var Title: String?
    var Description: String?
    var Date: String?
    var User: String?

    func initialiseFromJSON(data: JSON) {
        self.JpgBase64 = data["jpgbase64"].stringValue
        self.Title = data["title"].stringValue
        self.Description = data["description"].stringValue
        self.Date = data["date"].stringValue
        self.User = data["user"].stringValue
    }
    
    func toParameters() -> Parameters {
        let photoEncoded: Parameters = ["jpgBase64": self.JpgBase64, "title": self.Title, "description": self.Description, "date": self.Date, "user": self.User]
        return photoEncoded
    }
}
