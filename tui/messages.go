package tui

import "github.com/LarsEckart/huey/hue"

type lightsLoadedMsg struct {
	lights []hue.Light
}

type groupsLoadedMsg struct {
	groups []hue.Group
}

type errMsg struct {
	err error
}

type lightToggledMsg struct {
	id    string
	newOn bool
}

type groupToggledMsg struct {
	id    string
	newOn bool
}

type lightRenamedMsg struct {
	id      string
	newName string
}

type groupRenamedMsg struct {
	id      string
	newName string
}

type groupDeletedMsg struct {
	id string
}

type groupCreatedMsg struct {
	group hue.Group
}

type scenesLoadedMsg struct {
	scenes []hue.Scene
}

type sceneActivatedMsg struct {
	id   string
	name string
}

type sceneCreatedMsg struct {
	scene hue.Scene
}

type sceneDeletedMsg struct {
	id string
}
