import Foundation

struct ClientSession {
    var token: Data
    var url: URL

    init(url: URL, token: Data) {
        self.token = token
        self.url = url
    }

    func upload(session: Session, completionHandler: @escaping(Data?, URLResponse?, Error?) -> Void) {
        let config = URLSessionConfiguration.default
        var request = URLRequest(url: url, cachePolicy: .reloadIgnoringCacheData, timeoutInterval: TimeInterval(10))
        request.httpMethod = "POST"
        request.httpBody = session.data()
        request.addValue("application/x-protobuf", forHTTPHeaderField: "Content-Type")
        request.setValue(token.base64EncodedString(), forHTTPHeaderField: "Authentication")

        URLSession(configuration: config).dataTask(with: request, completionHandler: completionHandler).resume()
    }

}
