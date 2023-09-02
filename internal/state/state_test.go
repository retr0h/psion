package state

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	off StateType = "Off"
	on  StateType = "On"

	switchOff EventType = "SwitchOff"
	switchOn  EventType = "SwitchOn"
)

type offAction struct{}

func (a *offAction) Execute(eventCtx EventContext) EventType {
	return NoOp
}

type onAction struct{}

func (a *onAction) Execute(eventCtx EventContext) EventType {
	return NoOp
}

func newLightSwitchFSM() *StateMachine {
	return &StateMachine{
		States: States{
			Defaults: State{
				Events: Events{
					switchOff: off,
				},
			},
			off: State{
				Action: &offAction{},
				Events: Events{
					switchOn: on,
				},
			},
			on: State{
				Action: &onAction{},
				Events: Events{
					switchOff: off,
				},
			},
		},
	}
}

var _ = Describe("State", func() {
	When("in off state", func() {
		lightSwitchFsm := newLightSwitchFSM()
		err := lightSwitchFsm.SendEvent(switchOff, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(lightSwitchFsm.Current).To(Equal(off))

		When("setting to off state", func() {
			It("should error", func() {
				err = lightSwitchFsm.SendEvent(switchOff, nil)
				Expect(err).To(HaveOccurred())
			})

			When("transitioning to on state", func() {
				It("should set state to on", func() {
					err := lightSwitchFsm.SendEvent(switchOn, nil)
					Expect(err).ToNot(HaveOccurred())
					Expect(lightSwitchFsm.Current).To(Equal(on))
				})
			})
		})
	})

	When("not in off state", func() {
		lightSwitchFsm := newLightSwitchFSM()

		When("transitioning to on state", func() {
			It("should error", func() {
				err := lightSwitchFsm.SendEvent(switchOn, nil)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
