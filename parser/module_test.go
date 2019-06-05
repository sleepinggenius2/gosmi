package parser_test

const ModuleExample = `
FIZBIN-MIB DEFINITIONS ::= BEGIN

IMPORTS
    MODULE-IDENTITY, OBJECT-TYPE, experimental
        FROM SNMPv2-SMI;

fizbin MODULE-IDENTITY
    LAST-UPDATED "199505241811Z"
    ORGANIZATION "IETF SNMPv2 Working Group"
    CONTACT-INFO
            "        Marshall T. Rose

             Postal: Dover Beach Consulting, Inc.
                     420 Whisman Court
                     Mountain View, CA  94043-2186
                     US

                Tel: +1 415 968 1052
                Fax: +1 415 968 2510

             E-mail: mrose@dbc.mtview.ca.us"

    DESCRIPTION
            "The MIB module for entities implementing the xxxx
            protocol."
    REVISION      "9505241811Z"
    DESCRIPTION
            "The latest version of this MIB module."
    REVISION      "9210070433Z"
    DESCRIPTION
            "The initial version of this MIB module, published in
            RFC yyyy."
-- contact IANA for actual number
    ::= { experimental 101 }

END
`