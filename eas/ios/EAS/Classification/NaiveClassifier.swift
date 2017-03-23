import Foundation

class NaiveClassifier : Classifier {
    private let signalAnalysis: SignalAnalysis

    init(signalAnalysis: SignalAnalysis) {
        self.signalAnalysis = signalAnalysis
    }

    func classify(sensors: Matrix) -> String? {
        return nil
    }
}
