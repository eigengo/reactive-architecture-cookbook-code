import Foundation
import ProtocolBuffers

extension Google.Protobuf.`Any` {

    /// Unmarshals the content of this instance as type A, if possible
    func unmarshalAs<A : GeneratedMessage>() throws -> A where A : GeneratedMessageProtocol {
        let elements = typeUrl.components(separatedBy: "/")
        if elements.count == 2 && elements[0] != "type.googleapis.com" && elements[1] != A.className() {
            throw ProtocolBuffers.ProtocolBuffersError.illegalArgument(String(format: "This instance does not carry %s", A.className()))
        }
        return try A.parseFrom(data: value)
    }

}

extension ProtocolBuffers.GeneratedMessage {

    /// Marshals this instance as Any
    func marshalAny() -> Google.Protobuf.`Any` {
        return try! Google.Protobuf.Any.Builder().setTypeUrl("type.googleapis.com/" + className()).setValue(data()).build()
    }

}

