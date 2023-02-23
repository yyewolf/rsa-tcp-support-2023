use crate::message::Message;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Copy)]
pub enum PacketType {
    Hello = 0,
    Identify = 1,
    Success = 2,
    Error = 3,
    AgentCount = 4,
    Message = 5,
    Elevate = 6,
    ClientPresent = 7,
    ClientHistoryRequest = 8,
    ClientHistoryResponse = 9,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
pub struct Packet {
    #[serde(rename = "t")]
    pub packet_type: u8,

    #[serde(rename = "d")]
    pub data: PacketData,
}

#[derive(Serialize, Deserialize, Default, Debug, PartialEq, Clone)]
pub struct PacketData {
    #[serde(default, flatten)]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub message: Option<Message>,

    #[serde(default, flatten)]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub hello: Option<Hello>,

    #[serde(default, flatten)]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub identify: Option<Identify>,

    #[serde(default, flatten)]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub error: Option<Error>,

    #[serde(default, flatten)]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub agent_count: Option<AgentCount>,

    #[serde(default, flatten)]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub elevate: Option<Elevate>,

    #[serde(default, flatten)]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub client_present: Option<ClientPresent>,

    #[serde(default, flatten)]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub client_history_request: Option<ClientHistoryRequest>,

    #[serde(default, flatten)]
    #[serde(skip_serializing_if = "Option::is_none")]
    pub client_history: Option<ClientHistory>,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
pub struct Hello {
    #[serde(rename = "i")]
    id: u16,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
pub struct Identify {
    #[serde(rename = "l")]
    level: u16,

    #[serde(rename = "a")]
    auth: String,

    #[serde(rename = "n")]
    name: String,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
pub struct Error {
    #[serde(rename = "e")]
    pub error: String,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
pub struct AgentCount {
    #[serde(rename = "c")]
    pub count: u16,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
pub struct Elevate {
    #[serde(rename = "i")]
    #[serde(skip_serializing_if = "Option::is_none")]
    id: Option<u16>,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
pub struct ClientPresent {
    #[serde(rename = "i")]
    ids: Vec<u16>,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
pub struct ClientHistoryRequest {
    #[serde(rename = "i")]
    id: u16,
}

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
pub struct ClientHistory {
    #[serde(rename = "m")]
    messages: Vec<Message>,
}

impl Packet {
    pub fn message(message: Message) -> Self {
        Packet {
            packet_type: PacketType::Message as u8,
            data: PacketData {
                message: Some(message),
                ..Default::default()
            },
        }
    }

    pub fn identify(name: String, auth: String) -> Self {
        Packet {
            packet_type: PacketType::Identify as u8,
            data: PacketData {
                identify: Some(Identify {
                    level: 0,
                    auth,
                    name,
                }),
                ..Default::default()
            },
        }
    }

    pub fn elevate() -> Self {
        Packet {
            packet_type: PacketType::Elevate as u8,
            data: PacketData {
                elevate: Some(Elevate { id: None }),
                ..Default::default()
            },
        }
    }
}

impl PacketType {
    pub fn from_u8(value: u8) -> Option<Self> {
        match value {
            0 => Some(PacketType::Hello),
            1 => Some(PacketType::Identify),
            2 => Some(PacketType::Success),
            3 => Some(PacketType::Error),
            4 => Some(PacketType::AgentCount),
            5 => Some(PacketType::Message),
            6 => Some(PacketType::Elevate),
            7 => Some(PacketType::ClientPresent),
            8 => Some(PacketType::ClientHistoryRequest),
            9 => Some(PacketType::ClientHistoryResponse),
            _ => None,
        }
    }
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_error_serializer() {
        let e = Error {
            error: "test".to_string(),
        };
        let s = "{\"e\":\"test\"}";

        assert_eq!(serde_json::to_string(&e).unwrap(), s);
    }

    #[test]
    fn test_error_deserializer() {
        let e = Error {
            error: "test".to_string(),
        };
        let s = "{\"e\":\"test\"}";

        assert_eq!(serde_json::from_str::<Error>(s).unwrap(), e);
    }

    #[test]
    fn test_error_packet_serializer() {
        let e = Error {
            error: "test".to_string(),
        };
        let p = Packet {
            packet_type: PacketType::Hello as u8,
            data: PacketData {
                error: Some(e),
                ..Default::default()
            },
        };
        let s = "{\"t\":0,\"d\":{\"e\":\"test\"}}";

        assert_eq!(serde_json::to_string(&p).unwrap(), s);
    }

    #[test]
    fn test_error_packet_deserializer() {
        let e = Error {
            error: "test".to_string(),
        };
        let p = Packet {
            packet_type: PacketType::Hello as u8,
            data: PacketData {
                error: Some(e),
                ..Default::default()
            },
        };
        let s = "{\"t\":0,\"d\":{\"e\":\"test\"}}";

        assert_eq!(serde_json::from_str::<Packet>(s).unwrap(), p);

        let e = Error {
            error: "invalid character 'h' looking for beginning of value".to_string(),
        };
        let p = Packet {
            packet_type: PacketType::Error as u8,
            data: PacketData {
                error: Some(e),
                ..Default::default()
            },
        };
        let s = "{\"t\":3,\"d\":{\"e\":\"invalid character 'h' looking for beginning of value\"}}";

        assert_eq!(serde_json::from_str::<Packet>(s).unwrap(), p);
    }
}
