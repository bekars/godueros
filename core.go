package godueros

import (
//	"fmt"
)

type DuCore struct {
	directive *DuDirective
	event     *DuEvent
	mic       *DuMic
	speaker   *DuSpeaker
}

func NewDuCore()(core *DuCore, err error)  {
	err = nil
	core = &DuCore{
		directive: 	&DuDirective{},
		event: 		&DuEvent{},
		mic: 		&DuMic{},
		speaker: 	&DuSpeaker{},
	}

	return core, err
}

func (core *DuCore)Run() (err error) {
	err = nil
	/*
	directive: connect

	for {
	snowboy: wake up
		event: send voice server
		mic: record voice
		event: receive result
		directive: receive action
		speaker: play sound
	}
	 */

	return err
}
