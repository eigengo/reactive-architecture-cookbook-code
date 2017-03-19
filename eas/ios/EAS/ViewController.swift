//
//  ViewController.swift
//  EAS
//
//  Created by Jan Machacek on 19/03/2017.
//  Copyright © 2017 Jan Machacek. All rights reserved.
//

import UIKit
import ProtocolBuffers

class ViewController: UIViewController {

    override func viewDidLoad() {
        super.viewDidLoad()
        foo()
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }

    func foo() {
        let s = try! Session.Builder().build()
        Google.Protobuf.Any().
        let e = try! Envelope.Builder().setCorrelationId("fooo").setToken("bar").build()
        let se = e.data()
        let e2 = try! Envelope.parseFrom(data: se)
        print(e2)
    }

}

