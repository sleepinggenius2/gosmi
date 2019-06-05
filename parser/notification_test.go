package parser_test

const NotificationTypeExample = `
entityMIBTraps      OBJECT IDENTIFIER ::= { entityMIB 2 }
entityMIBTrapPrefix OBJECT IDENTIFIER ::= { entityMIBTraps 0 }

entConfigChange NOTIFICATION-TYPE
    STATUS             current
    DESCRIPTION
            "An entConfigChange trap is sent when the value of
            entLastChangeTime changes. It can be utilized by an NMS to
            trigger logical/physical entity table maintenance polls.
            An agent must not generate more than one entConfigChange
            'trap-event' in a five second period, where a 'trap-event'
            is the transmission of a single trap PDU to a list of
            trap destinations.  If additional configuration changes
            occur within the five second 'throttling' period, then
            these trap-events should be suppressed by the agent. An
            NMS should periodically check the value of
            entLastChangeTime to detect any missed entConfigChange
            trap-events, e.g. due to throttling or transmission loss."
        ::= { entityMIBTrapPrefix 1 }
`