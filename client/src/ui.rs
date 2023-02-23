use crate::client::Client;
use crate::message::MessageType;
use std::io;
use tui::{
    backend::Backend,
    layout::{Alignment, Constraint, Direction, Layout},
    style::{Color, Modifier, Style},
    text::{Span, Spans},
    widgets::{Block, BorderType, Borders, List, ListItem, Paragraph, Row, Table},
    Terminal,
};

pub fn ui<B: Backend>(terminal: &mut Terminal<B>, client: &mut Client) -> Result<(), io::Error> {
    let help = display_help()?;
    let console = display_console(client)?;
    let message_list = display_body(client)?;
    let logs = display_logs(client)?;
    let messages_state = &mut client.messages.state;

    terminal.draw(|f| {
        let size = f.size();
        let block = Block::default();
        f.render_widget(block, size);

        // window constraints
        let window = Layout::default()
            .direction(Direction::Horizontal)
            .constraints([Constraint::Min(20), Constraint::Length(32)].as_ref())
            .split(size);

        // right panel constraints
        let right_panel = Layout::default()
            .direction(Direction::Vertical)
            .constraints([Constraint::Ratio(1, 3), Constraint::Ratio(2, 3)].as_ref())
            .split(window[1]);

        // body constraints
        let body = Layout::default()
            .direction(Direction::Vertical)
            .constraints([Constraint::Min(20), Constraint::Length(3)].as_ref())
            .split(window[0]);

        // render widgets
        f.render_stateful_widget(message_list, body[0], messages_state);
        f.render_widget(console, body[1]);
        f.render_widget(logs, right_panel[0]);
        f.render_widget(help, right_panel[1]);
    })?;

    Ok(())
}

fn display_help<'a>() -> Result<Table<'a>, io::Error> {
    let rows = vec![
        Row::new(vec!["<Esc>", "Quit"]),
        Row::new(vec!["<Ctr+e>", "Elevate"]),
        Row::new(vec!["<Enter>", "Send message"]),
        Row::new(vec!["<Backspace>", "Delete character"]),
        Row::new(vec!["<Up>", "Select previous message"]),
        Row::new(vec!["<Down>", "Select next message"]),
        Row::new(vec!["<Ctrl+u>", "Unselect message"]),
    ];

    Ok(Table::new(rows)
        .block(
            Block::default()
                .title("Help")
                .borders(Borders::ALL)
                .border_type(BorderType::Plain),
        )
        .widths(&[Constraint::Length(11), Constraint::Min(20)])
        .column_spacing(1))
}

fn display_console<'a>(client: &Client) -> Result<Paragraph<'a>, io::Error> {
    let text = vec![Spans::from(vec![
        Span::styled(">> ", Style::default().fg(Color::Green)),
        Span::styled(client.input.clone(), Style::default().fg(Color::Blue)),
    ])];

    let paraph = Paragraph::new(text)
        .block(Block::default().title("Console").borders(Borders::ALL))
        .style(Style::default().fg(Color::White))
        .alignment(Alignment::Left);

    Ok(paraph)
}

fn display_body<'a>(client: &Client) -> Result<List<'a>, io::Error> {
    let items: Vec<ListItem> = client
        .messages
        .iter()
        .map(|msg| {
            let spans = match msg.message_type {
                MessageType::Incoming => Some(vec![
                    Span::styled("<< ", Style::default().fg(Color::Yellow)),
                    Span::styled(msg.msg.clone(), Style::default().fg(Color::Blue)),
                ]),
                MessageType::Outgoing => Some(vec![
                    Span::styled(">> ", Style::default().fg(Color::Green)),
                    Span::styled(msg.msg.clone(), Style::default().fg(Color::Blue)),
                ]),
            };

            if spans.is_none() {
                return ListItem::new(vec![]);
            }

            ListItem::new(Spans::from(spans.unwrap()))
        })
        .collect();

    let list = List::new(items)
        .block(Block::default().title("Messages").borders(Borders::ALL))
        .highlight_style(Style::default().add_modifier(Modifier::BOLD));

    Ok(list)
}

fn display_logs<'a>(client: &Client) -> Result<List<'a>, io::Error> {
    let items = vec![
        (
            String::from("status : "),
            match client.status {
                true => ("Connected".to_string(), Color::Green),
                false => ("Disconnected".to_string(), Color::Red),
            },
        ),
        (
            String::from("level  : "),
            (client.level.to_string(), Color::Blue),
        ),
        (
            String::from("agents : "),
            (client.agent_count.to_string(), Color::Blue),
        ),
    ];

    let list_item: Vec<ListItem> = items
        .iter()
        .map(|(title, value)| {
            let spans = vec![
                Span::styled(title.clone(), Style::default().fg(Color::Yellow)),
                Span::styled(value.0.clone(), Style::default().fg(value.1)),
            ];

            ListItem::new(Spans::from(spans))
        })
        .collect();

    let list = List::new(list_item)
        .block(Block::default().title("Logs").borders(Borders::ALL))
        .style(Style::default().fg(Color::White));

    Ok(list)
}
