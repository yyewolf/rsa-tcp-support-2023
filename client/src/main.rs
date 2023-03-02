use client::client::{run_client, Client};
use crossterm::event::{DisableMouseCapture, EnableMouseCapture};
use crossterm::execute;
use crossterm::terminal::{
    disable_raw_mode, enable_raw_mode, EnterAlternateScreen, LeaveAlternateScreen,
};
use std::io;
use tui::backend::CrosstermBackend;
use tui::Terminal;

fn main() -> Result<(), io::Error> {
    let mut stdout = io::stdout();

    // setup terminal
    enable_raw_mode()?;

    execute!(stdout, EnterAlternateScreen, EnableMouseCapture)?;
    let backend = CrosstermBackend::new(stdout);
    let mut terminal = Terminal::new(backend)?;

    let mut client = Client::default();
    let run_time = run_client(&mut terminal, &mut client);

    // restore terminal
    disable_raw_mode()?;
    execute!(
        terminal.backend_mut(),
        LeaveAlternateScreen,
        DisableMouseCapture
    )?;
    terminal.show_cursor()?;

    match run_time {
        Ok(()) => {}
        Err(e) => println!("An error occured : {}", e),
    }

    Ok(())
}
