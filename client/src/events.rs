use crate::client::{Client, Error};
use crossterm::event::{poll, read, Event, KeyCode, KeyEvent, KeyModifiers};
use std::{io, time::Duration};

pub enum Events {
    Event(Event),
    Tick,
}

pub fn read_events() -> Result<Events, io::Error> {
    if poll(Duration::from_millis(500))? {
        let event = read()?;
        match event {
            Event::Key(_) => Ok(Events::Event(event)),
            Event::Mouse(_) => Ok(Events::Event(event)),
            _ => Ok(Events::Tick),
        }
    } else {
        Ok(Events::Tick)
    }
}

pub fn handle_events(client: &mut Client) -> Result<(), Error> {
    let event = match read_events() {
        Ok(e) => e,
        Err(e) => return Err(Error::Err(e)),
    };

    match event {
        Events::Event(e) => match e {
            Event::Key(KeyEvent {
                code: KeyCode::Char('E'),
                modifiers: KeyModifiers::CONTROL,
            }) => client.elevate()?,
            Event::Key(KeyEvent {
                code: KeyCode::Char(c),
                modifiers: _,
            }) => client.input.push(c),
            Event::Key(KeyEvent {
                code: KeyCode::Enter,
                modifiers: _,
            }) => {
                client.send_message()?;
            }
            Event::Key(KeyEvent {
                code: KeyCode::Backspace,
                modifiers: _,
            }) => {
                client.input.pop();
            }
            Event::Key(KeyEvent {
                code: KeyCode::Esc,
                modifiers: _,
            }) => return Err(Error::Exit),
            _ => {}
        },
        Events::Tick => {}
    }

    Ok(())
}
