//
//  PhotoEditingViewController.swift
//  PhotoSharingApp
//
//  Created by Gustave Rousselet on 2017-03-30.
//  Copyright © 2017 Gustave Rousselet. All rights reserved.
//

import UIKit

class PhotoEditingViewController: UIViewController {
    
    // Image to be set from previous view controller (PhotoTableView)
    var image: UIImage!
    // ImageView displaying image taken by user (sent from previous view)
    @IBOutlet var imageView: UIImageView!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        // Set image passed from previous view controller (PhotoTableView)
        imageView.image = self.image
        // Do any additional setup after loading the view.
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }
    
    @IBAction func onUploadPhotoPress(_ sender: Any) {
        performSegue(withIdentifier: "PhotoTableView", sender: self)
    }

    /*
    // MARK: - Navigation

    // In a storyboard-based application, you will often want to do a little preparation before navigation
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        // Get the new view controller using segue.destinationViewController.
        // Pass the selected object to the new view controller.
    }
    */

}
