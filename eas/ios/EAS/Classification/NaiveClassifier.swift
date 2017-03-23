import Foundation

class NaiveExerciseClassifier: ExerciseClassifier {
    private let signalAnalysis: SignalAnalysis

    init(signalAnalysis: SignalAnalysis) {
        self.signalAnalysis = signalAnalysis
    }

    func classify(matrix: Matrix) -> String? {
        let x = signalAnalysis.dft(matrix: matrix, n: 2).flatten()
        let (px1, cx1) = (x[0], x[1])
        let peps = Float(10)
        let ceps = Float(0.5)

        if abs(px1 - 86) < peps && abs(cx1 - 1.2) < ceps {
            return "straight bar biceps curl"
        } else if abs(px1 - 250) < peps && abs(cx1 - 1.5) < ceps {
            return "chest dumbbell flyes"
        }

        return nil
    }
}
