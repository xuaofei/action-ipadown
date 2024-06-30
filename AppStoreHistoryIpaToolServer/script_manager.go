package main

import "sync"

type ScriptManager struct {
}

var instanceSM *ScriptManager
var onceSM sync.Once

func GetSMInstance() *ScriptManager {
	onceSM.Do(func() {
		instanceSM = &ScriptManager{}
	})
	return instanceSM
}

func (receiver *ScriptManager) StartTask() {

}
