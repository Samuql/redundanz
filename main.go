package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	env "github.com/Samuql/redundanz/environment"

	tea "github.com/charmbracelet/bubbletea"
)

// describes the application state
type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
	logger   *log.Logger
}

// stores initial application state
func initialModel(verbose bool) model {
	return model{

		choices: env.GetFolderSelection(), // currently string
		// A map which indicates which choices are selected. We're using
		// the map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
		logger:   log.New(os.Stdout, "", log.LstdFlags), // Writer can go to os.Stdout here
		// TODO UNDERSTAND THIS LINE
	}
}

// // Write writes len(b) bytes from b to the File.
// // It returns the number of bytes written and an error, if any.
// // Write returns a non-nil error when n != len(b).
// func (f *File) Write(b []byte) (n int, err error) {
// 	if err := f.checkValid("write"); err != nil {
// 		return 0, err
// 	}
// 	n, e := f.write(b)
// 	if n < 0 {
// 		n = 0
// 	}
// 	if n != len(b) {
// 		err = io.ErrShortWrite
// 	}

// 	epipecheck(f, e)

// 	if e != nil {
// 		err = f.wrapErr("write", e)
// 	}

// 	return n, err
// }

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("redundanz")
}

// handle updates and modify model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			// getFileNames for every selected element
			// every element ins the map selected should be iterated through
			files := make(chan []os.DirEntry)
			var wg sync.WaitGroup
			// for i := range m.selected {
			// 	fmt.Printf("== Folder %v/ ==\n", m.choices[i])
			// 	env.GetFileNames(m.choices[i], log.Default())
			// }
			for i := range m.selected {

				// Increment the WaitGroup counter
				wg.Add(1)
				// Start a new goroutine
				go func(i int) {
					// Decrement the WaitGroup counter when the goroutine finishes
					defer wg.Done()
					// Get the filenames
					env.GetFilesInDir(m.choices[i], log.Default(), files) // TODO FIX LOGGING

				}(i)
			}

			go func() {
				wg.Wait()
				close(files)
			}()

			i := 0
			for file := range files {
				fmt.Printf("\nfound files %s", file)
				i++
			}
			fmt.Printf("\n\nfound a total of %v files\n", i)

			// for filesl := range fnamesl {
			// 	// fmt.Printf("== Folder %s/ ==\n", filesl[0])
			// 	for _, file := range filesl {
			// 		fmt.Println(file)
			// 	}
			// }

			// ctn := []string{getWd()}
			// ctn = append(ctn, getContent(m.choices[m.cursor])...)
			// fmt.Printf("\n %v\n", ctn)
			return m, tea.Quit
		}
	}

	return m, nil
}

// render the application
func (m model) View() string {
	s := "Which folder do you want to eliminate redundancy in?\n[space] to select, [enter] to run.\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()
	p := tea.NewProgram(initialModel(*verbose))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
