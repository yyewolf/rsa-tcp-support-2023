use crate::packet::Packet;
use serde_json::Deserializer;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::net::{Ipv4Addr, TcpStream};

pub struct Conn {
    addr: Ipv4Addr,
    port: u16,
    stream: Option<TcpStream>,
}

impl Conn {
    pub fn connect(&mut self) -> Result<(), Error> {
        let stream = TcpStream::connect((self.addr, self.port))?;
        self.stream = Some(stream);
        Ok(())
    }

    pub fn read(&mut self) -> Result<Option<Vec<Packet>>, Error> {
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

    pub fn write(&mut self, packet: Packet) -> Result<(), Error> {
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
            stream,
        }
    }
}
