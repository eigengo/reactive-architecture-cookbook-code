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
        let cs = ClientSession(url: URL(string: "http://localhost:8080/session")!, token: Data())
        let sensor = try! Sensor.Builder().setLocation(SensorLocation.wrist).setDataTypes([SensorDataType.acceleration]).build()
        let values = [Float32](repeating: 0.987, count: 3600 * 50 * 3)
        let sensorData = try! SensorData.Builder().setSensors([sensor]).setValues(values).build()
        let s = try! Session.Builder().setSessionId(UUID().uuidString).setSensorData(sensorData).build()

        cs.upload(session: s) { data, response, error in
            if let data = data, let s = String(data: data, encoding: .utf8) { print("data = ", s) }
            if let response = response { print("response = ", response) }
            if let error = error { print("error = ", error) }
        }
    }

}

