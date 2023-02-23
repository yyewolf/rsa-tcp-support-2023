use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Clone, Default, PartialEq, Debug)]
pub struct Message {
    #[serde(rename = "m")]
    pub msg: String,

    #[serde(rename = "t")]
    pub to: String,

    #[serde(rename = "f")]
    pub from: String,

    #[serde(skip_serializing, skip_deserializing)]
    pub message_type: MessageType,
}

#[derive(Clone, Debug, PartialEq)]
pub enum MessageType {
    Incoming,
    Outgoing,
}

impl Default for MessageType {
    fn default() -> Self {
        MessageType::Incoming
    }
}
