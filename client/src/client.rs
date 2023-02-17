mod events;
mod messages;
mod packet;
mod ui;

use crate::client::messages::Message;
use crate::client::packet::{Packet, PacketType};
use crossterm::event::{KeyCode, KeyEvent};
use serde_json::Deserializer;
use std::io;
use std::io::prelude::*;
use std::io::BufReader;
use std::net::{Ipv4Addr, TcpStream};
use std::sync::mpsc::Receiver;
use std::thread;
use tui::backend::Backend;
use tui::Terminal;

#[derive(Clone, Default)]
pub struct Client {
    input: String,
    messages: Vec<Message>,
    conn: Conn,
}

pub struct Conn {
    addr: Ipv4Addr,
    port: u16,
    stream: Option<TcpStream>,
}

pub enum Error {
    Exit,
    Err(io::Error),
}

impl Client {
    pub fn connect(&mut self) -> Result<(), io::Error> {
        let stream = TcpStream::connect(format!("{}:{}", self.conn.addr, self.conn.port))?;
        self.conn.stream = Some(stream);
        Ok(())
    }

    pub fn wrap_message(&mut self) -> Result<Packet, io::Error> {
        let m = Message {
            msg: self.input.drain(..).collect(),
            to: String::default(),
            from: String::default(),
            message_type: messages::MessageType::Outgoing,
        };

        let p = Packet::message(m);

        Ok(p)
    }

    pub fn handle_packet(&mut self, packet: &Packet) -> Result<(), io::Error> {
        match PacketType::from_u8(packet.packet_type) {
            Some(PacketType::Hello) => {
                let identify = Packet::identify(String::new(), String::new());
                self.send(identify)?;
            }
            Some(PacketType::Message) => {
                self.messages.push(packet.data.message.clone().unwrap());
            }
            Some(PacketType::Error) => {
                // TODO: handle error
            }
            Some(PacketType::AgentCount) => {
                // TODO: handle agent count
            }
            Some(PacketType::Success) => {}
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

    pub fn handle_key_event(&mut self, key: KeyEvent) -> Result<(), Error> {
        match key.code {
            KeyCode::Char(c) => {
                self.input.push(c);
            }
            KeyCode::Backspace => {
                self.input.pop();
            }
            KeyCode::Enter => {
                let p = match self.wrap_message() {
                    Ok(p) => p,
                    Err(e) => return Err(Error::Err(e)),
                };
                let msg = p.data.message.clone().unwrap();
                match self.send(p) {
                    Ok(_) => self.messages.push(msg),
                    Err(e) => return Err(Error::Err(e)),
                }
                self.input.clear();
            }
            KeyCode::Esc => return Err(Error::Exit),
            _ => {}
        }

        Ok(())
    }
}

impl Conn {
    pub fn read(&mut self) -> Result<Option<Vec<Packet>>, io::Error> {
        let stream = self.stream.as_mut().unwrap();
        let mut reader = BufReader::new(stream);
        let mut buf = [0; 1024];

        reader.read(&mut buf)?;

        // Remove null bytes from buffer
        let slice = buf
            .iter()
            .cloned()
            .skip_while(|&x| x != 123)
            .take_while(|&x| x != 0)
            .collect::<Vec<u8>>();

        if slice.len() == 0 {
            return Ok(None);
        }

        // Deserialize
        let stream_de = Deserializer::from_slice(&slice).into_iter::<Packet>();

        let mut packets = Vec::new();
        stream_de.into_iter().for_each(|p| {
            if let Ok(p) = p {
                packets.push(p);
            }
        });

        Ok(Some(packets))
    }

    pub fn write(&mut self, packet: Packet) -> Result<(), io::Error> {
        let stream = self.stream.as_mut().unwrap();
        let s = serde_json::to_string(&packet)?;
        stream.write(s.as_bytes())?;
        stream.flush()?;
        Ok(())
    }
}

impl Default for Conn {
    fn default() -> Self {
        Self {
            addr: Ipv4Addr::new(127, 0, 0, 1),
            port: 8000,
            stream: None,
        }
    }
}

impl Clone for Conn {
    fn clone(&self) -> Self {
        let stream = match &self.stream {
            Some(stream) => Some(stream.try_clone().unwrap()),
            None => None,
        };

        Self {
            addr: self.addr,
            port: self.port,
            stream: stream,
        }
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
        ui::ui(term, client)?;

        let event = events::read_events()?;

        match event {
            events::Events::Key(key) => match client.handle_key_event(key) {
                Ok(_) => {}
                Err(Error::Exit) => break,
                Err(Error::Err(e)) => return Err(e),
            },
            _ => {}
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
