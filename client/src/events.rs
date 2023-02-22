use crossterm::event::{poll, read, Event, KeyEvent, MouseEvent};
use std::{io, time::Duration};

pub enum Events {
    Key(KeyEvent),
    Mouse(MouseEvent),
    Tick,
}

pub fn read_events() -> Result<Events, io::Error> {
    if poll(Duration::from_millis(500))? {
        match read()? {
            Event::Key(k) => Ok(Events::Key(k)),
            Event::Mouse(m) => Ok(Events::Mouse(m)),
            _ => Ok(Events::Tick),
        }
    } else {
        Ok(Events::Tick)
    }
}
