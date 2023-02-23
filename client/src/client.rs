use crate::conn::Conn;
use crate::events::handle_events;
use crate::message::{Message, MessageType};
use crate::packet::{Packet, PacketType};
use crate::statefull_list::StatefullList;
use crate::ui::ui;
use std::io;
use std::sync::mpsc::Receiver;
use std::thread;
use tui::backend::Backend;
use tui::Terminal;

#[derive(Clone, Default)]
pub struct Client {
    pub input: String,
    pub status: bool,
    pub level: u8,
    pub agent_count: u16,
    pub messages: StatefullList<Message>,

    conn: Conn,
}

pub enum Error {
    Exit,
    Err(io::Error),
}

impl Client {
    pub fn connect(&mut self) -> Result<(), io::Error> {
        self.conn.connect()
    }

    pub fn elevate(&mut self) -> Result<(), Error> {
        let p = Packet::elevate();
        match self.conn.write(p) {
            Ok(_) => {
                self.input.clear();
                Ok(())
            }
            Err(e) => Err(Error::Err(e)),
        }
    }

    fn wrap_input(&mut self) -> Result<Packet, io::Error> {
        let m = Message {
            msg: self.input.drain(..).collect(),
            to: String::default(),
            from: String::default(),
            message_type: MessageType::Outgoing,
        };

        let p = Packet::message(m);

        Ok(p)
    }

    pub fn send_message(&mut self) -> Result<(), Error> {
        let p = match self.wrap_input() {
            Ok(p) => p,
            Err(e) => return Err(Error::Err(e)),
        };
        let msg = p.data.message.clone().unwrap();
        match self.send(p) {
            Ok(_) => self.messages.push(msg),
            Err(e) => return Err(Error::Err(e)),
        }
        self.input.clear();
        Ok(())
    }

    pub fn handle_packet(&mut self, packet: &Packet) -> Result<(), io::Error> {
        match PacketType::from_u8(packet.packet_type) {
            Some(PacketType::Hello) => {
                let identify = Packet::identify(String::new(), String::new());
                self.send(identify)?;
            }
            Some(PacketType::Message) => {
                let msg = packet.data.message.clone().unwrap();
                match msg.from.as_str() {
                    "SYSTEM" => {
                        let (_, level) = msg.msg.split_once("level ").unwrap();
                        self.level = level.parse().unwrap();
                    }
                    _ => {
                        self.messages.push(msg);
                    }
                }
            }
            Some(PacketType::Error) => {
                // TODO: handle error
            }
            Some(PacketType::AgentCount) => {
                let msg = packet.data.agent_count.clone().unwrap();
                self.agent_count = msg.count;
            }
            Some(PacketType::Success) => {
                self.status = true;
            }
            _ => {
                println!("Unknown packet type received");
            }
        }

        Ok(())
    }

    pub fn send(&mut self, p: Packet) -> Result<(), io::Error> {
        self.conn.write(p)?;
        Ok(())
    }
}

pub fn run_client<B: Backend>(
    term: &mut Terminal<B>,
    client: &mut Client,
) -> Result<(), io::Error> {
    // connect to server
    client.connect()?;

    // start reading
    let c = client.clone();
    let rx = spawn_reading(&c)?;

    loop {
        // read
        if let Ok(packet) = rx.try_recv() {
            client.handle_packet(&packet)?;
        }

        // render ui
        ui(term, client)?;

        match handle_events(client) {
            Ok(_) => {}
            Err(e) => match e {
                Error::Exit => break,
                Error::Err(e) => return Err(e),
            },
        }
    }

    Ok(())
}

fn spawn_reading(client: &Client) -> Result<Receiver<Packet>, io::Error> {
    // start reading from server
    let mut conn = client.conn.clone();
    let (tx, rx) = std::sync::mpsc::channel::<Packet>();

    thread::spawn(move || loop {
        let packets = conn.read().unwrap();

        if packets.is_none() {
            continue;
        }

        packets.iter().for_each(|packets| {
            packets.iter().for_each(|packet| {
                tx.send(packet.clone()).unwrap();
            })
        });
    });

    Ok(rx)
}
