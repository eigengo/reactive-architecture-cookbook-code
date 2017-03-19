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
        let s = try! Session.Builder().setSessionId("sid").build()
        let e = try! Envelope.Builder().setCorrelationId("fooo").setToken("bar").setPayload(s.any()).build()
        let se = e.data()
        let e2 = try! Envelope.parseFrom(data: se)
        print(e2)
        let s2: Session = try! e2.payload.uas()
        print(s2)
    }

}

