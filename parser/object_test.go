package parser_test

const ObjectIdentityExample = `
fizbin69 OBJECT-IDENTITY
    STATUS  current
    DESCRIPTION
        "The authoritative identity of the Fizbin 69 chipset."
    := { fizbinChipSets 1 }
`

const ObjectDefvalExample = `
ObjectSyntax       DEFVAL clause
----------------   ------------
Integer32          DEFVAL { 1 }
                   -- same for Gauge32, TimeTicks, Unsigned32
INTEGER            DEFVAL { valid } -- enumerated value
OCTET STRING       DEFVAL { 'ffffffffffff'H }
DisplayString      DEFVAL { "SNMP agent" }
IpAddress          DEFVAL { 'c0210415'H } -- 192.33.4.21
OBJECT IDENTIFIER  DEFVAL { sysDescr }
BITS               DEFVAL { { primary, secondary } }
                   -- enumerated values that are set
BITS               DEFVAL { { } }
                   -- no enumerated values are set
`

const ObjectTypeExample = `
evalSlot OBJECT-TYPE
    SYNTAX      Integer32 (0..2147483647)
    MAX-ACCESS  read-only
    STATUS      current
    DESCRIPTION
            "The index number of the first unassigned entry in the
            evaluation table, or the value of zero indicating that
            all entries are assigned.

            A management station should create new entries in the
            evaluation table using this algorithm:  first, issue a
            management protocol retrieval operation to determine the
            value of evalSlot; and, second, issue a management
            protocol set operation to create an instance of the
            evalStatus object setting its value to createAndGo(4) or
            createAndWait(5).  If this latter operation succeeds,
            then the management station may continue modifying the
            instances corresponding to the newly created conceptual
            row, without fear of collision with other management
            stations."
        ::= { eval 1 }

evalTable OBJECT-TYPE
    SYNTAX      SEQUENCE OF EvalEntry
    MAX-ACCESS  not-accessible
    STATUS      current
    DESCRIPTION
            "The (conceptual) evaluation table."
        ::= { eval 2 }
 
evalEntry OBJECT-TYPE
    SYNTAX      EvalEntry
    MAX-ACCESS  not-accessible
    STATUS      current
    DESCRIPTION
            "An entry (conceptual row) in the evaluation table."
    INDEX   { evalIndex }
        ::= { evalTable 1 }

EvalEntry ::=
    SEQUENCE {
        evalIndex       Integer32,
        evalString      DisplayString,
        evalValue       Integer32,
        evalStatus      RowStatus
    }

evalIndex OBJECT-TYPE
    SYNTAX      Integer32 (1..2147483647)
    MAX-ACCESS  not-accessible
    STATUS      current
    DESCRIPTION
            "The auxiliary variable used for identifying instances of
            the columnar objects in the evaluation table."
        ::= { evalEntry 1 }

evalString OBJECT-TYPE
    SYNTAX      DisplayString
    MAX-ACCESS  read-create
    STATUS      current
    DESCRIPTION
            "The string to evaluate."
        ::= { evalEntry 2 }

evalValue OBJECT-TYPE
    SYNTAX      Integer32
    MAX-ACCESS  read-only
    STATUS      current
    DESCRIPTION
            "The value when evalString was last evaluated, or zero if
             no such value is available."
    DEFVAL  { 0 }
        ::= { evalEntry 3 }

evalStatus OBJECT-TYPE
    SYNTAX      RowStatus
    MAX-ACCESS  read-create
    STATUS      current
    DESCRIPTION
            "The status column used for creating, modifying, and
            deleting instances of the columnar objects in the
            evaluation table."
    DEFVAL  { active }
        ::= { evalEntry 4 }
`