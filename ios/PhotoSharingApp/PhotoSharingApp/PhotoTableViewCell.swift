//
//  PhotoTableViewCell.swift
//  PhotoSharingApp
//
//  Created by Gustave Rousselet on 2017-03-28.
//  Copyright © 2017 Gustave Rousselet. All rights reserved.
//

import UIKit

// Custome table view cell for photos. 
class PhotoTableViewCell: UITableViewCell {

    
    @IBOutlet var imageInCell: UIImageView!
    @IBOutlet var blurView: UIVisualEffectView!
    @IBOutlet var titleLabel: UILabel!
    @IBOutlet var dateLabel: UILabel!
    @IBOutlet var descriptionLabel: UILabel!
    
    
    override func awakeFromNib() {
        super.awakeFromNib()
        // Initialization code
    }

    override func setSelected(_ selected: Bool, animated: Bool) {
        super.setSelected(selected, animated: animated)

        // Configure the view for the selected state
    }

}
