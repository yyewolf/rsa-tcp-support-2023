use tui::widgets::ListState;

#[derive(Clone)]
pub struct StatefullList<T> {
    pub state: ListState,

    items: Vec<T>,
}

impl<T> StatefullList<T> {
    pub fn next(&mut self) {
        if self.items.len() == 0 {
            return;
        }

        let i = match self.state.selected() {
            Some(i) => {
                if i >= self.items.len() - 1 {
                    0
                } else {
                    i + 1
                }
            }
            None => 0,
        };
        self.state.select(Some(i));
    }

    pub fn prev(&mut self) {
        if self.items.len() == 0 {
            return;
        }

        let i = match self.state.selected() {
            Some(i) => {
                if i == 0 {
                    self.items.len() - 1
                } else {
                    i - 1
                }
            }
            None => 0,
        };
        self.state.select(Some(i));
    }

    pub fn unselect(&mut self) {
        self.state.select(None);
    }

    pub fn push(&mut self, item: T) {
        self.items.push(item);
    }

    pub fn iter(&self) -> std::slice::Iter<T> {
        self.items.iter()
    }
}

impl<T> Default for StatefullList<T> {
    fn default() -> Self {
        StatefullList {
            state: ListState::default(),
            items: Vec::new(),
        }
    }
}
