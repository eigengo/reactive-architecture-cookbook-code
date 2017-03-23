import Foundation

protocol Classifier {

    func classify(sensors: Matrix) -> String?

}
