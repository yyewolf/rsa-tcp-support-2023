use crossterm::event::{Event, KeyEvent, MouseEvent};
use std::io;

pub enum Events {
    Key(KeyEvent),
    Mouse(MouseEvent),
    Tick,
}

pub fn read_events() -> Result<Events, io::Error> {
    let event = crossterm::event::read()?;
    match event {
        Event::Key(key) => Ok(Events::Key(key)),
        Event::Mouse(mouse) => Ok(Events::Mouse(mouse)),
        _ => Ok(Events::Tick),
    }
}
