package main

func ng() {
	activeForm = &form{
		Fn: createNewGroup,
	}
}

func createNewGroup(f *form, input string) {
	if f.Stage == 3 {
		activeForm = nil
	}
}
